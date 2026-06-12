# Phase 2 — Real LIC API Integration

This document covers everything required to move from manual/CSV-based data entry to live
synchronization with LIC's official agent APIs. Every item below corresponds to a
`TODO(LIC-API — Phase 2)` marker already left in the codebase.

> **Note**: LIC does not currently publish a public agent API. The endpoint paths below
> (`https://licindia.in/api/v2/...`) are placeholders used in code comments to describe the
> *shape* of integration expected once LIC provides (or an authorized intermediary provides)
> such an API. Replace base URL/paths with the real contract once available.

---

## 1. Agent authentication via LIC (`domain/auth.go`)

**Current state**: agents authenticate purely against our local `agent` table
(bcrypt password hash, JWT issued by us).

**What's missing**:
- On first login with `agentCode`, call `POST /api/v2/auth/login` with agent credentials to
  obtain LIC's own `userkey` + `authtoken`.
- Persist `userkey`/`authtoken` on the `agent` row (columns already exist: `userkey` is now
  `TEXT` after migration 000025; `authtoken` would need a new column/migration).
- **Device binding**: compute a `machine_id` (hash of device fingerprint — e.g. hostname +
  MAC + app install ID) and send it with the LIC login call so LIC ties the token to this
  device.
- **MPIN offline re-auth**: allow agents to set a 4-6 digit MPIN; if LIC's servers are
  unreachable, validate MPIN locally (hashed, like password) and continue using the
  last-cached `authtoken` until it expires.
- Encrypt `authtoken`/`userkey` at rest using the existing AES-256-GCM helper in
  `domain/crypto.go` (`Encrypt`/`Decrypt`, key derived from `JWT_SECRET`) — same pattern
  already used for `aadhar`/`pan`/`ckyc`.

**To build**:
- New migration: `agent.authtoken TEXT`, `agent.machine_id VARCHAR(128)`, `agent.mpin_hash TEXT`.
- New `domain/licapi` package: HTTP client wrapping LIC auth, policy, plan, commission, FUP
  endpoints, with retry/backoff and token refresh.
- `AuthService.Login` extended: after local password check succeeds, if `agent.authtoken` is
  empty/expired and network available, call LIC login and store the result.

---

## 2. Policy sync (`domain/policy.go:41`)

**Current state**: policies are entered manually or via the bulk SQL import
(`POST /agent/import`).

**What's missing**:
- On agent login (or a manual "Sync Now" action), call `GET /api/v2/policies/list` using the
  stored `authtoken`.
- Diff returned policies against `policies` table by `policy_no`:
  - New policy numbers → INSERT.
  - Existing → UPDATE mutable fields (status, next_premium, sum_assured if endorsed, etc.)
  - LIC is authoritative for `status`/`stat_cd` and `next_premium` (see §6 and §5).
- Map LIC's response schema → our `CreatePolicyRequest`/`policies` columns (field-by-field
  mapping table needed once LIC's actual JSON schema is known — currently unknown).

**To build**:
- `PolicySyncService.SyncFromLIC(agentID)`:
  1. Fetch list from LIC.
  2. For each item, upsert into `policies` (transaction).
  3. After upsert, call `ReportRepo.RefreshCalendar`/`RefreshCashflow` for affected families
     (cached reports go stale otherwise — same caveat as manual edits, documented in
     FRONTEND_API_GUIDE.md §Reports).
- New endpoint: `POST /policies/sync` (auth required) — triggers the above on demand, returns
  `{ "added": n, "updated": n, "skipped": n }`.
- Scheduled job (cron or on-login) to run sync automatically.

---

## 3. Plan catalogue sync (`domain/plan.go:21`)

**Current state**: `plans` table is a static seed of 269 LIC plans (sbSchedule, lapsDays,
gstRates baked in at seed time).

**What's missing**:
- `GET /api/v2/plans/list` on agent login to pick up new plan launches / withdrawn plans /
  rate changes without redeploying the app.
- Plan updates must not break existing `policies.plan_no` foreign-key references — sync
  should be an upsert keyed on `plan_no`, never a delete of plans still referenced by policies.

**To build**:
- `PlanSyncService.SyncFromLIC()` — upsert into `plans` by `plan_no`.
- Run on a longer cadence (e.g. daily/weekly) rather than every login, since plan data
  changes infrequently — avoid hammering LIC's API.
- Add `plans.updated_at` column (new migration) to track sync freshness; surface "Plan data
  last updated: <date>" in the frontend.

---

## 4. Commission sync + rate table (`domain/commission.go:32` and `:60`)

**Current state**: `/commission/calculate` returns an **estimate** using hardcoded
approximate percentages; `/commission` records are entered manually as they're paid.

**What's missing**:
- **Actual commission sync**: `GET /api/v2/commission/list` — pull real paid-commission
  records (bill_date, first/second/third/bonus/single comm amounts) and reconcile against the
  `commission` table (insert missing bill periods; never overwrite a manually-entered record
  unless LIC's amount differs — flag discrepancies instead of silently overwriting).
- **Commission rate table**: build `commission_rates(plan_no, term, premium_year, base_pct,
  bonus_pct)` so `/commission/calculate` returns an *exact* figure instead of an approximation.
  Source for these rates: either scraped from LIC circulars or returned by a (currently
  hypothetical) `GET /api/v2/commission/rates?planNo=&term=` endpoint.

**To build**:
- New migration: `commission_rates` table (see schema above) + seed data once rates are
  sourced.
- `CommissionService.Calculate` updated to look up `commission_rates` first, fall back to the
  current approximation (with the existing `note` field) if no row matches.
- `CommissionSyncService.SyncFromLIC(agentID, year)` — reconciliation job, new endpoint
  `POST /commission/sync`.

---

## 5. FUP update sync queue (`domain/fup.go:27`)

**Current state**: `POST /fup/update` writes to `fuphistory` + `fupupdate` + updates
`policies.next_premium`/`fupstatus`/`lastfup` — all **local only**.

**What's missing**:
- `fupupdate` should function as an **outbox queue**: each row represents a FUP change made
  locally that hasn't yet been pushed to LIC.
- Background worker drains `fupupdate` (where `synced_at IS NULL`), calling
  `POST /api/v2/fup/update` for each row with `{ policyNo, oldFup, newFup }`.
- On success: mark `fupupdate.synced_at = now()`.
- On failure due to conflict (LIC's `nextPremium` differs from our `oldFup`): **LIC wins** —
  overwrite our `policies.next_premium`/`fupstatus`/`lastfup` with LIC's value, write a
  corrective `fuphistory` entry noting the conflict resolution, and mark the queue row as
  `conflict_resolved` rather than retrying indefinitely.

**To build**:
- New migration: `fupupdate.synced_at TIMESTAMPTZ NULL`, `fupupdate.sync_status VARCHAR(20)
  DEFAULT 'pending'` (`pending`/`synced`/`conflict_resolved`/`failed`).
- `FUPSyncWorker.Run()` — periodic job processing pending queue rows.
- Endpoint `GET /fup/sync-status` — lets the frontend show "3 updates pending sync to LIC".

---

## 6. Policy status refresh (`domain/policy.go:131`)

**Current state**: `policies.status`/`stat_cd` only change via manual `PUT /policies/:policyNo`.

**What's missing**:
- Status (IF/LA/PU/SU/MA/CL/EX) should be refreshed from LIC as part of the policy sync in
  §2 — LIC is authoritative for lapse/maturity/claim/surrender events, which often happen
  without the agent's direct action.
- When status changes via sync (e.g. IF → LA), this should also trigger:
  - A `fuphistory` note (system-generated, not agent-entered).
  - A report cache invalidation (`rpt_currstatus`, `rpt_cashflow`) for the affected family.

**To build**: covered by the same `PolicySyncService` in §2 — status is just one of the
fields included in the upsert diff. Add an `activities`/notification row when a policy's
status changes to LA/SU/MA/CL so the agent is alerted (ties into the existing
`prospect_activity`/leads table).

---

## 7. GST regime history (`domain/gst.go:33`)

**Current state**: `/gst/calculate` hardcodes a single cutover date (2025-09-22) and two GST
rate sets in code.

**What's missing**:
- A `gst_regime(effective_date DATE, rate_pct NUMERIC, ...)` table so future GST rate changes
  (which happen periodically via government notification, independent of LIC) don't require a
  code deploy.
- `GSTService.Calculate` looks up the applicable row by `effective_date <= premium due date`,
  ordered descending, limit 1 — replacing the current `if date.Before(cutover)` branch.

**To build**:
- New migration: `gst_regime` table, seeded with the two currently-hardcoded regimes
  (pre/post 2025-09-22).
- This item is **not** LIC-API-dependent — it can be implemented immediately and
  independently of the rest of Phase 2, and removes one of the two hardcoded-data TODOs.

---

## 8. Cross-cutting infrastructure needed for all of the above

1. **`domain/licapi` client package** — single place for: base URL config (`LIC_API_BASE_URL`
   env var), auth header injection (`authtoken`), retry/backoff, timeout, and structured
   error mapping (LIC error codes → our `{error, message}` envelope).
2. **Background job runner** — currently the app has no scheduler. Either:
   - a lightweight in-process ticker (`time.AfterFunc`/`time.Ticker` goroutine in `main.go`),
     guarded by `LIC_SYNC_ENABLED=true` env flag, or
   - an external cron hitting new `/policies/sync`, `/commission/sync`, `/fup/sync` endpoints.
3. **Sync audit log** — a `sync_log(agent_id, sync_type, started_at, finished_at, status,
   detail)` table so failures are visible/debuggable and the frontend can show "last synced".
4. **New env vars** (add to `.env.example`):
   ```
   LIC_API_BASE_URL=https://licindia.in/api/v2
   LIC_SYNC_ENABLED=false
   LIC_SYNC_INTERVAL_MINUTES=60
   ```
5. **Encryption**: all LIC tokens stored via the existing `domain.Encrypt`/`Decrypt`
   (AES-256-GCM, key from `JWT_SECRET`) — no new crypto needed, just apply it to the new
   `authtoken`/`machine_id`/`mpin_hash` columns.

---

## 9. Suggested implementation order

1. §7 GST regime table — no external dependency, quick win.
2. §8 infra (licapi client skeleton, sync_log table, env vars) — foundation for everything else.
3. §1 LIC auth (userkey/authtoken/machine_id/MPIN) — required before any other LIC call can
   be authenticated.
4. §2 policy sync + §6 status refresh (same service).
5. §3 plan catalogue sync.
6. §5 FUP outbox queue + worker.
7. §4 commission sync + rate table (most data-dependent — needs sourced rate table first).

Each numbered item above is independently shippable once its prerequisites are in place, and
each leaves the system in a working state if LIC's real API contract turns out to differ from
the placeholder paths assumed here — only the `domain/licapi` client and field-mapping layers
would need adjustment.

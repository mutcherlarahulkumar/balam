# Balam API — Frontend Integration Guide

Base URL: `/v1` (e.g. `https://your-host/v1/auth/login`)
Health check (no `/v1` prefix, no auth): `GET /health` → `{"status":"ok"}`

---

## 1. Authentication

All routes except `POST /auth/login` and `POST /auth/register` require:

```
Authorization: Bearer <jwt>
```

Missing header → `401 {"error":"missing_token","message":"Authorization header is required"}`
Invalid/expired → `401 {"error":"invalid_token","message":"Token is invalid or expired"}`

Tokens expire after `JWT_EXPIRY_HOURS` (default 24h). Call `POST /auth/refresh` (with the
current still-valid token) to get a new one before it expires.

---

## 2. Standard error envelope

Every error response (except validation errors) looks like:

```json
{ "error": "error_code", "message": "Human readable message." }
```

Validation errors (HTTP 400) look like:

```json
{
  "error": "validation_error",
  "errors": [
    { "field": "email", "message": "email must be a valid email address." },
    { "field": "mobile", "message": "mobile must be exactly 10 characters." }
  ]
}
```

Show `errors[].message` directly next to the corresponding form field (`errors[].field`
matches the JSON key you sent).

Malformed JSON body → `400 {"error":"invalid_json","message":"Request body is not valid JSON."}`

---

## 3. Routes

### Auth (public)

| Method | Path | Body | Notes |
|---|---|---|---|
| POST | `/auth/login` | `LoginRequest` | identifier = agent code OR email |
| POST | `/auth/register` | `RegisterRequest` | creates new agent |
| POST | `/auth/refresh` | — (auth) | returns new token |
| POST | `/auth/change-password` | `ChangePasswordRequest` (auth) | |

**LoginRequest**
```json
{ "identifier": "0045269T", "password": "Lic@8074" }
```
- `identifier`: required
- `password`: required, min 6 chars

**Login responses**
- `200` → `{ "token", "expiresAt", "agent": AgentPublic }`
- `401 invalid_credentials` — wrong identifier/password
- `401 account_locked` — 5+ failed attempts in last 15 min
- `401 account_terminated` — agent.terminated = true

**RegisterRequest**
```json
{
  "name": "M Balarama Murty",
  "email": "agent@example.com",
  "agentCode": "0045269T",
  "password": "Lic@8074!",
  "branch": "Vijayawada",
  "mobile": "9876543210",
  "licenceNo": "1234567890"
}
```
| Field | Rule |
|---|---|
| name | required, 2–80 chars |
| email | required, valid email |
| agentCode | required, 3–20 chars (alphanumeric LIC codes like `0045269T` allowed) |
| password | required, 8–72 chars |
| branch | required, 2–20 chars |
| mobile | required, exactly 10 digits |
| licenceNo | required, 3–100 chars |

- `201` → `{ "message", "agentId" }`
- `409 email_taken` / `409 agent_code_taken`

**ChangePasswordRequest**
```json
{ "currentPassword": "old", "newPassword": "new12345" }
```
- `400 wrong_current_password` if `currentPassword` doesn't match

---

### Agent profile (auth required)

| Method | Path | Body |
|---|---|---|
| GET | `/agent/profile` | — |
| PUT | `/agent/profile` | `UpdateProfileRequest` |

**UpdateProfileRequest** — all fields optional
```json
{ "name": "...", "mobile": "...", "email": "...", "photo": "...", "slogan": "...", "address": "..." }
```
| Field | Rule |
|---|---|
| name | 2–80 |
| mobile | 10–15 |
| email | valid email |
| slogan | max 200 |
| address | max 200 |

---

### Agent data import (auth required)

| Method | Path | Body |
|---|---|---|
| POST | `/agent/import` | multipart `file` field |

Clears **all** tables (agent, families, clients, policies, plans, commission, sb, fup*, loan,
leads, activities, documents, bankdetails, report caches), then runs your INSERT-only `.sql`
file in one transaction.

- File must contain **only** `INSERT` and `SELECT setval(...)` statements. No `TRUNCATE`,
  `DROP`, `ALTER`, etc — server handles the wipe.
- Max 50 MB.
- `403 disallowed_sql` if any other statement type is present
- `422 import_failed` with the raw Postgres error message if INSERTs fail (constraint
  violation etc.) — whole import is rolled back, nothing is partially applied.

There is also `POST /admin/import-sql` which accepts a **full seed file including
`TRUNCATE`/`BEGIN`/`COMMIT`** — used for the one-time SQLite→Postgres migration dump.

---

### Families (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/families?page=&limit=&search=` | paginated, `search` matches head name or family code |
| POST | `/families` | `CreateFamilyRequest` |
| GET | `/families/:familyCode` | returns `FamilyDetail` (members + policies) |
| PUT | `/families/:familyCode` | `UpdateFamilyRequest` |

**CreateFamilyRequest**
| Field | Rule |
|---|---|
| familyCode | optional, max 15 (auto-generated if omitted) |
| headName | **required**, 2–80 |
| address | max 250 |
| mobile | 10–15 |
| email | valid email |
| pincode | max 10 |
| religion | max 20 |
| designation | max 80 |

`UpdateFamilyRequest` — same fields, all optional.

List response: `{ "data": [FamilyListItem], "total", "page", "limit" }`
`limit` is capped at 100.

---

### Clients (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/clients?page=&limit=&familyCode=` | paginated |
| POST | `/clients` | `CreateClientRequest` |
| GET | `/clients/search?q=&limit=` | search by name/mobile, `q` min 2 chars |
| GET | `/clients/:id` | returns `ClientDetail` (policies + bank details + documents) |
| PUT | `/clients/:id` | `UpdateClientRequest` |

**CreateClientRequest**
| Field | Rule |
|---|---|
| familyCode | **required**, max 15 |
| persCode | **required**, max 10 |
| name | **required**, 2–80 |
| dob | optional, `YYYY-MM-DD` |
| sex | optional, one of `M F O` |
| mobile | 10–15 |
| email | valid email |
| occupation | max 50 |
| clientType | one of `C P N` (Customer/Prospect/New) |
| address | max 255 |

`UpdateClientRequest` — same fields, all optional.

Bank details (`bankDetails[]` in `ClientDetail`) — `aadhar`, `pan`, `ckyc` are encrypted at
rest (AES-256-GCM) and decrypted automatically in the response.

---

### Plans (auth required)

| Method | Path |
|---|---|
| GET | `/plans` |

Returns `{ "data": [Plan] }` — 269 LIC plans with `sbSchedule` (survival benefit %
per year for money-back plans), `lapsDays`, and current `gstRates`.

---

### Policies (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/policies?page=&limit=&familyCode=&status=&fupStatus=` | paginated |
| POST | `/policies` | `CreatePolicyRequest` |
| GET | `/policies/:policyNo` | returns `PolicyDetail` (FUP history + loans + SB records) |
| PUT | `/policies/:policyNo` | `UpdatePolicyRequest` |

**Status codes** (`policies.status`):
| Code | Meaning |
|---|---|
| IF | In Force |
| LA | Lapsed |
| PU | Paid-Up |
| SU | Surrendered |
| MA | Matured |
| CL | Claim |
| EX | Expired |

**FUP status** (`policies.fupstatus`, computed): `PAID`, `DUE`, `OVERDUE`, `LAPSED`

**Payment modes**: `Y` yearly, `H` half-yearly, `Q` quarterly, `M` monthly, `S` single

**New response field — `premiumEndDate`**: returned on `GET /policies` (list items) and
`GET /policies/:policyNo` (detail). Computed as `issueDate + ppt years` — the date premium
payments stop (premium paying term end date). `null` if `issueDate` is not set. Existing
fields `plan` (numeric plan code) and `planName` were already present in both responses.

**New response field — `planDetails`** (detail only): `GET /policies/:policyNo` now also
returns a `planDetails` object — the full plan record from `GET /plans` (same shape as the
`PlanResponse` schema: `planNo`, `planName`, `planType`, `termPpt`, `sbSchedule`, `stax`,
`lapsDays`, `gstRates`). `null` if the policy's `plan` code doesn't match any row in `plans`.
Use this instead of doing a separate `/plans` lookup per policy.

**CreatePolicyRequest**
| Field | Rule |
|---|---|
| policyNo | **required**, integer, unique |
| familyCode | **required** |
| persCode | **required** |
| planNo | **required** (LIC plan number, string) |
| issueDate | **required**, `YYYY-MM-DD` |
| matDate | **required**, `YYYY-MM-DD` |
| term | **required**, ≥1 |
| ppt | **required**, ≥1 (premium paying term) |
| sumAssured | **required**, > 0 |
| premium | **required**, > 0 |
| paymentMode | **required**, one of `Y H Q M S` |
| nextPremium | **required**, `YYYY-MM-DD` |
| nominee | **required**, max 80 |
| relation | **required**, max 20 |
| branch | optional, max 15 |
| neft | optional, `YES` or `NO` |
| dab | optional, ≥0 |
| termRider | optional, ≥0 |

**UpdatePolicyRequest** — all optional
| Field | Rule |
|---|---|
| status | one of `IF LA PU SU MA CL EX` |
| nominee | max 80 |
| relation | max 20 |
| neft | `YES`/`NO` |
| nextPremium | `YYYY-MM-DD` |
| fupStatus | free text |

---

### FUP — Follow-Up Premium (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/fup/due?year=&month=&overdueDays=` | due/overdue premium list |
| POST | `/fup/update` | `UpdateFUPRequest` |
| GET | `/fup/history/:policyNo` | full audit trail |
| GET | `/fup/multipledue/:policyNo` | instalment arrears |

**`/fup/due` query params** (all optional, combine freely):
- `year` — `YYYY`, filters by `next_premium` year
- `month` — `1`–`12`, filters by `next_premium` month
- `overdueDays` — only policies overdue by ≥ N days

Each item includes `daysOverdue`, `lapseDate`, `daysUntilLapse` (computed from
`plans.lapsdays`, defaults to 180 if plan not found).

**Calendar-style month/year browsing**: when `year` and/or `month` is supplied, the endpoint
returns **all** policies whose `nextPremium` falls in that calendar month/year — past,
present, or future — so the frontend can build a month picker (e.g. "June 2025") and show
everything due that month, not just overdue items. When `year`/`month` are both omitted, the
endpoint defaults to the original behaviour: only premiums due **on or before today**
(current due/overdue list).

**UpdateFUPRequest**
```json
{ "policyNo": 123456789, "oldFup": "2024-04-01", "newFup": "2025-04-01", "reason": "" }
```
- `oldFup` **must exactly match** the policy's current `nextPremium` (validated server-side)
- `404 policy_not_found`
- `422 fup_mismatch` — "Provided oldFup does not match the current next premium on the policy"
  → re-fetch the policy and retry with the correct `oldFup`

On success, writes to both `fuphistory` and `fupupdate`, updates `policies.next_premium`,
`fupstatus`, and `lastfup`.

---

### Commission (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/commission?year=&month=` | filter by bill_date |
| POST | `/commission` | `CreateCommissionRequest` |
| GET | `/commission/summary` | current month + last 5 years aggregated |
| GET | `/commission/calculate?policyNo=&year=` | estimate |

`year`/`month` are plain integers (e.g. `?year=2024&month=4`).

**CreateCommissionRequest**
| Field | Rule |
|---|---|
| policyNo | **required** |
| billDate | **required**, `YYYY-MM-DD` |
| firstComm, secondComm, thirdComm, bonusComm, singleComm | optional, ≥0 |
| payDate | optional, `YYYY-MM-DD` |

**`/commission/calculate`**: `year` = premium year (1 = first year). Returns
`baseCommissionPct`, `bonusCommissionPct`, `totalPct`, `estimatedCommission`, and a `note`
explaining the rates are approximate (see Phase 2 doc — needs a proper rate table).

---

### GST (auth required)

| Method | Path |
|---|---|
| GET | `/gst/calculate?policyNo=&premiumYear=` |

Returns GST breakdown for a premium. `premiumYear` defaults to 1. Uses a hardcoded GST
regime cutover date (2025-09-22) — see Phase 2 doc for making this data-driven.

---

### Loans (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/loans?policyNo=` | optional filter |
| POST | `/loans` | `CreateLoanRequest` |

**CreateLoanRequest**
```json
{ "policyNo": 123456789, "loanDate": "2024-01-15", "loanAmount": 50000, "interestDueDate": "2025-01-15", "loanInterest": 5000 }
```
| Field | Rule |
|---|---|
| policyNo | **required** |
| loanDate | **required**, `YYYY-MM-DD` |
| loanAmount | **required**, integer > 0 (whole rupees — DB column is INTEGER) |
| interestDueDate | **required**, `YYYY-MM-DD` |
| loanInterest | optional, integer ≥0 |

> ⚠️ Amounts are **integers** (no paise/decimals). Round on the client before sending.

---

### Survival Benefits (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/sb?year=&month=&unpaidOnly=` | filter |
| POST | `/sb` | `CreateSBRequest` |
| PUT | `/sb/:id/mark-paid` | `MarkSBPaidRequest` |

`year` (`YYYY`), `month` (`1`-`12`), `unpaidOnly` (`true`/`false`) — filters on `sb_duedate`.

**CreateSBRequest**
| Field | Rule |
|---|---|
| policyNo | **required** |
| sbDueDate | **required**, `YYYY-MM-DD` |
| sbAmount | **required**, > 0 |
| instalmentNo | **required**, ≥1 |

**MarkSBPaidRequest**
```json
{ "paidDate": "2026-04-05", "chequeNo": "123456" }
```
- `404 sb_not_found`
- `409 sb_already_paid` — record already has `sbPayDate` set

---

### Leads & Activities (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/leads` | all leads |
| POST | `/leads` | `CreateLeadRequest` |
| GET | `/activities?clientId=` | optional filter |
| POST | `/activities` | `CreateActivityRequest` |
| GET | `/activities/today` | today's scheduled activities |

**CreateLeadRequest**
| Field | Rule |
|---|---|
| name | **required**, 2–150 |
| mobile | **required**, 10–15 |
| address | max 500 |
| searchTerm | max 100 |

**CreateActivityRequest**
| Field | Rule |
|---|---|
| clientId | **required** |
| activityType | **required**, one of `CALL MEETING DEMO EMAIL PROPOSAL MEDICAL REMINDER` |
| activityDate | **required**, `YYYY-MM-DD` |
| activityTime | optional, `HH:MM` |
| details | max 1000 |
| reminderDate | optional, `YYYY-MM-DD` |
| reminderTime | optional |
| status | optional, one of `PENDING DONE CANCELLED` |

---

### Dashboard (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/dashboard` | aggregated home-screen summary |

Single call combining everything the home/landing screen typically needs, instead of 5
separate requests:

```json
{
  "duePremiums": { "total": 12, "preview": [ /* up to 5 FUPDueItem, see /fup/due */ ] },
  "todayActivities": [ /* Activity[], see /activities/today */ ],
  "commissionThisMonth": { "month": "2026-06", "totalCommission": 12345.0, ... },
  "unpaidSB": { "total": 3, "preview": [ /* up to 5 SB, see /sb */ ] },
  "leads": { "total": 7, "preview": [ /* up to 5 Lead, most recent first */ ] }
}
```

`total` reflects the full count; `preview` is capped at 5 items — for the full list, call the
underlying endpoint (`/fup/due`, `/sb?unpaidOnly=true`, `/leads`).

---

### Reports (auth required)

| Method | Path | Notes |
|---|---|---|
| GET | `/reports/cashflow?familyCode=` | SB + maturity cashflow timeline |
| GET | `/reports/cashinout` | yearly cash in/out across all policies |
| GET | `/reports/status?familyCode=` | current status snapshot |
| GET | `/reports/calendar?familyCode=` | 12-month premium calendar per policy |
| POST | `/reports/refresh` | `{ "familyCode": "..." }` — rebuilds cashflow + calendar caches |

`familyCode` is **required** on cashflow/status/calendar — `400 missing_family_code` if
omitted. `404 family_not_found` if the family doesn't exist.

**Important**: Cashflow and calendar are **cached** tables (`rpt_cashflow`,
`rpt_prcalendar`). Call `POST /reports/refresh` after creating/editing policies or SB
records for a family, otherwise the cached report will be stale.

Calendar logic (fixed): for payment mode `Q`/`H`/`Y`, the months shown are computed
relative to the policy's `next_premium` month (e.g. a quarterly policy due in March shows
Mar/Jun/Sep/Dec, not a fixed set of months).

---

## 4. Pagination

All paginated list endpoints return:
```json
{ "data": [...], "total": 123, "page": 1, "limit": 20 }
```
`limit` defaults to 20, capped at 100. `page` is 1-indexed.

---

## 5. Date/time formats

- All request dates: `"YYYY-MM-DD"` (e.g. `"2024-04-01"`)
- All response timestamps: ISO-8601 `"2024-04-01T00:00:00Z"`, or `null` if not set
- `activityTime`: `"HH:MM"` 24-hour

---

## 6. Full machine-readable spec

See `swagger.yaml` (OpenAPI 3.0.3) in the repo root — import into Postman/Swagger UI for
request builders and example payloads for every route.

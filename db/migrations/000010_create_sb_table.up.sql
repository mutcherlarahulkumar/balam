CREATE TABLE IF NOT EXISTS "sb" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "sb_duedate"    TIMESTAMP,
    "sb_amount"     DOUBLE PRECISION,
    "i_no"          INTEGER,
    "sb_paydate"    TIMESTAMP,
    "ch_no"         VARCHAR(10),
    "details"       TEXT
);

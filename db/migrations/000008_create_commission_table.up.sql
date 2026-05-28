CREATE TABLE IF NOT EXISTS "commission" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "bill_date"     TIMESTAMP,
    "first_comm"    DECIMAL(15, 2),
    "second_comm"   DECIMAL(15, 2),
    "third_comm"    DECIMAL(15, 2),
    "bonus_comm"    DECIMAL(15, 2),
    "single_comm"   DECIMAL(15, 2),
    "sub_comm"      DECIMAL(15, 2),
    "paydate"       TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "multipledue" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "duedate"       TIMESTAMP,
    "inst_no"       INTEGER,
    "int_amt"       DECIMAL(15, 2),
    "valid_upto"    TIMESTAMP,
    "total_amt"     DECIMAL(15, 2),

    CONSTRAINT "fk_multipledue_policy"
        FOREIGN KEY ("policy_no") REFERENCES "policies" ("policy_no")
);

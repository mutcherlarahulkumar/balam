CREATE TABLE IF NOT EXISTS "loan" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "ldate"         TIMESTAMP,
    "loan"          INTEGER,
    "intduedate"    TIMESTAMP,
    "loaninterest"  INTEGER,
    "xcharge"       INTEGER,
    "revivalded"    INTEGER,
    "details"       TEXT,

    CONSTRAINT "fk_loan_policy"
        FOREIGN KEY ("policy_no") REFERENCES "policies" ("policy_no")
);

CREATE TABLE IF NOT EXISTS "nbupdate" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "name"          VARCHAR(100),
    "dateupdated"   TIMESTAMP,
    "policyid"      INTEGER
);

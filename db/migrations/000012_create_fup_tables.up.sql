CREATE TABLE IF NOT EXISTS "fuphistory" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "oldfup"        TIMESTAMP,
    "newfup"        TIMESTAMP,
    "name"          VARCHAR(100),
    "dateupdated"   TIMESTAMP,

    CONSTRAINT "fk_fuphistory_policy"
        FOREIGN KEY ("policy_no") REFERENCES "policies" ("policy_no")
);

CREATE TABLE IF NOT EXISTS "fupupdate" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "oldfup"        TIMESTAMP,
    "newfup"        TIMESTAMP,
    "name"          VARCHAR(100),
    "dateupdated"   TIMESTAMP,

    CONSTRAINT "fk_fupupdate_policy"
        FOREIGN KEY ("policy_no") REFERENCES "policies" ("policy_no")
);

CREATE TABLE IF NOT EXISTS "fuphistory" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "oldfup"        TIMESTAMP,
    "newfup"        TIMESTAMP,
    "name"          VARCHAR(100),
    "dateupdated"   TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "fupupdate" (
    "_id"           SERIAL PRIMARY KEY,
    "policy_no"     INTEGER,
    "oldfup"        TIMESTAMP,
    "newfup"        TIMESTAMP,
    "name"          VARCHAR(100),
    "dateupdated"   TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "bankdetails" (
    "_id"               SERIAL PRIMARY KEY,
    "clientid"          INTEGER,
    "bankname"          VARCHAR(200),
    "accountnumber"     VARCHAR(50),
    "micrcode"          VARCHAR(30),
    "ifsecode"          VARCHAR(20),
    "bankcode"          VARCHAR(50),
    "branchcode"        VARCHAR(20),
    "familycode"        VARCHAR(20),
    "perscode"          VARCHAR(20),
    "policy_no"         VARCHAR(20),
    "aadhar"            VARCHAR(15),
    "pan"               VARCHAR(10),
    "ckyc"              VARCHAR(50),

    CONSTRAINT "fk_bankdetails_client"
        FOREIGN KEY ("clientid") REFERENCES "client" ("_id")
        ON DELETE CASCADE
);

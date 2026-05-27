CREATE TABLE IF NOT EXISTS "client" (
    "_id"           INTEGER PRIMARY KEY,
    "familycode"    VARCHAR(15) NOT NULL,
    "perscode"      VARCHAR(10),
    "name"          VARCHAR(80),
    "mobileno"      VARCHAR(15),
    "dob_r"         TIMESTAMP,
    "dob_a"         TIMESTAMP,
    "sex"           VARCHAR(2),
    "address"       VARCHAR(255),
    "email"         VARCHAR(80),
    "mar_date"      TIMESTAMP,
    "street"        VARCHAR(50),
    "city"          VARCHAR(50),
    "state"         VARCHAR(50),
    "designation"   VARCHAR(10),
    "occupation"    VARCHAR(50),
    "height"        VARCHAR(15),
    "weight"        VARCHAR(15),
    "name_updated"  TIMESTAMP,
    "age"           INTEGER DEFAULT 0,
    "client_type"   VARCHAR(2) DEFAULT 'N',
    "comment"       TEXT,
    "source"        VARCHAR(30),
    "lead_type"     VARCHAR(30),
    "pstatus"       VARCHAR(30),
    "custid"        VARCHAR(50),

    CONSTRAINT "fk_client_family"
        FOREIGN KEY ("familycode") REFERENCES "family" ("familycode")
        ON DELETE CASCADE
);

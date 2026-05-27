CREATE TABLE IF NOT EXISTS "family" (
    "_id"          INTEGER PRIMARY KEY,
    "familycode"   VARCHAR(15) UNIQUE NOT NULL,
    "address"      VARCHAR(250),
    "email"        VARCHAR(80),
    "mobileno"     VARCHAR(15),
    "headname"     VARCHAR(80),
    "pincode"      VARCHAR(10),
    "religion"     VARCHAR(20),
    "designation"  VARCHAR(80),
    "pin"          VARCHAR(5),
    "last_update"  TIMESTAMP,
    "last_upload"  TIMESTAMP
);

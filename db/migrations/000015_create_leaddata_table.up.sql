CREATE TABLE IF NOT EXISTS "leaddata" (
    "_id"           SERIAL PRIMARY KEY,
    "lead_id"       VARCHAR(20),
    "search_term"   VARCHAR(100),
    "name"          VARCHAR(150),
    "mobile"        VARCHAR(50),
    "address"       TEXT,
    "status"        INTEGER,
    "date_added"    TIMESTAMP
);

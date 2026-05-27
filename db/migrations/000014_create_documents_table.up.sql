CREATE TABLE IF NOT EXISTS "documents" (
    "_id"           SERIAL PRIMARY KEY,
    "clientid"      INTEGER,
    "policy_no"     INTEGER,
    "title"         VARCHAR(200),
    "photo"         TEXT,
    "profile"       BOOLEAN,

    CONSTRAINT "fk_documents_client"
        FOREIGN KEY ("clientid") REFERENCES "client" ("_id")
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "prospect_activity" (
    "_id"               SERIAL PRIMARY KEY,
    "clientid"          INTEGER,
    "activity_date"     TIMESTAMP,
    "activity_time"     TIMESTAMP,
    "activity_type"     VARCHAR(25),
    "specs"             VARCHAR(30),
    "details"           TEXT,
    "reminder"          VARCHAR(30),
    "reminder_date"     TIMESTAMP,
    "reminder_time"     TIMESTAMP,
    "status"            VARCHAR(30),
    "eventuri"          VARCHAR(200),
    "reminderuri"       VARCHAR(200),

    CONSTRAINT "fk_prospect_activity_client"
        FOREIGN KEY ("clientid") REFERENCES "client" ("_id")
        ON DELETE CASCADE
);

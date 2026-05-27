CREATE TABLE IF NOT EXISTS "msg_category" (
    "_id"       INTEGER PRIMARY KEY,
    "name"      VARCHAR(40),
    "status"    SMALLINT
);

CREATE TABLE IF NOT EXISTS "messages" (
    "_id"           INTEGER PRIMARY KEY,
    "category_id"   INTEGER,
    "message"       TEXT,
    "status"        SMALLINT,
    "msg_type"      SMALLINT,
    "defaultsms"    SMALLINT DEFAULT 0,

    CONSTRAINT "fk_messages_category"
        FOREIGN KEY ("category_id") REFERENCES "msg_category" ("_id")
);

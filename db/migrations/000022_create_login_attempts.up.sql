CREATE TABLE IF NOT EXISTS "login_attempts" (
    "_id"          SERIAL PRIMARY KEY,
    "identifier"   VARCHAR(120) NOT NULL,
    "attempted_at" TIMESTAMP    NOT NULL DEFAULT NOW(),
    "succeeded"    BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_login_attempts_identifier ON "login_attempts" ("identifier");
CREATE INDEX idx_login_attempts_attempted_at ON "login_attempts" ("attempted_at");

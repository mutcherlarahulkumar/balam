CREATE TABLE IF NOT EXISTS "plans" (
    "plan_no"       VARCHAR(10),
    "plan_name"     VARCHAR(200),
    "plan_type"     VARCHAR(80),
    "term_ppt"      VARCHAR(10),
    "sb_year"       VARCHAR(30),
    "sb_benefit"    VARCHAR(30),
    "stax"          VARCHAR(15),
    "lapsdays"      INTEGER DEFAULT 180
);

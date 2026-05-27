CREATE TABLE IF NOT EXISTS "rpt_cashflow" (
    "familycode"    VARCHAR(20),
    "perscode"      VARCHAR(20),
    "policy_no"     VARCHAR(20),
    "due_date"      TIMESTAMP,
    "name"          VARCHAR(80),
    "age"           INTEGER,
    "sb_amount"     DECIMAL(15, 2),
    "bonus"         DECIMAL(15, 2),
    "total"         DECIMAL(15, 2),
    "type"          VARCHAR(20),
    "plan"          INTEGER,
    "term"          INTEGER,
    "ppt"           INTEGER
);

CREATE TABLE IF NOT EXISTS "rpt_cashinout" (
    "year"          INTEGER,
    "cashin_mat"    DECIMAL(15, 2),
    "cashin_sb"     DECIMAL(15, 2),
    "cashout"       DECIMAL(15, 2),
    "netamt"        DECIMAL(15, 2)
);

CREATE TABLE IF NOT EXISTS "rpt_currstatus" (
    "familycode"        VARCHAR(20),
    "perscode"          VARCHAR(20),
    "policy_no"         VARCHAR(20),
    "due_status"        VARCHAR(20),
    "premium_paid"      DECIMAL(15, 2),
    "vested_bonus"      DECIMAL(15, 2),
    "surr_value"        DECIMAL(15, 2),
    "loan_taken"        DECIMAL(15, 2),
    "loan_date"         VARCHAR(50),
    "fuli_date"         VARCHAR(30),
    "loan_amt"          DECIMAL(15, 2)
);

CREATE TABLE IF NOT EXISTS "rpt_prcalendar" (
    "familycode"    VARCHAR(20),
    "perscode"      VARCHAR(20),
    "policy_no"     VARCHAR(20),
    "jan"           DECIMAL(15, 2),
    "feb"           DECIMAL(15, 2),
    "mar"           DECIMAL(15, 2),
    "apr"           DECIMAL(15, 2),
    "may"           DECIMAL(15, 2),
    "jun"           DECIMAL(15, 2),
    "jul"           DECIMAL(15, 2),
    "aug"           DECIMAL(15, 2),
    "sep"           DECIMAL(15, 2),
    "oct"           DECIMAL(15, 2),
    "nov"           DECIMAL(15, 2),
    "dec"           DECIMAL(15, 2)
);

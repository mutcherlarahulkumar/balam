-- policies
CREATE INDEX idx_policies_policy_no    ON "policies" ("policy_no");
CREATE INDEX idx_policies_familycode   ON "policies" ("familycode");
CREATE INDEX idx_policies_perscode     ON "policies" ("perscode");
CREATE INDEX idx_policies_next_premium ON "policies" ("next_premium");
CREATE INDEX idx_policies_mat_date     ON "policies" ("mat_date");
CREATE INDEX idx_policies_status       ON "policies" ("status");

-- client
CREATE INDEX idx_client_familycode     ON "client" ("familycode");
CREATE INDEX idx_client_perscode       ON "client" ("perscode");
CREATE INDEX idx_client_mobileno       ON "client" ("mobileno");

-- family
CREATE INDEX idx_family_familycode     ON "family" ("familycode");

-- bankdetails
CREATE INDEX idx_bankdetails_clientid  ON "bankdetails" ("clientid");

-- financial tables
CREATE INDEX idx_commission_policy_no  ON "commission" ("policy_no");
CREATE INDEX idx_loan_policy_no        ON "loan" ("policy_no");
CREATE INDEX idx_sb_policy_no          ON "sb" ("policy_no");
CREATE INDEX idx_sb_duedate            ON "sb" ("sb_duedate");
CREATE INDEX idx_multipledue_policy_no ON "multipledue" ("policy_no");

-- fup
CREATE INDEX idx_fuphistory_policy_no  ON "fuphistory" ("policy_no");
CREATE INDEX idx_fupupdate_policy_no   ON "fupupdate" ("policy_no");

-- crm
CREATE INDEX idx_prospect_clientid     ON "prospect_activity" ("clientid");
CREATE INDEX idx_leaddata_mobile       ON "leaddata" ("mobile");

-- agent
CREATE INDEX idx_agmast_agcode         ON "agmast" ("agcode");

-- settings
CREATE INDEX idx_settings_key          ON "settings" ("key");

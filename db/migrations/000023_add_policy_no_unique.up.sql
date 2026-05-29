-- policies.policy_no must be unique: it is the LIC-issued policy number and
-- is referenced as a business key from commission, loan, sb, multipledue, and fup tables.
ALTER TABLE "policies" ADD CONSTRAINT "policies_policy_no_unique" UNIQUE ("policy_no");

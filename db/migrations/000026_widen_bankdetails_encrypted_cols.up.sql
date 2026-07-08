-- aadhar, pan, and ckyc are stored AES-encrypted (base64), which far exceeds
-- their original plaintext-sized VARCHAR limits. Widen to TEXT.
ALTER TABLE bankdetails ALTER COLUMN aadhar TYPE TEXT;
ALTER TABLE bankdetails ALTER COLUMN pan   TYPE TEXT;
ALTER TABLE bankdetails ALTER COLUMN ckyc  TYPE TEXT;

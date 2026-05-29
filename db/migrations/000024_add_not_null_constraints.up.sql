-- Backfill NULLs before adding NOT NULL constraints
UPDATE agent SET
    name        = COALESCE(name, ''),
    address     = COALESCE(address, ''),
    mobile      = COALESCE(mobile, ''),
    email       = COALESCE(email, ''),
    login       = COALESCE(login, ''),
    password    = COALESCE(password, ''),
    branch      = COALESCE(branch, ''),
    club        = COALESCE(club, ''),
    licence_no  = COALESCE(licence_no, ''),
    pan         = COALESCE(pan, ''),
    photo       = COALESCE(photo, ''),
    slogan      = COALESCE(slogan, '');

ALTER TABLE agent
    ALTER COLUMN name       SET NOT NULL,
    ALTER COLUMN email      SET NOT NULL,
    ALTER COLUMN login      SET NOT NULL,
    ALTER COLUMN password   SET NOT NULL,
    ALTER COLUMN mobile     SET NOT NULL;

UPDATE client SET
    familycode  = COALESCE(familycode, ''),
    perscode    = COALESCE(perscode, ''),
    name        = COALESCE(name, ''),
    mobileno    = COALESCE(mobileno, ''),
    sex         = COALESCE(sex, ''),
    address     = COALESCE(address, ''),
    email       = COALESCE(email, ''),
    occupation  = COALESCE(occupation, ''),
    client_type = COALESCE(client_type, ''),
    comment     = COALESCE(comment, ''),
    source      = COALESCE(source, ''),
    city        = COALESCE(city, ''),
    state       = COALESCE(state, ''),
    age         = COALESCE(age, 0);

ALTER TABLE client
    ALTER COLUMN familycode  SET NOT NULL,
    ALTER COLUMN perscode    SET NOT NULL,
    ALTER COLUMN name        SET NOT NULL;

UPDATE family SET
    familycode  = COALESCE(familycode, ''),
    headname    = COALESCE(headname, ''),
    address     = COALESCE(address, ''),
    email       = COALESCE(email, ''),
    mobileno    = COALESCE(mobileno, ''),
    pincode     = COALESCE(pincode, ''),
    religion    = COALESCE(religion, ''),
    designation = COALESCE(designation, '');

ALTER TABLE family
    ALTER COLUMN familycode SET NOT NULL,
    ALTER COLUMN headname   SET NOT NULL;

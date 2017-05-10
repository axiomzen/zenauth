
-- EXTENSIONS
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--
-- TRIGGERS
CREATE OR REPLACE FUNCTION update_row_modified_function_()
RETURNS TRIGGER 
AS 
$$
BEGIN
    -- ASSUMES the table has a column named exactly "updated_at".
    -- Fetch date-time of actual current moment from clock, rather than start of statement or start of transaction.
    NEW.updated_at = now();
    -- NEW.num_revisions = NEW.num_revisions + 1; 
    RETURN NEW;
END;
$$ 
language 'plpgsql';
--
-- USERS TABLE
CREATE TABLE users (
  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  first_name   VARCHAR(256),
  last_name    VARCHAR(256),
  email        VARCHAR(256),
  verified     BOOLEAN NOT NULL DEFAULT false,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  preferences  JSON,
  -- these will be jwt, so there is no size limit really
  reset_token  TEXT,
  -- reset_token_expiry   TIMESTAMP WITH TIME ZONE,
  hash         VARCHAR(256),
  
  
  -- social stuff
  social_login BOOLEAN NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX users_email_idx ON users (lower(email) varchar_pattern_ops);
-- optional indexes
--CREATE INDEX users_hash_idx ON users(hash);
--CREATE INDEX users_created_at_idx ON users(created_at DESC);
--CREATE INDEX users_updated_at_idx ON users(updated_at DESC);
CREATE TRIGGER row_mod_on_users_trigger_
BEFORE UPDATE
ON users 
FOR EACH ROW 
EXECUTE PROCEDURE update_row_modified_function_();
--
-- ADMINS TABLE
CREATE TABLE admins (
  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  first_name   VARCHAR(256),
  last_name    VARCHAR(256),
  email        VARCHAR(256) UNIQUE,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  preferences  JSON,
  --reset_token  VARCHAR(256),
  --reset_token_expiry   TIMESTAMP WITH TIME ZONE,
  -- TODO: see what the exact length is that we need for hash
  hash         VARCHAR(256)
);

CREATE UNIQUE INDEX admins_email_idx ON admins (lower(email) varchar_pattern_ops);
-- optional indexes
-- since the id is already unique, we don't need an index on the hash
--CREATE INDEX admins_hash_idx ON admins(hash);
--CREATE INDEX admins_created_at_idx ON admins(created_at DESC);
--CREATE INDEX admins_updated_at_idx ON admins(updated_at DESC);
CREATE TRIGGER row_mod_on_admins_trigger_
BEFORE UPDATE
ON admins 
FOR EACH ROW 
EXECUTE PROCEDURE update_row_modified_function_();
--

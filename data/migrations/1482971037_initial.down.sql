
-- USERS TABLE
DROP INDEX IF EXISTS users_email_idx;
-- optional indexes
--DROP INDEX IF EXISTS users_hash_idx;
--DROP INDEX IF EXISTS users_created_at_idx;
--DROP INDEX IF EXISTS users_updated_at_idx;
DROP TRIGGER IF EXISTS row_mod_on_users_trigger_ ON users;
DROP TABLE IF EXISTS users CASCADE;
--
-- ADMINS TABLE
DROP INDEX IF EXISTS admins_email_idx;
-- optional indexes
--DROP INDEX IF EXISTS admins_hash_idx;
--DROP INDEX IF EXISTS admins_created_at_idx;
--DROP INDEX IF EXISTS admins_updated_at_idx;
DROP TRIGGER IF EXISTS row_mod_on_admins_trigger_ ON admins;
DROP TABLE IF EXISTS admins CASCADE;
--
-- TRIGGERS
DROP FUNCTION IF EXISTS update_row_modified_function_();
--
-- EXTENSIONS
DROP EXTENSION IF EXISTS "uuid-ossp";

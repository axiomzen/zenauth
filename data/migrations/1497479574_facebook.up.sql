ALTER TABLE users 
ADD COLUMN facebook_id TEXT UNIQUE,
ADD COLUMN facebook_username TEXT,
ADD COLUMN facebook_token TEXT,
ADD COLUMN facebook_email VARCHAR(256);

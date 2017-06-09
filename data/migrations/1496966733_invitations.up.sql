CREATE TABLE invitations (
  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email        VARCHAR(256) NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX invitations_email_idx ON invitations (lower(email) varchar_pattern_ops);

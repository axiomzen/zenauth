CREATE TYPE code_type AS ENUM ('email', 'facebook', 'url');

CREATE TABLE invitations (
  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  type         code_type NOT NULL,
  code         VARCHAR(256) NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX invitation_code_idx ON invitations (type, code);

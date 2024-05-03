CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
    expiry timestamp(0) with time zone NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
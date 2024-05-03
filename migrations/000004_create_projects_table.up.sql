CREATE TABLE IF NOT EXISTS projects (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    detail text NOT NULL,
    version integer NOT NULL DEFAULT 1,
    user_id bigserial NOT NULL REFERENCES users ON DELETE RESTRICT
);
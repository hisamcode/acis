
CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    nominal numeric NOT NULL,
    detail text NOT NULL,
    emoji_id text NOT NULL,
    version integer NOT NULL DEFAULT 1,
    wts_id serial NOT NULL REFERENCES wts ON DELETE RESTRICT,
    project_id bigserial NOT NULL REFERENCES projects ON DELETE RESTRICT,
    created_by bigserial NOT NULL REFERENCES users ON DELETE RESTRICT
)

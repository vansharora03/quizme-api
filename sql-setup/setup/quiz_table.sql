CREATE TABLE IF NOT EXISTS quiz (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,
    title text NOT NULL
);


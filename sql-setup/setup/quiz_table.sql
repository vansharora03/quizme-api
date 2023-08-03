CREATE TABLE IF NOT EXISTS quiz (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,
    user_id bigint NOT NULL REFERENCES user_account(id) ON DELETE CASCADE,
    title text NOT NULL
);


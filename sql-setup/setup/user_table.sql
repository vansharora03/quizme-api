CREATE TABLE IF NOT EXISTS user_account (
    id bigserial PRIMARY KEY,
    username text NOT NULL,
    hashed_password bytea NOT NULL,
    email citext UNIQUE NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
    

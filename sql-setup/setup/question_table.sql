CREATE TABLE IF NOT EXISTS question (
    id bigserial PRIMARY KEY,
    quiz_id bigint NOT NULL REFERENCES quiz(id) ON DELETE CASCADE,
    prompt text NOT NULL,
    choices text[] NOT NULL,
    correct_index integer NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

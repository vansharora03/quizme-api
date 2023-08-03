CREATE TABLE IF NOT EXISTS score (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES user_account(id) ON DELETE CASCADE,
    quiz_id bigint NOT NULL REFERENCES quiz(id) ON DELETE CASCADE,
    chosen_choices_indices integer[] NOT NULL,
    chosen_choices_correctness boolean[] NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

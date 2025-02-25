CREATE TABLE talk_session_generated_images (
    talk_session_id UUID PRIMARY KEY,
    wordmap_url TEXT NOT NULL,
    tsnc_url TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT (now()),
    updated_at timestamp NOT NULL DEFAULT (now())
);

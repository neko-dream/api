CREATE TABLE IF NOT EXISTS user_group_dimension_reduction_info (
    talk_session_id UUID NOT NULL REFERENCES talk_sessions(talk_session_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dimension_reduction_type TEXT NOT NULL,
    summary TEXT NOT NULL,
    x_desc TEXT NOT NULL,
    y_desc TEXT NOT NULL,
    PRIMARY KEY (talk_session_id)
);

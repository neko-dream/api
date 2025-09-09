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

CREATE TRIGGER update_user_group_dimension_reduction_info_updated_at
    BEFORE UPDATE ON user_group_dimension_reduction_info
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

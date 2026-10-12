CREATE TABLE report_feedback (
    report_feedback_id uuid PRIMARY KEY,
    talk_session_report_history_id uuid NOT NULL,
    user_id uuid NOT NULL,
    -- good/bad
    feedback_type INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(talk_session_report_history_id, user_id)
);

CREATE INDEX idx_report_feedback_user_id_talk_session_report_history_id  ON report_feedback(user_id, talk_session_report_history_id);

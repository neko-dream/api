CREATE TABLE talk_session_report_histories (
    talk_session_report_history_id uuid PRIMARY KEY,
    talk_session_id uuid NOT NULL,
    report text NOT NULL,
    created_at timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE talk_session_report {
    talk_session_report_id uuid PRIMARY KEY,
    talk_session_id uuid NOT NULL,
    report text NOT NULL,
    created_at timestamp NOT NULL DEFAULT (now()),
    updated_at timestamp NOT NULL DEFAULT (now())
};

CREATE UNIQUE INDEX talk_session_report_talk_session_id_index ON talk_session_report (talk_session_id);

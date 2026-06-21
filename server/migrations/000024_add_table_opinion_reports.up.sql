CREATE TABLE IF NOT EXISTS opinion_reports (
    opinion_report_id UUID PRIMARY KEY,
    opinion_id UUID NOT NULL,
    talk_session_id UUID NOT NULL,
    reporter_id UUID NOT NULL,
    reason INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'unconfirmed',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX opinion_reports_opinion_id_idx ON opinion_reports(opinion_id);
CREATE INDEX opinion_reports_talk_session_id_idx ON opinion_reports(talk_session_id);

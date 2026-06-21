-- セッション参加前に回答必須のアンケート（任意 1:0..1）
CREATE TABLE talk_session_surveys (
    survey_id UUID PRIMARY KEY,
    talk_session_id UUID NOT NULL UNIQUE REFERENCES talk_sessions(talk_session_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_talk_session_surveys_talk_session_id ON talk_session_surveys (talk_session_id);

CREATE TYPE talk_session_survey_question_type AS ENUM (
    'free_text',
    'single_choice',
    'multi_choice',
    'rating',
    'date',
    'dropdown'
);

CREATE TABLE talk_session_survey_questions (
    question_id UUID PRIMARY KEY,
    survey_id UUID NOT NULL REFERENCES talk_session_surveys(survey_id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type talk_session_survey_question_type NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT FALSE,
    allow_other BOOLEAN NOT NULL DEFAULT FALSE,
    max_length INT,
    min_value INT,
    max_value INT,
    min_label VARCHAR(255),
    max_label VARCHAR(255),
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_talk_session_survey_questions_survey_id ON talk_session_survey_questions (survey_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_talk_session_survey_questions_order ON talk_session_survey_questions (survey_id, display_order) WHERE deleted_at IS NULL;

CREATE TABLE talk_session_survey_choices (
    choice_id UUID PRIMARY KEY,
    question_id UUID NOT NULL REFERENCES talk_session_survey_questions(question_id) ON DELETE CASCADE,
    choice_text VARCHAR(500) NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_talk_session_survey_choices_question_id ON talk_session_survey_choices (question_id) WHERE deleted_at IS NULL;

CREATE TABLE talk_session_survey_submissions (
    submission_id UUID PRIMARY KEY,
    survey_id UUID NOT NULL REFERENCES talk_session_surveys(survey_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (survey_id, user_id)
);
CREATE INDEX idx_talk_session_survey_submissions_user_id ON talk_session_survey_submissions (user_id);

CREATE TABLE talk_session_survey_answers (
    answer_id UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES talk_session_survey_submissions(submission_id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES talk_session_survey_questions(question_id),
    answer_value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_talk_session_survey_answers_submission_id ON talk_session_survey_answers (submission_id);
CREATE INDEX idx_talk_session_survey_answers_question_id ON talk_session_survey_answers (question_id);

-- enum を旧 6 種に戻す。
ALTER TYPE talk_session_survey_question_type RENAME TO talk_session_survey_question_type_old;
CREATE TYPE talk_session_survey_question_type AS ENUM (
    'free_text',
    'single_choice',
    'multi_choice',
    'rating',
    'date',
    'dropdown'
);
ALTER TABLE talk_session_survey_questions
    ALTER COLUMN question_type TYPE talk_session_survey_question_type
    USING question_type::text::talk_session_survey_question_type;
DROP TYPE talk_session_survey_question_type_old;

-- 削除した列を復元する（データは復元できない）。
ALTER TABLE talk_session_survey_questions
    ADD COLUMN IF NOT EXISTS min_value INT,
    ADD COLUMN IF NOT EXISTS max_value INT,
    ADD COLUMN IF NOT EXISTS min_label VARCHAR(255),
    ADD COLUMN IF NOT EXISTS max_label VARCHAR(255),
    ADD COLUMN IF NOT EXISTS contains_pii BOOLEAN NOT NULL DEFAULT FALSE;

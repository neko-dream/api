ALTER TABLE talk_session_survey_answers DROP COLUMN IF EXISTS is_encrypted;
ALTER TABLE talk_session_survey_questions DROP COLUMN IF EXISTS contains_pii;

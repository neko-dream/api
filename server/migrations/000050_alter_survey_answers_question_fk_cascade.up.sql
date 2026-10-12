-- talk_session_survey_answers.question_id の FK を ON DELETE CASCADE に変更する。
-- これにより、回答ゼロの質問を hard delete した際に CASCADE で answers も消える。
ALTER TABLE talk_session_survey_answers
    DROP CONSTRAINT talk_session_survey_answers_question_id_fkey;

ALTER TABLE talk_session_survey_answers
    ADD CONSTRAINT talk_session_survey_answers_question_id_fkey
    FOREIGN KEY (question_id) REFERENCES talk_session_survey_questions(question_id) ON DELETE CASCADE;

-- 質問単位で PII（個人情報）フラグを持たせる。
-- contains_pii=true の質問への回答は text/other/date を AES-GCM で暗号化して保存する。
ALTER TABLE talk_session_survey_questions
    ADD COLUMN IF NOT EXISTS contains_pii BOOLEAN NOT NULL DEFAULT FALSE;

-- 回答行単位で「保存時に暗号化されたかどうか」を持つ。
-- 質問の contains_pii を後から変更しても、行ごとに正しく復号できるようにする。
ALTER TABLE talk_session_survey_answers
    ADD COLUMN IF NOT EXISTS is_encrypted BOOLEAN NOT NULL DEFAULT FALSE;

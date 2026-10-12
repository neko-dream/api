-- アンケートを選択肢系（single_choice / dropdown）のみに絞り込み、回答を全件暗号化する方針へ変更する。
-- 旧マイグレーション 49/51 は既存環境で適用済みのため書き換えず、本マイグレーションで前進的に変換する。

-- 1. 廃止する質問タイプ（free_text / multi_choice / rating / date）の質問を削除する。
--    answers への FK は ON DELETE CASCADE（migration 50）なので、紐づく回答も連鎖削除される。
DELETE FROM talk_session_survey_questions
WHERE question_type::text IN ('free_text', 'multi_choice', 'rating', 'date');

-- 2. rating 専用の設定列と PII フラグ列を削除する（max_length は「その他」入力用に残す）。
ALTER TABLE talk_session_survey_questions
    DROP COLUMN IF EXISTS min_value,
    DROP COLUMN IF EXISTS max_value,
    DROP COLUMN IF EXISTS min_label,
    DROP COLUMN IF EXISTS max_label,
    DROP COLUMN IF EXISTS contains_pii;

-- 3. enum を single_choice / dropdown のみに作り直す。
ALTER TYPE talk_session_survey_question_type RENAME TO talk_session_survey_question_type_old;
CREATE TYPE talk_session_survey_question_type AS ENUM (
    'single_choice',
    'dropdown'
);
ALTER TABLE talk_session_survey_questions
    ALTER COLUMN question_type TYPE talk_session_survey_question_type
    USING question_type::text::talk_session_survey_question_type;
DROP TYPE talk_session_survey_question_type_old;

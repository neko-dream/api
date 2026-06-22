-- name: CreateTalkSessionSurvey :exec
INSERT INTO talk_session_surveys (survey_id, talk_session_id, created_at, updated_at)
VALUES ($1, $2, $3, $4);

-- name: GetSurveyByID :one
SELECT survey_id, talk_session_id, created_at, updated_at
FROM talk_session_surveys
WHERE survey_id = $1;

-- name: GetSurveyByTalkSessionID :one
SELECT survey_id, talk_session_id, created_at, updated_at
FROM talk_session_surveys
WHERE talk_session_id = $1;

-- name: TouchSurvey :exec
UPDATE talk_session_surveys SET updated_at = $2 WHERE survey_id = $1;

-- name: ListSurveyQuestions :many
SELECT
    question_id, survey_id, question_text, question_type::text AS question_type,
    is_required, allow_other, max_length, display_order, deleted_at, created_at, updated_at
FROM talk_session_survey_questions
WHERE survey_id = $1
ORDER BY display_order ASC, created_at ASC;

-- name: CreateSurveyQuestion :exec
INSERT INTO talk_session_survey_questions (
    question_id, survey_id, question_text, question_type,
    is_required, allow_other, max_length, display_order, created_at, updated_at
) VALUES (
    sqlc.arg('question_id')::uuid,
    sqlc.arg('survey_id')::uuid,
    sqlc.arg('question_text')::text,
    sqlc.arg('question_type')::talk_session_survey_question_type,
    sqlc.arg('is_required')::boolean,
    sqlc.arg('allow_other')::boolean,
    sqlc.narg('max_length')::int,
    sqlc.arg('display_order')::int,
    sqlc.arg('created_at')::timestamptz,
    sqlc.arg('updated_at')::timestamptz
);

-- name: UpdateSurveyQuestion :exec
UPDATE talk_session_survey_questions
SET question_text = sqlc.arg('question_text')::text,
    question_type = sqlc.arg('question_type')::talk_session_survey_question_type,
    is_required = sqlc.arg('is_required')::boolean,
    allow_other = sqlc.arg('allow_other')::boolean,
    max_length = sqlc.narg('max_length')::int,
    display_order = sqlc.arg('display_order')::int,
    updated_at = sqlc.arg('updated_at')::timestamptz
WHERE question_id = sqlc.arg('question_id')::uuid;

-- name: SoftDeleteSurveyQuestion :exec
UPDATE talk_session_survey_questions
SET deleted_at = $2, updated_at = $2
WHERE question_id = $1 AND deleted_at IS NULL;

-- name: HardDeleteSurveyQuestion :exec
DELETE FROM talk_session_survey_questions WHERE question_id = $1;

-- name: ListSurveyChoicesByQuestionIDs :many
SELECT choice_id, question_id, choice_text, display_order, deleted_at, created_at, updated_at
FROM talk_session_survey_choices
WHERE question_id = ANY(sqlc.arg('question_ids')::uuid[])
ORDER BY display_order ASC, created_at ASC;

-- name: CreateSurveyChoice :exec
INSERT INTO talk_session_survey_choices (
    choice_id, question_id, choice_text, display_order, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateSurveyChoice :exec
UPDATE talk_session_survey_choices
SET choice_text = $2,
    display_order = $3,
    updated_at = $4
WHERE choice_id = $1;

-- name: SoftDeleteSurveyChoice :exec
UPDATE talk_session_survey_choices
SET deleted_at = $2, updated_at = $2
WHERE choice_id = $1 AND deleted_at IS NULL;

-- name: HardDeleteSurveyChoice :exec
DELETE FROM talk_session_survey_choices WHERE choice_id = $1;

-- name: CountAnswersForQuestion :one
SELECT COUNT(*) AS answer_count
FROM talk_session_survey_answers
WHERE question_id = $1;

-- name: CountAnswersForChoice :one
-- 回答はすべて暗号化されており選択肢 ID で直接絞り込めないため、
-- 当該選択肢が属する質問への回答件数を返す（hard delete 可否判定用の保守的なカウント）。
SELECT COUNT(*) AS answer_count
FROM talk_session_survey_answers a
INNER JOIN talk_session_survey_choices c ON c.question_id = a.question_id
WHERE c.choice_id = sqlc.arg('choice_id')::uuid;

-- name: CreateSurveySubmission :exec
INSERT INTO talk_session_survey_submissions (submission_id, survey_id, user_id, submitted_at)
VALUES ($1, $2, $3, $4);

-- name: CreateSurveyAnswer :exec
INSERT INTO talk_session_survey_answers (answer_id, submission_id, question_id, answer_value, is_encrypted, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: HasSubmittedSurveyBySessionID :one
SELECT EXISTS (
    SELECT 1
    FROM talk_session_surveys ts
    INNER JOIN talk_session_survey_submissions sub ON sub.survey_id = ts.survey_id
    WHERE ts.talk_session_id = $1 AND sub.user_id = $2
) AS has_submitted;

-- name: GetSurveySubmissionByUserID :one
SELECT submission_id, survey_id, user_id, submitted_at
FROM talk_session_survey_submissions
WHERE survey_id = $1 AND user_id = $2;

-- name: ListSurveySubmissions :many
SELECT submission_id, survey_id, user_id, submitted_at
FROM talk_session_survey_submissions
WHERE survey_id = $1
ORDER BY submitted_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAllSurveySubmissions :many
SELECT submission_id, survey_id, user_id, submitted_at
FROM talk_session_survey_submissions
WHERE survey_id = $1
ORDER BY submitted_at ASC;

-- name: CountSurveySubmissions :one
SELECT COUNT(*) AS submission_count
FROM talk_session_survey_submissions
WHERE survey_id = $1;

-- name: ListSurveyAnswersBySubmissionIDs :many
SELECT answer_id, submission_id, question_id, answer_value, is_encrypted, created_at
FROM talk_session_survey_answers
WHERE submission_id = ANY(sqlc.arg('submission_ids')::uuid[])
ORDER BY created_at ASC;

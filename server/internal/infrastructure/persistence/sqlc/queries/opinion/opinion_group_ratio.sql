-- name: GetGroupRatioByOpinionID :many
SELECT
  sqlc.embed(representative_opinions)
FROM representative_opinions
WHERE representative_opinions.opinion_id = sqlc.narg('opinion_id')::uuid
ORDER BY representative_opinions.created_at DESC;

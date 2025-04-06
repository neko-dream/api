-- name: FindOpinionsByOpinionIDs :many
SELECT
    sqlc.embed(opinions),
    sqlc.embed(users)
FROM
    opinions
LEFT JOIN users
    ON opinions.user_id = users.user_id
WHERE
    opinions.opinion_id IN(sqlc.slice('opinion_ids'))
;

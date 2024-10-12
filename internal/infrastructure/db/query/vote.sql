-- name: CreateVote :exec
INSERT INTO votes (
    vote_id,
    opinion_id,
    user_id,
    vote_type,
    created_at
) VALUES ($1, $2, $3, $4, $5);

-- name: GetVoteByUserIDAndOpinionID :one
SELECT
    vote_id,
    opinion_id,
    user_id,
    vote_type,
    created_at
FROM votes
WHERE opinion_id = $1 AND user_id = $2
LIMIT 1;

-- name: CreateVote :exec
INSERT INTO votes (
    vote_id,
    opinion_id,
    user_id,
    vote_type,
    created_at
) VALUES ($1, $2, $3, $4, $5);

-- name: FindVoteByUserIDAndOpinionID :one
SELECT * FROM votes WHERE user_id = $1 AND opinion_id = $2;

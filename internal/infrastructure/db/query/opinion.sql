-- name: CreateOpinion :exec
INSERT INTO opinions (
    opinion_id,
    talk_session_id,
    user_id,
    parent_opinion_id,
    content,
    vote_id,
    created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7);


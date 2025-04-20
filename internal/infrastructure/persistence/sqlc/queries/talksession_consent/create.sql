-- name: CreateTSConsent :exec
INSERT INTO talksession_consents (
    talksession_id,
    user_id,
    consented_at,
    restrictions
) VALUES (
    $1,
    $2,
    $3,
    $4
);

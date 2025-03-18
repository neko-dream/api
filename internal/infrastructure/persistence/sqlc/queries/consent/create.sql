-- name: CreatePolicyConsent :exec
INSERT INTO policy_consents (
    policy_consent_id,
    user_id,
    policy_version,
    consented_at,
    ip_address,
    user_agent
) VALUES ($1, $2, $3, $4, $5, $6);

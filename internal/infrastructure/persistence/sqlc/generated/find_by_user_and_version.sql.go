// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: find_by_user_and_version.sql

package model

import (
	"context"

	"github.com/google/uuid"
)

const findConsentByUserAndVersion = `-- name: FindConsentByUserAndVersion :one
SELECT
    policy_consents.policy_consent_id, policy_consents.user_id, policy_consents.policy_version, policy_consents.consented_at, policy_consents.ip_address, policy_consents.user_agent, policy_consents.index
FROM
    policy_consents
WHERE
    user_id = $1
    AND policy_version = $2
`

type FindConsentByUserAndVersionParams struct {
	UserID        uuid.UUID
	PolicyVersion string
}

type FindConsentByUserAndVersionRow struct {
	PolicyConsent PolicyConsent
}

// FindConsentByUserAndVersion
//
//	SELECT
//	    policy_consents.policy_consent_id, policy_consents.user_id, policy_consents.policy_version, policy_consents.consented_at, policy_consents.ip_address, policy_consents.user_agent, policy_consents.index
//	FROM
//	    policy_consents
//	WHERE
//	    user_id = $1
//	    AND policy_version = $2
func (q *Queries) FindConsentByUserAndVersion(ctx context.Context, arg FindConsentByUserAndVersionParams) (FindConsentByUserAndVersionRow, error) {
	row := q.db.QueryRowContext(ctx, findConsentByUserAndVersion, arg.UserID, arg.PolicyVersion)
	var i FindConsentByUserAndVersionRow
	err := row.Scan(
		&i.PolicyConsent.PolicyConsentID,
		&i.PolicyConsent.UserID,
		&i.PolicyConsent.PolicyVersion,
		&i.PolicyConsent.ConsentedAt,
		&i.PolicyConsent.IpAddress,
		&i.PolicyConsent.UserAgent,
		&i.PolicyConsent.Index,
	)
	return i, err
}

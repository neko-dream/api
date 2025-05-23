// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: find_by_talksession_id_and_user_id.sql

package model

import (
	"context"

	"github.com/google/uuid"
)

const findTSConsentByTalksessionIdAndUserId = `-- name: FindTSConsentByTalksessionIdAndUserId :one
SELECT
    talksession_consents.talksession_id, talksession_consents.user_id, talksession_consents.restrictions, talksession_consents.consented_at
FROM
    talksession_consents
WHERE
    talksession_id = $1
    AND user_id = $2
`

type FindTSConsentByTalksessionIdAndUserIdParams struct {
	TalksessionID uuid.UUID
	UserID        uuid.UUID
}

type FindTSConsentByTalksessionIdAndUserIdRow struct {
	TalksessionConsent TalksessionConsent
}

// FindTSConsentByTalksessionIdAndUserId
//
//	SELECT
//	    talksession_consents.talksession_id, talksession_consents.user_id, talksession_consents.restrictions, talksession_consents.consented_at
//	FROM
//	    talksession_consents
//	WHERE
//	    talksession_id = $1
//	    AND user_id = $2
func (q *Queries) FindTSConsentByTalksessionIdAndUserId(ctx context.Context, arg FindTSConsentByTalksessionIdAndUserIdParams) (FindTSConsentByTalksessionIdAndUserIdRow, error) {
	row := q.db.QueryRowContext(ctx, findTSConsentByTalksessionIdAndUserId, arg.TalksessionID, arg.UserID)
	var i FindTSConsentByTalksessionIdAndUserIdRow
	err := row.Scan(
		&i.TalksessionConsent.TalksessionID,
		&i.TalksessionConsent.UserID,
		&i.TalksessionConsent.Restrictions,
		&i.TalksessionConsent.ConsentedAt,
	)
	return i, err
}

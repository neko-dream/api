// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: find_by_version.sql

package model

import (
	"context"
)

const findPolicyByVersion = `-- name: FindPolicyByVersion :one
SELECT
    policy_versions.version, policy_versions.created_at, policy_versions.index
FROM
    policy_versions
WHERE
    version = $1
LIMIT 1
`

type FindPolicyByVersionRow struct {
	PolicyVersion PolicyVersion
}

// FindPolicyByVersion
//
//	SELECT
//	    policy_versions.version, policy_versions.created_at, policy_versions.index
//	FROM
//	    policy_versions
//	WHERE
//	    version = $1
//	LIMIT 1
func (q *Queries) FindPolicyByVersion(ctx context.Context, version string) (FindPolicyByVersionRow, error) {
	row := q.db.QueryRowContext(ctx, findPolicyByVersion, version)
	var i FindPolicyByVersionRow
	err := row.Scan(&i.PolicyVersion.Version, &i.PolicyVersion.CreatedAt, &i.PolicyVersion.Index)
	return i, err
}

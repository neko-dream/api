package repository

import (
	"fmt"
	legacy "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgconn"
	"testing"
)

func TestSingleMembershipConstraintErrorBothPGXVersions(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want bool
	}{
		{&legacy.PgError{Code: "23505", ConstraintName: "organization_users_one_per_user"}, true},
		{&pgconn.PgError{Code: "23505", ConstraintName: "organization_users_one_per_user"}, true},
		{fmt.Errorf("wrapped: %w", &pgconn.PgError{Code: "23505", ConstraintName: "organization_users_one_per_user"}), true},
		{&pgconn.PgError{Code: "23505", ConstraintName: "organization_users_pkey"}, false},
		{&legacy.PgError{Code: "23503", ConstraintName: "organization_users_one_per_user"}, false},
		{fmt.Errorf("unavailable"), false},
	} {
		if got := isSingleMembershipConflict(tc.err); got != tc.want {
			t.Fatalf("%T: got %v want %v", tc.err, got, tc.want)
		}
	}
}

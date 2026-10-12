package txtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
)

type TransactionalTestCase[T any] struct {
	Name    string
	SetupFn func(context.Context, *T) error
	TestFn  func(context.Context, *T) error
	WantErr bool
}

func RunTransactionalTests[T any](t *testing.T, dbManager *db.DBManager, initialData *T, testCases []*TransactionalTestCase[T]) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := dbManager.TestTx(t.Context(), func(ctx context.Context) error {
				if tc.SetupFn != nil {
					if err := tc.SetupFn(ctx, initialData); err != nil {
						return fmt.Errorf("setup error: %w", err)
					}
				}
				return tc.TestFn(ctx, initialData)
			})

			if (err != nil) != tc.WantErr {
				t.Errorf("%s error = %v, wantErr %v", tc.Name, err, tc.WantErr)
			}
		})
	}
}

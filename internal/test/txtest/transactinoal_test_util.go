package txtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type TransactionalTestCase[T any] struct {
	Name    string
	SetupFn func(*TestContext[T]) error
	TestFn  func(*TestContext[T]) error
	WantErr bool
}

func RunTransactionalTests[T any](t *testing.T, dbManager *db.DBManager, initialData T, testCases []*TransactionalTestCase[T]) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := dbManager.TestTx(context.Background(), func(ctx context.Context) error {
				testCtx := NewTestContext(ctx, initialData)
				if tc.SetupFn != nil {
					if err := tc.SetupFn(testCtx); err != nil {
						return fmt.Errorf("setup error: %w", err)
					}
				}
				return tc.TestFn(testCtx)
			})

			if (err != nil) != tc.WantErr {
				t.Errorf("%s error = %v, wantErr %v", tc.Name, err, tc.WantErr)
			}
		})
	}
}

type TestContext[T any] struct {
	context.Context
	Data T
}

func NewTestContext[T any](ctx context.Context, data T) *TestContext[T] {
	ctx, span := otel.Tracer("txtest").Start(ctx, "NewTestContext")
	defer span.End()

	return &TestContext[T]{
		Context: ctx,
		Data:    data,
	}
}

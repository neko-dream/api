package consent

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
)

type PolicyRepository interface {
	Save(ctx context.Context, policy *Policy) error
	FetchLatestPolicy(ctx context.Context) (*Policy, error)
	FindByVersion(ctx context.Context, version string) (*Policy, error)
}

type Policy struct {
	ID        shared.UUID[Policy]
	Version   string
	CreatedAt time.Time
}

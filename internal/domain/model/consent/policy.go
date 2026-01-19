package consent

//go:generate go tool mockgen -source=$GOFILE -package=mock_consent -destination=../../mock/consent/policy.go

import (
	"context"
	"time"
)

type PolicyRepository interface {
	Save(ctx context.Context, policy *Policy) error
	FetchLatestPolicy(ctx context.Context) (*Policy, error)
	FindByVersion(ctx context.Context, version string) (*Policy, error)
}

type Policy struct {
	Version   string
	CreatedAt time.Time
}

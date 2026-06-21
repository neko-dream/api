package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/neko-dream/api/internal/domain/model/consent"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type policyRepository struct {
	*db.DBManager
	*config.Config
}

func NewPolicyRepository(
	DBManager *db.DBManager,
	conf *config.Config,
) consent.PolicyRepository {
	return &policyRepository{
		DBManager: DBManager,
		Config:    conf,
	}
}

// FetchLatestPolicy implements consent.PolicyRepository.
func (p *policyRepository) FetchLatestPolicy(ctx context.Context) (*consent.Policy, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "policyRepository.FetchLatestPolicy")
	defer span.End()

	_, err := p.DBManager.GetQueries(ctx).FindPolicyByVersion(ctx, p.Config.POLICY_VERSION)
	if errors.Is(err, sql.ErrNoRows) {
		// なかったら作る
		if err := p.Save(ctx, &consent.Policy{
			Version:   p.Config.POLICY_VERSION,
			CreatedAt: time.Now(),
		}); err != nil {
			utils.HandleError(ctx, err, "ポリシーを作成できませんでした。")
			return nil, err
		}
	}

	return &consent.Policy{
		Version:   p.Config.POLICY_VERSION,
		CreatedAt: time.Now(),
	}, nil
}

// FindByVersion implements consent.PolicyRepository.
func (p *policyRepository) FindByVersion(ctx context.Context, version string) (*consent.Policy, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "policyRepository.FindByVersion")
	defer span.End()

	// resRow, err := p.DBManager.GetQueries(ctx).FindPolicyByVersion(ctx, version)
	// if err != nil {
	// 	utils.HandleError(ctx, err, "ポリシーを取得できませんでした。")
	// 	return nil, err
	// }

	return &consent.Policy{
		Version:   p.POLICY_VERSION,
		CreatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
	}, nil
}

// Save implements consent.PolicyRepository.
func (p *policyRepository) Save(ctx context.Context, policy *consent.Policy) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "policyRepository.Save")
	defer span.End()

	if err := p.DBManager.GetQueries(ctx).CreatePolicy(ctx, model.CreatePolicyParams{
		Version:   policy.Version,
		CreatedAt: policy.CreatedAt,
	}); err != nil {
		utils.HandleError(ctx, err, "ポリシーを作成できませんでした。")
		return err
	}

	return nil
}

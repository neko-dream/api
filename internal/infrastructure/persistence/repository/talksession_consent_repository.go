package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/talksession/talksession_consent"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type talkSessionConsentRepository struct {
	*db.DBManager
}

func NewTalkSessionConsentRepository(dbManager *db.DBManager) talksession_consent.TalkSessionConsentRepository {
	return &talkSessionConsentRepository{
		DBManager: dbManager,
	}
}

func (r *talkSessionConsentRepository) Store(ctx context.Context, consent talksession_consent.TalkSessionConsent) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionConsentRepository.Store")
	defer span.End()

	if err := r.GetQueries(ctx).CreateTSConsent(ctx, model.CreateTSConsentParams{
		TalksessionID: consent.TalkSessionID.UUID(),
		UserID:        consent.UserID.UUID(),
		ConsentedAt:   consent.ConsentedAt,
		Restrictions:  talksession.Restrictions(consent.Restrictions),
	}); err != nil {
		utils.HandleError(ctx, err, "Consentの保存に失敗しました。")
		return err
	}

	return nil
}

func (r *talkSessionConsentRepository) FindByTalkSessionIDAndUserID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], userID shared.UUID[user.User]) (*talksession_consent.TalkSessionConsent, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionConsentRepository.FindByUserID")
	defer span.End()
	row, err := r.GetQueries(ctx).FindTSConsentByTalksessionIdAndUserId(ctx, model.FindTSConsentByTalksessionIdAndUserIdParams{
		TalksessionID: talkSessionID.UUID(),
		UserID:        userID.UUID(),
	})
	if err != nil {
		utils.HandleError(ctx, err, "Consentの取得に失敗しました。")
		return nil, err
	}

	tc, err := talksession_consent.NewTalkSessionConsent(
		shared.UUID[talksession.TalkSession](row.TalksessionConsent.TalksessionID),
		shared.UUID[user.User](row.TalksessionConsent.UserID),
		row.TalksessionConsent.ConsentedAt,
		row.TalksessionConsent.Restrictions,
	)
	if err != nil {
		return nil, err
	}

	return &tc, nil
}

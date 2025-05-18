package talksession_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/talksession/talksession_consent"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type hasConsentQuery struct {
	talkSessionConsentService talksession_consent.TalkSessionConsentService
}

func NewHasConsentQuery(
	talkSessionConsentService talksession_consent.TalkSessionConsentService,
) talksession.HasConsentQuery {
	return &hasConsentQuery{
		talkSessionConsentService: talkSessionConsentService,
	}
}

func (q *hasConsentQuery) Execute(ctx context.Context, input talksession.HasConsentQueryInput) (bool, error) {
	ctx, span := otel.Tracer("talksession").Start(ctx, "hasConsentQuery.Execute")
	defer span.End()

	hasConsented, err := q.talkSessionConsentService.HasConsented(ctx, input.TalkSessionID, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "Consentの取得に失敗しました。")
		return false, err
	}

	return hasConsented, nil
}

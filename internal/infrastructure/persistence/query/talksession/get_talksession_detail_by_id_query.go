package talksession_query

import (
	"context"
	"log"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type getTalkSessionDetailByIDQuery struct {
	*db.DBManager
}

func NewGetTalkSessionDetailByIDQueryHandler(tm *db.DBManager) talksession.GetTalkSessionDetailByIDQuery {
	return &getTalkSessionDetailByIDQuery{
		DBManager: tm,
	}
}

// Execute トークセッションの詳細を取得する
func (h *getTalkSessionDetailByIDQuery) Execute(ctx context.Context, input talksession.GetTalkSessionDetailInput) (*talksession.GetTalkSessionDetailOutput, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "getTalkSessionDetailByIDQuery.Execute")
	defer span.End()

	talkSessionRow, err := h.GetQueries(ctx).GetTalkSessionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, messages.TalkSessionNotFound
	}
	log.Printf("talkSessionRow: %+v", talkSessionRow)

	var result dto.TalkSessionWithDetail
	if err := copier.CopyWithOption(&result, talkSessionRow, copier.Option{
		DeepCopy: true,
	}); err != nil {
		return nil, messages.TalkSessionNotFound
	}

	return &talksession.GetTalkSessionDetailOutput{
		TalkSessionWithDetail: result,
	}, nil
}

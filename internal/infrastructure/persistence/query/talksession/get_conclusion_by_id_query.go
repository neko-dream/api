package talksession_query

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"go.opentelemetry.io/otel"
)

type getConclusionByIDQuery struct {
	*db.DBManager
}

func NewGetConclusionByIDQueryHandler(tm *db.DBManager) talksession.GetConclusionByIDQuery {
	return &getConclusionByIDQuery{
		DBManager: tm,
	}
}

func (h *getConclusionByIDQuery) Execute(ctx context.Context, input talksession.GetConclusionByIDQueryRequest) (*talksession.GetConclusionByIDQueryResponse, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "getConclusionByIDQuery.Execute")
	defer span.End()

	ts, err := h.GetQueries(ctx).GetTalkSessionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	// まだ終了していないトークセッションに対しては結論を取得できない
	if ts.TalkSession.ScheduledEndTime.After(clock.Now(ctx)) {
		return nil, messages.TalkSessionNotFinished
	}

	res, err := h.GetQueries(ctx).GetConclusionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var result talksession.GetConclusionByIDQueryResponse
	if err := copier.CopyWithOption(&result, res, copier.Option{
		DeepCopy: true,
	}); err != nil {
		return nil, err
	}

	return &result, nil
}

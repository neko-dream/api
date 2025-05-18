package talksession_query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type BrowseJoinedTalkSessionQueryHandler struct {
	*db.DBManager
}

func NewBrowseJoinedTalkSessionQueryHandler(tm *db.DBManager) talksession.BrowseJoinedTalkSessionsQuery {
	return &BrowseJoinedTalkSessionQueryHandler{
		DBManager: tm,
	}
}

// Execute ユーザーが参加しているトークセッションを検索する
func (h *BrowseJoinedTalkSessionQueryHandler) Execute(ctx context.Context, input talksession.BrowseJoinedTalkSessionsQueryInput) (*talksession.BrowseJoinedTalkSessionsQueryOutput, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "BrowseJoinedTalkSessionQueryHandler.Execute")
	defer span.End()

	if err := input.Validate(); err != nil {
		return nil, err
	}

	var status, theme sql.NullString
	// ValidateでStatusが空文字の場合はStatusOpenを設定しているので、if分岐は不要
	status = sql.NullString{String: string(input.Status), Valid: input.Status != ""}
	if input.Theme != nil {
		theme = sql.NullString{String: *input.Theme, Valid: true}
	}

	talkSessionRows, err := h.GetQueries(ctx).GetRespondTalkSessionByUserID(ctx, model.GetRespondTalkSessionByUserIDParams{
		Limit:  int32(*input.Limit),
		Offset: int32(*input.Offset),
		UserID: uuid.NullUUID{UUID: input.UserID.UUID(), Valid: true},
		Status: status,
		Theme:  theme,
	})
	if err != nil {
		utils.HandleError(ctx, err, "Listクエリに失敗")
		return nil, err
	}
	totalCount, err := h.GetQueries(ctx).CountTalkSessions(ctx, model.CountTalkSessionsParams{
		Theme:  theme,
		Status: status,
		UserID: uuid.NullUUID{UUID: input.UserID.UUID(), Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "Countクエリに失敗")
		return nil, err
	}

	var talkSessions []dto.TalkSessionWithDetail
	if err := copier.CopyWithOption(&talkSessions, talkSessionRows, copier.Option{
		DeepCopy: true,
	}); err != nil {
		utils.HandleError(ctx, err, "マッピングに失敗")
		return nil, err
	}

	return &talksession.BrowseJoinedTalkSessionsQueryOutput{
		TalkSessions: talkSessions,
		TotalCount:   int(totalCount.TalkSessionCount),
	}, nil
}

package talksession_query

import (
	"context"
	"database/sql"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type BrowseOpenedByUserQueryImpl struct {
	*db.DBManager
}

func NewBrowseOpenedByUserQueryHandler(tm *db.DBManager) talksession.BrowseOpenedByUserQuery {
	return &BrowseOpenedByUserQueryImpl{
		DBManager: tm,
	}
}

// Execute ユーザーが開いているトークセッションを検索する
func (h *BrowseOpenedByUserQueryImpl) Execute(ctx context.Context, input talksession.BrowseOpenedByUserInput) (*talksession.BrowseOpenedByUserOutput, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "BrowseOpenedByUserQueryImpl.Execute")
	defer span.End()

	if err := input.Validate(); err != nil {
		return nil, err
	}

	var out talksession.BrowseOpenedByUserOutput

	talkSessionRow, err := h.GetQueries(ctx).GetOwnTalkSessionByDisplayIDWithCount(ctx, model.GetOwnTalkSessionByDisplayIDWithCountParams{
		DisplayID: input.DisplayID,
		Limit:     utils.ToSQLNull[sql.NullInt32](input.Limit),
		Offset:    utils.ToSQLNull[sql.NullInt32](input.Offset),
		Theme:     utils.ToSQLNull[sql.NullString](input.Theme),
		Status:    utils.ToSQLNull[sql.NullString](input.Status),
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetOwnTalkSessionByIDでエラー")
		return nil, messages.TalkSessionNotFound
	}
	if len(talkSessionRow) <= 0 {
		return &out, nil
	}
	var talkSessions []dto.TalkSessionWithDetail
	if err := copier.CopyWithOption(&talkSessions, talkSessionRow, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	}); err != nil {
		utils.HandleError(ctx, err, "copier.CopyWithOptionでエラー")
		return nil, err
	}
	out.TalkSessions = talkSessions
	out.TotalCount = int32(talkSessionRow[0].TotalCount)

	return &out, nil
}

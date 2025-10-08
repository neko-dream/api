package talksession_query

import (
	"context"
	"database/sql"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/application/query/talksession"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
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
	status := sql.NullString{Valid: false}
	if input.Status != "" {
		status = sql.NullString{
			String: string(input.Status),
			Valid:  true,
		}
	}

	talkSessionRow, err := h.GetQueries(ctx).GetOwnTalkSessionByDisplayIDWithCount(ctx, model.GetOwnTalkSessionByDisplayIDWithCountParams{
		DisplayID: input.DisplayID,
		Limit:     utils.ToNullableSQL[sql.NullInt32](input.Limit),
		Offset:    utils.ToNullableSQL[sql.NullInt32](input.Offset),
		Theme:     utils.ToNullableSQL[sql.NullString](input.Theme),
		Status:    status,
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

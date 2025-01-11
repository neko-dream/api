package queryimpl

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	"github.com/neko-dream/server/internal/usecase/query/talksession"
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
	talkSessionRow, err := h.GetQueries(ctx).GetTalkSessionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	var result dto.TalkSessionWithDetail
	if err := copier.CopyWithOption(&result, talkSessionRow, copier.Option{
		DeepCopy: true,
	}); err != nil {
		return nil, err
	}

	return &talksession.GetTalkSessionDetailOutput{
		TalkSessionWithDetail: result,
	}, nil
}

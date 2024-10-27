package talk_session_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetTalkSessionConclusionQuery interface {
		Execute(context.Context, GetTalkSessionConclusionInput) (*GetTalkSessionConclusionOutput, error)
	}

	GetTalkSessionConclusionInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetTalkSessionConclusionOutput struct {
		User       UserDTO
		Conclusion string
	}

	getTalkSessionConclusionInteractor struct {
		*db.DBManager
	}
)

func NewGetTalkSessionConclusionQuery(
	DBManager *db.DBManager,
) GetTalkSessionConclusionQuery {
	return &getTalkSessionConclusionInteractor{
		DBManager: DBManager,
	}
}

// Execute implements GetTalkSessionConclusionQuery.
func (g *getTalkSessionConclusionInteractor) Execute(ctx context.Context, input GetTalkSessionConclusionInput) (*GetTalkSessionConclusionOutput, error) {
	ts, err := g.DBManager.GetQueries(ctx).GetTalkSessionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	// まだ終了していないトークセッションに対しては結論を取得できない
	if ts.ScheduledEndTime.After(time.Now(ctx).Time) {
		return nil, messages.TalkSessionNotFinished
	}

	res, err := g.DBManager.GetQueries(ctx).GetTalkSessionConclusionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &GetTalkSessionConclusionOutput{
		Conclusion: res.Content,
		User: UserDTO{
			DisplayID:   res.DisplayID.String,
			DisplayName: res.DisplayName.String,
			IconURL: utils.ToPtrIfNotNullValue(
				!res.IconUrl.Valid,
				res.IconUrl.String,
			),
		},
	}, nil
}

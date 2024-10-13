package opinion_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	GetSwipeOpinionsQueryHandler interface {
		Execute(context.Context, GetSwipeOpinionsQuery) (*GetSwipeOpinionsOutput, error)
	}

	GetSwipeOpinionsQuery struct {
		UserID        shared.UUID[user.User]
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetSwipeOpinionsOutput struct {
		Opinions []SwipeOpinionDTO
	}

	SwipeOpinionDTO struct {
		Opinion      OpinionDTO
		User         UserDTO
		CommentCount int
	}

	getSwipeOpinionsQueryHandler struct {
		*db.DBManager
	}
)

func NewGetSwipeOpinionsQueryHandler(
	dbManager *db.DBManager,
) GetSwipeOpinionsQueryHandler {
	return &getSwipeOpinionsQueryHandler{
		DBManager: dbManager,
	}
}

func (h *getSwipeOpinionsQueryHandler) Execute(ctx context.Context, q GetSwipeOpinionsQuery) (*GetSwipeOpinionsOutput, error) {

	return nil, nil
}

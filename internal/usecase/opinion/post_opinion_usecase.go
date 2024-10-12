package opinion_usecase

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	PostOpinionUseCase interface {
		Execute(context.Context, PostOpinionInput) (PostOpinionOutput, error)
	}

	PostOpinionInput struct {
		TalkSessionID   shared.UUID[talksession.TalkSession]
		OwnerID         shared.UUID[user.User]
		ParentOpinionID *shared.UUID[opinion.Opinion]
		VoteStatus      *vote.VoteStatus
		Content         string
		ReferenceURL    *string
		Picture         *multipart.FileHeader
	}

	PostOpinionOutput struct {
	}

	postOpinionInteractor struct {
		*db.DBManager
	}
)

func NewPostOpinionInteractor(
	dbManager *db.DBManager,
) PostOpinionUseCase {
	return &postOpinionInteractor{
		DBManager: dbManager,
	}
}

func (i *postOpinionInteractor) Execute(ctx context.Context, input PostOpinionInput) (PostOpinionOutput, error) {

	return PostOpinionOutput{}, nil
}

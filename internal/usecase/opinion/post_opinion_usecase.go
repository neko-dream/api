package opinion_usecase

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	PostOpinionUseCase interface {
		Execute(context.Context, PostOpinionInput) (*PostOpinionOutput, error)
	}

	PostOpinionInput struct {
		TalkSessionID   shared.UUID[talksession.TalkSession]
		OwnerID         shared.UUID[user.User]
		ParentOpinionID *shared.UUID[opinion.Opinion]
		Title           *string
		Content         string
		ReferenceURL    *string
		Picture         *multipart.FileHeader
	}

	PostOpinionOutput struct {
	}

	postOpinionInteractor struct {
		opinion.OpinionRepository
		opinion.OpinionService
		vote.VoteRepository
		*db.DBManager
	}
)

func NewPostOpinionUseCase(
	opinionRepository opinion.OpinionRepository,
	opinionService opinion.OpinionService,
	voteRepository vote.VoteRepository,
	dbManager *db.DBManager,
) PostOpinionUseCase {
	return &postOpinionInteractor{
		DBManager:         dbManager,
		OpinionService:    opinionService,
		OpinionRepository: opinionRepository,
		VoteRepository:    voteRepository,
	}
}

func (i *postOpinionInteractor) Execute(ctx context.Context, input PostOpinionInput) (*PostOpinionOutput, error) {
	if err := i.ExecTx(ctx, func(ctx context.Context) error {
		opinion, err := opinion.NewOpinion(
			shared.NewUUID[opinion.Opinion](),
			input.TalkSessionID,
			input.OwnerID,
			input.ParentOpinionID,
			input.Title,
			input.Content,
			time.Now(),
			input.ReferenceURL,
		)
		if err != nil {
			utils.HandleError(ctx, err, "NewOpinion")
			return err
		}
		if input.Picture != nil {
			if err := opinion.SetReferenceImage(ctx, input.Picture); err != nil {
				utils.HandleError(ctx, err, "SetReferenceImage")
				return err
			}
		}

		if err := i.OpinionRepository.Create(ctx, *opinion); err != nil {
			utils.HandleError(ctx, err, "OpinionRepository.Create")
			return messages.OpinionCreateFailed
		}

		// 自分の意見には必ず投票を紐付ける
		v, err := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			opinion.OpinionID(),
			input.TalkSessionID,
			input.OwnerID,
			vote.Agree,
			time.Now(),
		)
		if err != nil {
			utils.HandleError(ctx, err, "NewVote")
			return err
		}
		if err := i.VoteRepository.Create(ctx, *v); err != nil {
			return messages.VoteFailed
		}

		return nil
	}); err != nil {
		return nil, err

	}

	out := &PostOpinionOutput{}
	return out, nil
}

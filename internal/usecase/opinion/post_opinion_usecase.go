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
		VoteType        *string
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
		)
		if err != nil {
			return err
		}
		if err := i.OpinionRepository.Create(ctx, *opinion); err != nil {
			return messages.OpinionCreateFailed
		}

		// 自分の意見には必ず投票を紐付ける
		v, err := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			opinion.OpinionID(),
			input.OwnerID,
			vote.Agreed,
			time.Now(),
		)
		if err != nil {
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

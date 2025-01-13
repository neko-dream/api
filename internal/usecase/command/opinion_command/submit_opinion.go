package opinion_command

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	SubmitOpinion interface {
		Execute(context.Context, SubmitOpinionInput) error
	}

	SubmitOpinionInput struct {
		TalkSessionID   shared.UUID[talksession.TalkSession]
		OwnerID         shared.UUID[user.User]
		ParentOpinionID *shared.UUID[opinion.Opinion]
		Title           *string
		Content         string
		ReferenceURL    *string
		Picture         *multipart.FileHeader
	}

	submitOpinionHandler struct {
		opinion.OpinionRepository
		opinion.OpinionService
		vote.VoteRepository
		*db.DBManager
	}
)

func NewSubmitOpinionHandler(
	opinionRepository opinion.OpinionRepository,
	opinionService opinion.OpinionService,
	voteRepository vote.VoteRepository,
	dbManager *db.DBManager,
) SubmitOpinion {
	return &submitOpinionHandler{
		DBManager:         dbManager,
		OpinionService:    opinionService,
		OpinionRepository: opinionRepository,
		VoteRepository:    voteRepository,
	}
}

func (h *submitOpinionHandler) Execute(ctx context.Context, input SubmitOpinionInput) error {
	ctx, span := otel.Tracer("opinion_command").Start(ctx, "submitOpinionHandler.Execute")
	defer span.End()

	if err := h.ExecTx(ctx, func(ctx context.Context) error {
		opinion, err := opinion.NewOpinion(
			shared.NewUUID[opinion.Opinion](),
			input.TalkSessionID,
			input.OwnerID,
			input.ParentOpinionID,
			input.Title,
			input.Content,
			clock.Now(ctx),
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

		if err := h.OpinionRepository.Create(ctx, *opinion); err != nil {
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
			clock.Now(ctx),
		)
		if err != nil {
			utils.HandleError(ctx, err, "NewVote")
			return err
		}
		if err := h.VoteRepository.Create(ctx, *v); err != nil {
			return messages.VoteFailed
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

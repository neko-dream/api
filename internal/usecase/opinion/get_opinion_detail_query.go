package opinion_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetOpinionDetailUseCase interface {
		Execute(context.Context, GetOpinionDetailInput) (*GetOpinionDetailOutput, error)
	}

	GetOpinionDetailInput struct {
		OpinionID shared.UUID[opinion.Opinion]
		UserID    *shared.UUID[user.User]
	}

	GetOpinionDetailOutput struct {
		Opinion ReplyDTO
	}

	getOpinionDetailInteractor struct {
		*db.DBManager
	}
)

func NewGetOpinionDetailUseCase(
	dbManager *db.DBManager,
) GetOpinionDetailUseCase {
	return &getOpinionDetailInteractor{
		DBManager: dbManager,
	}
}

// Execute implements GetOpinionDetailUseCase.
func (g *getOpinionDetailInteractor) Execute(ctx context.Context, input GetOpinionDetailInput) (*GetOpinionDetailOutput, error) {
	var userID uuid.NullUUID
	if input.UserID != nil {
		userID = uuid.NullUUID{UUID: input.UserID.UUID(), Valid: true}
	}

	opinionRow, err := g.DBManager.GetQueries(ctx).GetOpinionByID(ctx, model.GetOpinionByIDParams{
		OpinionID: input.OpinionID.UUID(),
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}
	opinion := OpinionDTO{
		OpinionID:       opinionRow.OpinionID.String(),
		TalkSessionID:   opinionRow.TalkSessionID.String(),
		UserID:          opinionRow.UserID.String(),
		ParentOpinionID: utils.ToPtrIfNotNullValue(!opinionRow.ParentOpinionID.Valid, opinionRow.ParentOpinionID.UUID.String()),
		Title:           utils.ToPtrIfNotNullValue(!opinionRow.Title.Valid, opinionRow.Title.String),
		Content:         opinionRow.Content,
		CreatedAt:       opinionRow.CreatedAt,
		VoteType:        vote.VoteTypeFromInt(int(opinionRow.VoteType)).String(),
		ReferenceURL:    utils.ToPtrIfNotNullValue(!opinionRow.ReferenceUrl.Valid, opinionRow.ReferenceUrl.String),
		PictureURL:      utils.ToPtrIfNotNullValue(!opinionRow.PictureUrl.Valid, opinionRow.PictureUrl.String),
	}
	user := UserDTO{
		ID:   opinionRow.UserID.String(),
		Name: opinionRow.DisplayName.String,
		Icon: utils.ToPtrIfNotNullValue(!opinionRow.IconUrl.Valid, opinionRow.IconUrl.String),
	}
	opinionDTO := ReplyDTO{
		Opinion:    opinion,
		User:       user,
		MyVoteType: vote.VoteTypeFromInt(int(opinionRow.CurrentVoteType)).String(),
	}

	return &GetOpinionDetailOutput{
		Opinion: opinionDTO,
	}, nil
}

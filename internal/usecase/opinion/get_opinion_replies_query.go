package opinion_usecase

import (
	"context"
	"time"

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
	GetOpinionRepliesUseCase interface {
		Execute(context.Context, GetOpinionRepliesInput) (*GetOpinionRepliesOutput, error)
	}

	GetOpinionRepliesInput struct {
		OpinionID shared.UUID[opinion.Opinion]
		UserID    *shared.UUID[user.User]
	}

	GetOpinionRepliesOutput struct {
		RootOpinion ReplyDTO
		Replies     []ReplyDTO
	}

	// RepliesDTO 親意見に対するリプライ意見
	ReplyDTO struct {
		Opinion    OpinionDTO
		User       UserDTO
		MyVoteType string
	}
	UserDTO struct {
		ID   string
		Name string
		Icon *string
	}
	OpinionDTO struct {
		OpinionID       string
		TalkSessionID   string
		UserID          string
		ParentOpinionID *string
		Title           *string
		Content         string
		CreatedAt       time.Time
		VoteType        string
		PictureURL      *string
		ReferenceURL    *string
	}

	GetOpinionRepliesInteractor struct {
		opinion.OpinionRepository
		opinion.OpinionService
		vote.VoteRepository
		*db.DBManager
	}
)

func NewGetOpinionRepliesUseCase(
	opinionRepository opinion.OpinionRepository,
	opinionService opinion.OpinionService,
	voteRepository vote.VoteRepository,
	dbManager *db.DBManager,
) GetOpinionRepliesUseCase {
	return &GetOpinionRepliesInteractor{
		DBManager:         dbManager,
		OpinionService:    opinionService,
		OpinionRepository: opinionRepository,
		VoteRepository:    voteRepository,
	}
}

func (i *GetOpinionRepliesInteractor) Execute(ctx context.Context, input GetOpinionRepliesInput) (*GetOpinionRepliesOutput, error) {
	var userID uuid.NullUUID
	if input.UserID != nil {
		userID = uuid.NullUUID{Valid: true, UUID: input.UserID.UUID()}
	}

	// 親意見を取得
	opinionRow, err := i.GetQueries(ctx).GetOpinionByID(ctx, model.GetOpinionByIDParams{
		OpinionID: input.OpinionID.UUID(),
		UserID:    userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetOpinionRepliesInteractor.Execute")
		return nil, err
	}
	rootOpinion := OpinionDTO{
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
	rootUser := UserDTO{
		ID:   opinionRow.UserID.String(),
		Name: opinionRow.DisplayName.String,
		Icon: utils.ToPtrIfNotNullValue(!opinionRow.IconUrl.Valid, opinionRow.IconUrl.String),
	}
	rootOpinionDTO := ReplyDTO{
		Opinion:    rootOpinion,
		User:       rootUser,
		MyVoteType: vote.VoteTypeFromInt(int(opinionRow.CurrentVoteType)).String(),
	}

	row, err := i.GetQueries(ctx).GetOpinionReplies(ctx, model.GetOpinionRepliesParams{
		OpinionID: input.OpinionID.UUID(),
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}
	replies := make([]ReplyDTO, 0, len(row))
	for _, r := range row {
		replies = append(replies, ReplyDTO{
			Opinion: OpinionDTO{
				OpinionID:       r.OpinionID.String(),
				TalkSessionID:   r.TalkSessionID.String(),
				UserID:          r.UserID.String(),
				ParentOpinionID: utils.ToPtrIfNotNullValue(!r.ParentOpinionID.Valid, r.ParentOpinionID.UUID.String()),
				Title:           utils.ToPtrIfNotNullValue(!r.Title.Valid, r.Title.String),
				Content:         r.Content,
				CreatedAt:       r.CreatedAt,
				VoteType:        vote.VoteTypeFromInt(int(r.VoteType)).String(),
				ReferenceURL:    utils.ToPtrIfNotNullValue(!r.ReferenceUrl.Valid, r.ReferenceUrl.String),
				PictureURL:      utils.ToPtrIfNotNullValue(!r.PictureUrl.Valid, r.PictureUrl.String),
			},
			User: UserDTO{
				ID:   r.UserID.String(),
				Name: r.DisplayName.String,
				Icon: utils.ToPtrIfNotNullValue(!r.IconUrl.Valid, r.IconUrl.String),
			},
			MyVoteType: vote.VoteTypeFromInt(int(r.CurrentVoteType)).String(),
		})
	}

	return &GetOpinionRepliesOutput{
		RootOpinion: rootOpinionDTO,
		Replies:     replies,
	}, nil
}

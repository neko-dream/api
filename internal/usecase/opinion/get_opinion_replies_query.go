package opinion_usecase

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetOpinionRepliesUseCase interface {
		Execute(context.Context, GetOpinionRepliesInput) (*GetOpinionRepliesOutput, error)
	}

	GetOpinionRepliesInput struct {
		OpinionID shared.UUID[opinion.Opinion]
	}

	GetOpinionRepliesOutput struct {
		RootOpinion ReplyDTO
		Replies     []ReplyDTO
	}

	// RepliesDTO 親意見に対するリプライ意見
	ReplyDTO struct {
		Opinion OpinionDTO
		User    UserDTO
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

	// // 親意見を取得
	opinionRow, err := i.GetQueries(ctx).GetOpinionByID(ctx, input.OpinionID.UUID())
	if err != nil {
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
	}
	rootUser := UserDTO{
		ID:   opinionRow.UserID.String(),
		Name: opinionRow.DisplayName.String,
		Icon: utils.ToPtrIfNotNullValue(!opinionRow.IconUrl.Valid, opinionRow.IconUrl.String),
	}
	rootOpinionDTO := ReplyDTO{
		Opinion: rootOpinion,
		User:    rootUser,
	}

	row, err := i.GetQueries(ctx).GetOpinionReplies(ctx, input.OpinionID.UUID())
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
			},
			User: UserDTO{
				ID:   r.UserID.String(),
				Name: r.DisplayName.String,
				Icon: utils.ToPtrIfNotNullValue(!r.IconUrl.Valid, r.IconUrl.String),
			},
		})
	}

	return &GetOpinionRepliesOutput{
		RootOpinion: rootOpinionDTO,
		Replies:     replies,
	}, nil
}

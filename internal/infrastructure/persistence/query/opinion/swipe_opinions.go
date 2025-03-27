package opinion_query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type GetSwipeOpinionsQueryHandler struct {
	*db.DBManager
}

func NewSwipeOpinionsQueryHandler(
	dbManager *db.DBManager,
) opinion_query.GetSwipeOpinionsQuery {
	return &GetSwipeOpinionsQueryHandler{
		DBManager: dbManager,
	}
}

func (g *GetSwipeOpinionsQueryHandler) Execute(ctx context.Context, in opinion_query.GetSwipeOpinionsQueryInput) (*opinion_query.GetSwipeOpinionsQueryOutput, error) {
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "GetSwipeOpinionsQueryHandler.Execute")
	defer span.End()

	swipeableOpinionCount, err := g.GetQueries(ctx).CountSwipeableOpinions(ctx, model.CountSwipeableOpinionsParams{
		UserID:        in.UserID.UUID(),
		TalkSessionID: in.TalkSessionID.UUID(),
	})
	if err != nil {
		utils.HandleError(ctx, err, "SwipeableOpinionのカウントに失敗")
		return nil, err
	}
	if swipeableOpinionCount == 0 {
		return &opinion_query.GetSwipeOpinionsQueryOutput{
			Opinions:          []dto.SwipeOpinion{},
			RemainingOpinions: 0,
		}, nil
	}

	// top,randomを1:2の比率で取得する
	// limitが3以上の場合、2件はtop, 1件はrandomで取得する
	topLimit := in.Limit / 3

	var swipeRows []model.GetRandomOpinionsRow
	if in.Limit >= 3 {
		// top & random
		// top
		swipeRow, err := g.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
			UserID:        in.UserID.UUID(),
			TalkSessionID: in.TalkSessionID.UUID(),
			Limit:         int32(topLimit),
			SortKey:       sql.NullString{String: "top", Valid: true},
		})
		if err != nil {
			utils.HandleError(ctx, err, "TopN意見の取得に失敗")
			return nil, err
		}
		swipeRows = swipeRow
		// random
		// randomはlimitより取得できたtopの数を引いた数だけ取得する
		randomLimit := in.Limit - len(swipeRows)
		randomSwipeRow, err := g.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
			UserID:        in.UserID.UUID(),
			TalkSessionID: in.TalkSessionID.UUID(),
			Limit:         int32(randomLimit),
		})
		if err != nil {
			utils.HandleError(ctx, err, "ランダムな意見の取得に失敗")
			return nil, err
		}
		// topとrandomを結合する
		swipeRows = append(swipeRows, randomSwipeRow...)

		// 同様の意見を除外する
		swipeRows = lo.UniqBy(swipeRows, func(swipe model.GetRandomOpinionsRow) uuid.UUID {
			return swipe.Opinion.OpinionID
		})
	} else {
		// randomのみ
		swipeRow, err := g.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
			UserID:        in.UserID.UUID(),
			TalkSessionID: in.TalkSessionID.UUID(),
			Limit:         int32(in.Limit),
		})
		if err != nil {
			utils.HandleError(ctx, err, "ランダムな意見の取得に失敗")
			return nil, err
		}
		swipeRows = swipeRow
	}

	var swipeOpinions []dto.SwipeOpinion
	for _, swipeRow := range swipeRows {
		var swipeOpinion dto.SwipeOpinion
		if err := copier.CopyWithOption(&swipeOpinion, &swipeRow, copier.Option{
			DeepCopy:    true,
			IgnoreEmpty: true,
		}); err != nil {
			utils.HandleError(ctx, err, "マッピングに失敗")
			return nil, err
		}
		swipeOpinions = append(swipeOpinions, swipeOpinion)
	}

	return &opinion_query.GetSwipeOpinionsQueryOutput{
		Opinions:          swipeOpinions,
		RemainingOpinions: int(swipeableOpinionCount),
	}, nil
}

package opinion_query

import (
	"context"
	"database/sql"
	"errors"

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

	// swipeAbleOpinionCountよりlimitの方が大きい場合、limitを上書きする
	if int64(in.Limit) > swipeableOpinionCount {
		in.Limit = int(swipeableOpinionCount)
	}

	var swipeOpinions []dto.SwipeOpinion
	// 一旦seed意見を取得する
	rows, err := g.GetQueries(ctx).GetSeedOpinions(ctx, model.GetSeedOpinionsParams{
		UserID:        in.UserID.UUID(),
		TalkSessionID: in.TalkSessionID.UUID(),
		Limit:         int32(in.Limit),
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(ctx, err, "Seed意見の取得に失敗")
			return nil, err
		}
	}

	// seedでlimitを満たしたらそれを返す
	for _, swipeRow := range rows {
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
	if len(rows) >= in.Limit {
		return &opinion_query.GetSwipeOpinionsQueryOutput{
			Opinions:          swipeOpinions,
			RemainingOpinions: int(swipeableOpinionCount),
		}, nil
	}
	seedOpinionIDs := lo.Map(rows, func(swipe model.GetSeedOpinionsRow, _ int) uuid.UUID {
		return swipe.Opinion.OpinionID
	})

	// 超えなかったら、seedを除いた残りの意見を取得する
	remainingLimit := in.Limit - len(rows)

	// top,randomを1:2の比率で取得する
	// limitが3以上の場合、2件はtop, 1件はrandomで取得する
	// シード意見取得後の残りで、1/3をトップ意見に割り当てる
	topLimit := remainingLimit / 3

	// top
	topRows, err := g.GetQueries(ctx).GetOpinionsByRank(ctx, model.GetOpinionsByRankParams{
		UserID:        in.UserID.UUID(),
		TalkSessionID: in.TalkSessionID.UUID(),
		Rank:          1,
		Limit:         int32(topLimit),
	})
	if err != nil {
		utils.HandleError(ctx, err, "TopN意見の取得に失敗")
		return nil, err
	}
	for _, swipeRow := range topRows {
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
	topOpinionIDs := lo.Map(topRows, func(swipe model.GetOpinionsByRankRow, _ int) uuid.UUID {
		return swipe.Opinion.OpinionID
	})

	// random
	// randomはlimitより取得できたtopの数を引いた数だけ取得する
	randomLimit := in.Limit - len(swipeOpinions)
	if randomLimit > 0 && (swipeableOpinionCount-int64(len(topRows))) > 0 {
		excludes := append(topOpinionIDs, seedOpinionIDs...)
		randomSwipeRow, err := g.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
			UserID:            in.UserID.UUID(),
			TalkSessionID:     in.TalkSessionID.UUID(),
			Limit:             int32(randomLimit),
			ExcludeOpinionIds: excludes,
			ExcludesLen:       int32(len(excludes)),
		})
		if err != nil {
			utils.HandleError(ctx, err, "ランダムな意見の取得に失敗")
			return nil, err
		}
		for _, swipeRow := range randomSwipeRow {
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
	}

	return &opinion_query.GetSwipeOpinionsQueryOutput{
		Opinions:          swipeOpinions,
		RemainingOpinions: int(swipeableOpinionCount),
	}, nil
}

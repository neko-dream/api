package opinion_query

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/api/internal/application/query/dto"
	opinion_query "github.com/neko-dream/api/internal/application/query/opinion"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
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

func convertToSwipeOpinion(source any) (dto.SwipeOpinion, error) {
	var swipeOpinion dto.SwipeOpinion
	if err := copier.CopyWithOption(&swipeOpinion, source, copier.Option{
		DeepCopy:    true,
		IgnoreEmpty: true,
	}); err != nil {
		return dto.SwipeOpinion{}, err
	}
	return swipeOpinion, nil
}

func (g *GetSwipeOpinionsQueryHandler) Execute(ctx context.Context, in opinion_query.GetSwipeOpinionsQueryInput) (*opinion_query.GetSwipeOpinionsQueryOutput, error) {
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "GetSwipeOpinionsQueryHandler.Execute")
	defer span.End()

	talkSession, err := g.GetQueries(ctx).GetTalkSessionByID(ctx, in.TalkSessionID.UUID())
	if err != nil {
		utils.HandleError(ctx, err, "トークセッションの取得に失敗")
		return nil, err
	}

	if talkSession.TalkSession.ScheduledEndTime.Before(clock.Now(ctx)) {
		return &opinion_query.GetSwipeOpinionsQueryOutput{
			Opinions:          []dto.SwipeOpinion{},
			RemainingOpinions: 0,
		}, nil
	}

	// スワイプ可能な意見の総数を取得
	swipeableOpinionCount, err := g.GetQueries(ctx).CountSwipeableOpinions(ctx, model.CountSwipeableOpinionsParams{
		UserID:        in.UserID.UUID(),
		TalkSessionID: in.TalkSessionID.UUID(),
	})
	if err != nil {
		utils.HandleError(ctx, err, "スワイプ可能な意見のカウントに失敗")
		return nil, err
	}

	// スワイプ可能な意見がない場合は空の結果を返す
	if swipeableOpinionCount == 0 {
		return &opinion_query.GetSwipeOpinionsQueryOutput{
			Opinions:          []dto.SwipeOpinion{},
			RemainingOpinions: 0,
		}, nil
	}

	// 取得限度数の調整：要求limitが利用可能な意見数より多い場合は調整
	requestLimit := in.Limit
	if int64(requestLimit) > swipeableOpinionCount {
		requestLimit = int(swipeableOpinionCount)
	}

	var allSwipeOpinions []dto.SwipeOpinion
	var collectedOpinionIDs []uuid.UUID

	// シード意見の取得（初期データとなる意見）
	seedOpinions, seedOpinionIDs, err := g.fetchSeedOpinions(ctx, in.UserID.UUID(), in.TalkSessionID.UUID(), requestLimit)
	if err != nil {
		return nil, err // エラーは内部関数で既にログ記録されている
	}
	allSwipeOpinions = append(allSwipeOpinions, seedOpinions...)
	collectedOpinionIDs = append(collectedOpinionIDs, seedOpinionIDs...)

	// シード意見だけで要求数を満たした場合はそのまま返す
	if len(allSwipeOpinions) >= requestLimit {
		return &opinion_query.GetSwipeOpinionsQueryOutput{
			Opinions:          allSwipeOpinions[:requestLimit],
			RemainingOpinions: int(swipeableOpinionCount),
		}, nil
	}

	// 残りの枠を埋めるためにトップ意見とランダム意見を取得
	remainingLimit := requestLimit - len(allSwipeOpinions)

	// トップ意見は残りの1/3を割り当て
	topLimit := remainingLimit / 3
	if topLimit > 0 {
		topOpinions, topOpinionIDs, err := g.fetchTopOpinions(
			ctx,
			in.UserID.UUID(),
			in.TalkSessionID.UUID(),
			topLimit,
			collectedOpinionIDs,
		)
		if err != nil {
			return nil, err
		}
		allSwipeOpinions = append(allSwipeOpinions, topOpinions...)
		collectedOpinionIDs = append(collectedOpinionIDs, topOpinionIDs...)
	}

	// ランダム意見を取得して残りを埋める
	randomLimit := requestLimit - len(allSwipeOpinions)
	if randomLimit > 0 && (swipeableOpinionCount-int64(len(allSwipeOpinions))) > 0 {
		randomOpinions, err := g.fetchRandomOpinions(
			ctx,
			in.UserID.UUID(),
			in.TalkSessionID.UUID(),
			randomLimit,
			collectedOpinionIDs,
		)
		if err != nil {
			return nil, err
		}
		allSwipeOpinions = append(allSwipeOpinions, randomOpinions...)
	}

	return &opinion_query.GetSwipeOpinionsQueryOutput{
		Opinions:          allSwipeOpinions,
		RemainingOpinions: int(swipeableOpinionCount),
	}, nil
}

// シード意見を取得する
func (g *GetSwipeOpinionsQueryHandler) fetchSeedOpinions(
	ctx context.Context,
	userID uuid.UUID,
	talkSessionID uuid.UUID,
	limit int,
) ([]dto.SwipeOpinion, []uuid.UUID, error) {
	seedRows, err := g.GetQueries(ctx).GetSeedOpinions(ctx, model.GetSeedOpinionsParams{
		UserID:        userID,
		TalkSessionID: talkSessionID,
		Limit:         int32(limit),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 結果がない場合は空のスライスを返す
			return []dto.SwipeOpinion{}, []uuid.UUID{}, nil
		}
		utils.HandleError(ctx, err, "シード意見の取得に失敗")
		return nil, nil, err
	}

	seedOpinions := make([]dto.SwipeOpinion, 0, len(seedRows))
	seedOpinionIDs := make([]uuid.UUID, 0, len(seedRows))

	for _, row := range seedRows {
		swipeOpinion, err := convertToSwipeOpinion(row)
		if err != nil {
			utils.HandleError(ctx, err, "シード意見のマッピングに失敗")
			return nil, nil, err
		}
		seedOpinions = append(seedOpinions, swipeOpinion)
		seedOpinionIDs = append(seedOpinionIDs, row.Opinion.OpinionID)
	}

	return seedOpinions, seedOpinionIDs, nil
}

// トップランクの意見を取得する
func (g *GetSwipeOpinionsQueryHandler) fetchTopOpinions(
	ctx context.Context,
	userID uuid.UUID,
	talkSessionID uuid.UUID,
	limit int,
	excludeOpinionIDs []uuid.UUID,
) ([]dto.SwipeOpinion, []uuid.UUID, error) {
	topRows, err := g.GetQueries(ctx).GetOpinionsByRank(ctx, model.GetOpinionsByRankParams{
		UserID:            userID,
		TalkSessionID:     talkSessionID,
		Rank:              1, // ランク1のトップ意見
		Limit:             int32(limit),
		ExcludeOpinionIds: excludeOpinionIDs,
	})
	if err != nil {
		utils.HandleError(ctx, err, "トップ意見の取得に失敗")
		return nil, nil, err
	}

	topOpinions := make([]dto.SwipeOpinion, 0, len(topRows))
	topOpinionIDs := make([]uuid.UUID, 0, len(topRows))

	for _, row := range topRows {
		swipeOpinion, err := convertToSwipeOpinion(row)
		if err != nil {
			utils.HandleError(ctx, err, "トップ意見のマッピングに失敗")
			return nil, nil, err
		}
		topOpinions = append(topOpinions, swipeOpinion)
		topOpinionIDs = append(topOpinionIDs, row.Opinion.OpinionID)
	}

	return topOpinions, topOpinionIDs, nil
}

// ランダム意見を取得する
func (g *GetSwipeOpinionsQueryHandler) fetchRandomOpinions(
	ctx context.Context,
	userID uuid.UUID,
	talkSessionID uuid.UUID,
	limit int,
	excludeOpinionIDs []uuid.UUID,
) ([]dto.SwipeOpinion, error) {
	randomRows, err := g.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
		UserID:            userID,
		TalkSessionID:     talkSessionID,
		Limit:             int32(limit),
		ExcludeOpinionIds: excludeOpinionIDs,
	})
	if err != nil {
		utils.HandleError(ctx, err, "ランダム意見の取得に失敗")
		return nil, err
	}

	randomOpinions := make([]dto.SwipeOpinion, 0, len(randomRows))

	for _, row := range randomRows {
		swipeOpinion, err := convertToSwipeOpinion(row)
		if err != nil {
			utils.HandleError(ctx, err, "ランダム意見のマッピングに失敗")
			return nil, err
		}
		randomOpinions = append(randomOpinions, swipeOpinion)
	}

	return randomOpinions, nil
}

package talksession_usecase

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	ToggleAnalysisUseCase interface {
		Execute(ctx context.Context, input ToggleAnalysisInput) (*ToggleAnalysisOutput, error)
	}

	ToggleAnalysisInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		UserID        shared.UUID[user.User]
	}

	ToggleAnalysisOutput struct {
		DisableAnalysis bool
	}

	toggleAnalysisHandler struct {
		talksession.TalkSessionRepository
		analysis_usecase.AnalysisScheduler
		*db.DBManager
		*config.Config
	}
)

func NewToggleAnalysisUseCase(
	talkSessionRepository talksession.TalkSessionRepository,
	analysisScheduler analysis_usecase.AnalysisScheduler,
	DBManager *db.DBManager,
	config *config.Config,
) ToggleAnalysisUseCase {
	return &toggleAnalysisHandler{
		TalkSessionRepository: talkSessionRepository,
		AnalysisScheduler:     analysisScheduler,
		DBManager:             DBManager,
		Config:                config,
	}
}

func (h *toggleAnalysisHandler) Execute(ctx context.Context, input ToggleAnalysisInput) (*ToggleAnalysisOutput, error) {
	ctx, span := otel.Tracer("talksession_usecase").Start(ctx, "toggleAnalysisHandler.Execute")
	defer span.End()

	// セッションを取得
	talkSession, err := h.TalkSessionRepository.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "TalkSessionRepository.FindByID")
		return nil, messages.TalkSessionNotFound
	}

	// ローカル環境以外ではオーナーのみ編集可能
	if h.Config.Env != config.LOCAL && talkSession.OwnerUserID() != input.UserID {
		return nil, messages.ForbiddenError
	}

	// 変更前の状態を保存
	wasDisabled := talkSession.DisableAnalysis()

	// トグル
	newState := !wasDisabled
	talkSession.SetDisableAnalysis(newState)

	// 保存
	if err := h.DBManager.ExecTx(ctx, func(ctx context.Context) error {
		if err := h.TalkSessionRepository.Update(ctx, talkSession); err != nil {
			utils.HandleError(ctx, err, "TalkSessionRepository.Update")
			return messages.TalkSessionUpdateFailed
		}
		return nil
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	// 分析が有効化された場合（disable: true → false）、分析を開始
	if wasDisabled && !newState {
		h.AnalysisScheduler.ScheduleAnalysis(ctx, input.TalkSessionID, analysis_usecase.ScheduleImmediate)
	}

	return &ToggleAnalysisOutput{
		DisableAnalysis: newState,
	}, nil
}

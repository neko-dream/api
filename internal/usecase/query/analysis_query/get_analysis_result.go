package analysis_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type (
	GetAnalysisResult interface {
		Execute(context.Context, GetAnalysisResultInput) (*GetAnalysisResultOutput, error)
	}

	GetAnalysisResultInput struct {
		UserID        *shared.UUID[user.User]
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetAnalysisResultOutput struct {
		// ユーザーがログインしていれば、自分の位置情報を返す
		MyPosition *dto.UserPosition
		// トークセッションの全てのポジション情報を返す
		Positions []dto.UserPosition
		// トークセッションの全てのグループの意見情報を返す
		GroupOpinions []dto.OpinionGroup
	}
)

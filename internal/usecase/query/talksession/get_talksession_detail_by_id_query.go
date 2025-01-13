package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type (
	GetTalkSessionDetailByIDQuery interface {
		Execute(context.Context, GetTalkSessionDetailInput) (*GetTalkSessionDetailOutput, error)
	}

	GetTalkSessionDetailInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetTalkSessionDetailOutput struct {
		dto.TalkSessionWithDetail
	}
)

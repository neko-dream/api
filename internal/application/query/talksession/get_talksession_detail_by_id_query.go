package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/application/query/dto"
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

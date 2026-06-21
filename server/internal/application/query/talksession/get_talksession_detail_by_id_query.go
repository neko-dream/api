package talksession

import (
	"context"

	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
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

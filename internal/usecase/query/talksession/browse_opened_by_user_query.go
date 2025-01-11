package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type (
	BrowseOpenedByUserQuery interface {
		Execute(context.Context, BrowseOpenedByUserInput) (*BrowseOpenedByUserOutput, error)
	}

	BrowseOpenedByUserInput struct {
		UserID shared.UUID[user.User]
		Limit  int
		Offset int
		Status Status
		Theme  *string
	}

	BrowseOpenedByUserOutput struct {
		TalkSessions []dto.TalkSessionWithDetail
	}
)

func (h *BrowseOpenedByUserInput) Validate() error {
	if h.Status == "" {
		h.Status = StatusOpen
	}
	if h.Status != StatusOpen && h.Status != StatusClosed {
		return messages.TalkSessionValidationFailed
	}

	return nil
}

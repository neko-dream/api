package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	"github.com/samber/lo"
)

type (
	BrowseJoinedTalkSessionsQuery interface {
		Execute(context.Context, BrowseJoinedTalkSessionsQueryInput) (*BrowseJoinedTalkSessionsQueryOutput, error)
	}

	BrowseJoinedTalkSessionsQueryInput struct {
		UserID shared.UUID[user.User]
		Limit  *int
		Offset *int
		Theme  *string
		Status Status
	}

	BrowseJoinedTalkSessionsQueryOutput struct {
		TalkSessions []dto.TalkSessionWithDetail
		TotalCount   int
	}
)

func (h *BrowseJoinedTalkSessionsQueryInput) Validate() error {
	if h.Status == "" {
		h.Status = StatusOpen
	}
	if h.Status != StatusOpen && h.Status != StatusClosed {
		return messages.TalkSessionValidationFailed
	}
	if h.Limit == nil {
		h.Limit = lo.ToPtr(10)
	}
	if h.Offset == nil {
		h.Offset = lo.ToPtr(0)
	}

	return nil
}

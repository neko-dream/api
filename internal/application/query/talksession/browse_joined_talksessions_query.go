package talksession

import (
	"context"
	"errors"
	"fmt"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/application/query/dto"
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
	var err error

	if h.Status == "" {
		h.Status = StatusOpen
	}
	if h.Status != StatusOpen && h.Status != StatusClosed {
		err = errors.Join(err, fmt.Errorf("無効なステータスです。: %s", h.Status))
	}
	if h.Limit == nil {
		h.Limit = lo.ToPtr(10)
	} else if *h.Limit <= 0 || *h.Limit > 100 {
		err = errors.Join(err, fmt.Errorf("Limitは1から100の間で指定してください"))
	}

	if h.Offset == nil {
		h.Offset = lo.ToPtr(0)
	} else if *h.Offset < 0 {
		err = errors.Join(err, fmt.Errorf("Offsetは0以上の値を指定してください"))
	}

	return err
}

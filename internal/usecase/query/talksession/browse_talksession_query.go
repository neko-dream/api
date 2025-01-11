package talksession

import (
	"context"
	"errors"
	"fmt"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	"github.com/samber/lo"
)

type (
	BrowseTalkSessionQuery interface {
		Execute(context.Context, BrowseTalkSessionQueryInput) (*BrowseTalkSessionQueryOutput, error)
	}

	BrowseTalkSessionQueryInput struct {
		Limit     *int
		Offset    *int
		Theme     *string
		Status    Status
		SortKey   *SortKey // デフォルト: latest
		Latitude  *float64
		Longitude *float64
	}

	BrowseTalkSessionQueryOutput struct {
		TalkSessions []dto.TalkSessionWithDetail
		TotalCount   int
		Limit        int
		Offset       int
	}
)

type SortKey string

const (
	SortKeyLatest      SortKey = "latest"
	SortKeyNearest     SortKey = "nearest"
	SortKeyMostReplies SortKey = "mostReplies"
	SortKeyOldest      SortKey = "oldest"
)

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "finished"
)

func (h *BrowseTalkSessionQueryInput) Validate() error {
	if h.SortKey != nil {
		switch *h.SortKey {
		case SortKeyLatest, SortKeyNearest, SortKeyMostReplies, SortKeyOldest:
		default:
			return messages.TalkSessionValidationFailed
		}
	}
	// SortKeyがnilの場合はlatestにする
	if h.SortKey == nil {
		h.SortKey = new(SortKey)
		*h.SortKey = SortKeyLatest
	}
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

	return nil
}

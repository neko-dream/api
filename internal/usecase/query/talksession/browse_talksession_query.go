package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type (
	BrowseTalkSessionQuery interface {
		Execute(context.Context, BrowseTalkSessionQueryInput) (*BrowseTalkSessionQueryOutput, error)
	}

	BrowseTalkSessionQueryInput struct {
		Limit     int
		Offset    int
		Theme     *string
		Status    string
		SortKey   *SortKey // デフォルト: latest
		Latitude  *float64
		Longitude *float64
	}

	BrowseTalkSessionQueryOutput struct {
		TalkSessions []dto.TalkSessionWithDetail
		TotalCount   int
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

func (b *BrowseTalkSessionQueryInput) Validate() error {
	if b.SortKey != nil {
		switch *b.SortKey {
		case SortKeyLatest, SortKeyNearest, SortKeyMostReplies, SortKeyOldest:
		default:
			return messages.TalkSessionValidationFailed
		}
	}
	// SortKeyがnilの場合はlatestにする
	if b.SortKey == nil {
		b.SortKey = new(SortKey)
		*b.SortKey = SortKeyLatest
	}
	// statusが空の場合はopenにする
	if b.Status == "" {
		b.Status = string(StatusOpen)
	}
	if b.Status != string(StatusOpen) && b.Status != string(StatusClosed) {
		return messages.TalkSessionValidationFailed
	}

	return nil
}

package talksession

import (
	"context"
	"errors"
	"fmt"

	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/pkg/sort"
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
		Status    *Status
		SortKey   *sort.SortKey
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

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "finished"
)

func (h *BrowseTalkSessionQueryInput) Validate() error {
	var errs []error

	if h.SortKey != nil && h.SortKey.IsValid() {
		errs = append(errs, fmt.Errorf("無効なSortKeyです。: %s", h.SortKey))
	}

	if h.Status != nil && (*h.Status == "" || (*h.Status != StatusOpen && *h.Status != StatusClosed)) {
		errs = append(errs, fmt.Errorf("無効なステータスです。: %s", *h.Status))
	}
	if h.Limit == nil {
		h.Limit = lo.ToPtr(10)
	} else if *h.Limit <= 0 {
		errs = append(errs, fmt.Errorf("Limitは1以上で指定してください"))
	}

	if h.Offset == nil {
		h.Offset = lo.ToPtr(0)
	} else if *h.Offset < 0 {
		errs = append(errs, fmt.Errorf("Offsetは0以上の値を指定してください"))
	}

	return errors.Join(errs...)
}

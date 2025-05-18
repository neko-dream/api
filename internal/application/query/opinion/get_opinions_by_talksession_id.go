package opinion_query

import (
	"context"
	"errors"
	"fmt"

	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/sort"
	"github.com/samber/lo"
)

type (
	GetOpinionsByTalkSessionQuery interface {
		Execute(context.Context, GetOpinionsByTalkSessionInput) (*GetOpinionsByTalkSessionOutput, error)
	}

	GetOpinionsByTalkSessionInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		UserID        *shared.UUID[user.User]
		SortKey       sort.SortKey
		Limit         *int
		Offset        *int
		IsSeed        bool
	}

	GetOpinionsByTalkSessionOutput struct {
		Opinions   []dto.SwipeOpinion
		TotalCount int
	}
)

func (i *GetOpinionsByTalkSessionInput) Validate() error {
	var err error

	// if !i.SortKey.IsValid() {
	// 	err = errors.Join(err, fmt.Errorf("ソートキーが不正です: %s", i.SortKey))
	// }

	if i.Limit != nil {
		if *i.Limit < 0 {
			err = errors.Join(err, fmt.Errorf("limitは0以上である必要があります: %d", *i.Limit))
		}
		// if *i.Limit > 100 {
		// 	err = errors.Join(err, fmt.Errorf("limitは100以下である必要があります: %d", *i.Limit))
		// }
	} else {
		i.Limit = lo.ToPtr(10)
	}

	if i.Offset != nil {
		if *i.Offset < 0 {
			err = errors.Join(err, fmt.Errorf("offsetは0以上である必要があります: %d",
				*i.Offset))
		}
	} else {
		i.Offset = lo.ToPtr(0)
	}

	return err
}

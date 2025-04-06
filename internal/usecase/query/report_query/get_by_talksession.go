package report_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type GetByTalkSessionQuery interface {
	Execute(ctx context.Context, input GetByTalkSessionInput) (*GetByTalkSessionOutput, error)
}

type GetByTalkSessionInput struct {
	TalkSessionID shared.UUID[talksession.TalkSession]
	UserID        shared.UUID[user.User]
	Status        string
}

type GetByTalkSessionOutput struct {
	Reports []dto.ReportDetail
}

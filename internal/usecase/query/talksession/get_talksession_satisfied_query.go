package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/service"
	"go.opentelemetry.io/otel"
)

type (
	// TalkSessionへの参加資格を満たしているかどうかを判定するクエリ
	IsTalkSessionSatisfiedQuery interface {
		Execute(context.Context, IsTalkSessionSatisfiedInput) (*IsTalkSessionSatisfiedOutput, error)
	}

	IsTalkSessionSatisfiedInput struct {
		// TalkSessionID
		TalkSessionID shared.UUID[talksession.TalkSession]
		// UserID
		UserID shared.UUID[user.User]
	}

	IsTalkSessionSatisfiedOutput struct {
		Attributes []talksession.RestrictionAttribute
	}

	isTalkSessionSatisfiedInteractor struct {
		restrictionService service.TalkSessionAccessControl
	}
)

func NewIsTalkSessionSatisfiedInteractor(restrictionService service.TalkSessionAccessControl) IsTalkSessionSatisfiedQuery {
	return &isTalkSessionSatisfiedInteractor{
		restrictionService: restrictionService,
	}
}

func (i *isTalkSessionSatisfiedInteractor) Execute(ctx context.Context, input IsTalkSessionSatisfiedInput) (*IsTalkSessionSatisfiedOutput, error) {
	ctx, span := otel.Tracer("talksession").Start(ctx, "isTalkSessionSatisfiedInteractor.Execute")
	defer span.End()

	attributes, err := i.restrictionService.UserSatisfiesRestriction(ctx, input.TalkSessionID, input.UserID)
	if err != nil {
		return nil, err
	}

	return &IsTalkSessionSatisfiedOutput{
		Attributes: attributes,
	}, nil
}

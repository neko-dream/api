package talksession

import (
	"context"

	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
)

type (
	GetRestrictionsQuery interface {
		Execute(context.Context) (*GetRestrictionsOutput, error)
	}

	GetRestrictionsOutput struct {
		Restrictions []talksession.RestrictionAttribute
	}
)

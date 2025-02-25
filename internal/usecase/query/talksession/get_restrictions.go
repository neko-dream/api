package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/talksession"
)

type (
	GetRestrictionsQuery interface {
		Execute(context.Context) (*GetRestrictionsOutput, error)
	}

	GetRestrictionsOutput struct {
		Restrictions []talksession.RestrictionAttribute
	}
)

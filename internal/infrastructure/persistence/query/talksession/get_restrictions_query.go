package talksession_query

import (
	"context"

	ts "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/usecase/query/talksession"
)

type getRestrictionsQuery struct {
}

func NewGetRestrictionsQuery() talksession.GetRestrictionsQuery {
	return &getRestrictionsQuery{}
}

func (q *getRestrictionsQuery) Execute(ctx context.Context) (*talksession.GetRestrictionsOutput, error) {

	restrictionAttributeMap := ts.RestrictionAttributeKeyMap

	var restrictions []ts.RestrictionAttribute
	for _, restriction := range restrictionAttributeMap {
		restrictions = append(restrictions, restriction)
	}

	return &talksession.GetRestrictionsOutput{
		Restrictions: restrictions,
	}, nil
}

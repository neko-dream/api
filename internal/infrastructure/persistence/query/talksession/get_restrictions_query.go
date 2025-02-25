package talksession_query

import (
	"context"

	ts "github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/usecase/query/talksession"
	"go.opentelemetry.io/otel"
)

type getRestrictionsQuery struct {
}

func NewGetRestrictionsQuery() talksession.GetRestrictionsQuery {
	return &getRestrictionsQuery{}
}

func (q *getRestrictionsQuery) Execute(ctx context.Context) (*talksession.GetRestrictionsOutput, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "getRestrictionsQuery.Execute")
	defer span.End()

	_ = ctx

	restrictionAttributeMap := ts.RestrictionAttributeKeyMap

	var restrictions []ts.RestrictionAttribute
	for _, restriction := range restrictionAttributeMap {
		restrictions = append(restrictions, restriction)
	}

	return &talksession.GetRestrictionsOutput{
		Restrictions: restrictions,
	}, nil
}

package user_query

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	user_query "github.com/neko-dream/server/internal/usecase/query/user"
	"go.opentelemetry.io/otel"
)

type DetailHandler struct {
	*db.DBManager
}

func NewDetailHandler(dbManager *db.DBManager) user_query.Detail {
	return &DetailHandler{
		DBManager: dbManager,
	}
}

func (d *DetailHandler) Execute(ctx context.Context, input user_query.DetailInput) (*user_query.DetailOutput, error) {
	ctx, span := otel.Tracer("user_query").Start(ctx, "DetailHandler.Execute")
	defer span.End()

	userRow, err := d.GetQueries(ctx).GetUserDetailByID(ctx, input.UserID.UUID())
	if err != nil {
		return nil, err
	}

	var userDetail dto.UserDetail
	if err := copier.CopyWithOption(&userDetail, userRow, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return nil, err
	}

	return &user_query.DetailOutput{
		User: userDetail,
	}, nil
}

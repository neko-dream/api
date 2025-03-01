package user_query

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/domain/model/crypto"
	crypto_infra "github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	user_query "github.com/neko-dream/server/internal/usecase/query/user"
	"go.opentelemetry.io/otel"
)

type DetailHandler struct {
	*db.DBManager
	encryptor crypto.Encryptor
}

func NewDetailHandler(dbManager *db.DBManager, encryptor crypto.Encryptor) user_query.Detail {
	return &DetailHandler{
		DBManager: dbManager,
		encryptor: encryptor,
	}
}

func (d *DetailHandler) Execute(ctx context.Context, input user_query.DetailInput) (*user_query.DetailOutput, error) {
	ctx, span := otel.Tracer("user_query").Start(ctx, "DetailHandler.Execute")
	defer span.End()

	userRow, err := d.GetQueries(ctx).GetUserDetailByID(ctx, input.UserID.UUID())
	if err != nil {
		return nil, err
	}

	userDemographic, err := crypto_infra.DecryptUserDemographicsDTO(ctx, d.encryptor, &userRow.UserDemographic)
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
	if userDemographic != nil {
		userDetail.UserDemographic = userDemographic
	}

	return &user_query.DetailOutput{
		User: userDetail,
	}, nil
}

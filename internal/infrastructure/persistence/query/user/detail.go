package user_query

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	user_query "github.com/neko-dream/server/internal/application/query/user"
	"github.com/neko-dream/server/internal/domain/model/crypto"
	crypto_infra "github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
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
	if userRow.User.Email.Valid {
		decrypted, err := d.encryptor.DecryptString(ctx, userRow.User.Email.String)
		if err != nil {
			utils.HandleError(ctx, err, "メアドの復号に失敗しました。")
			userDetail.UserAuth.Email = nil
			userDetail.UserAuth.EmailVerified = false
			return &user_query.DetailOutput{
				User: userDetail,
			}, nil
		}
		userDetail.UserAuth.Email = lo.ToPtr(decrypted)
		userDetail.UserAuth.EmailVerified = userRow.User.EmailVerified
	}

	return &user_query.DetailOutput{
		User: userDetail,
	}, nil
}

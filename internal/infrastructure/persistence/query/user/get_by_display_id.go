package user_query

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/api/internal/application/query/dto"
	user_query "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type GetByDisplayIDHandler struct {
	*db.DBManager
	userRepository user.UserRepository
}

func NewGetByDisplayIDHandler(dbManager *db.DBManager, userRepository user.UserRepository) user_query.GetByDisplayID {
	return &GetByDisplayIDHandler{
		DBManager:      dbManager,
		userRepository: userRepository,
	}
}

func (h *GetByDisplayIDHandler) Execute(ctx context.Context, input user_query.GetByDisplayIDInput) (*user_query.GetByDisplayIDOutput, error) {
	ctx, span := otel.Tracer("user_query").Start(ctx, "GetByDisplayIDHandler.Execute")
	defer span.End()

	user, err := h.userRepository.FindByDisplayID(ctx, input.DisplayID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &user_query.GetByDisplayIDOutput{
			User: nil,
		}, nil
	}

	var userDTO dto.User
	if err := copier.CopyWithOption(&userDTO, user, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return nil, err
	}

	return &user_query.GetByDisplayIDOutput{
		User: &userDTO,
	}, nil
}

package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	GetUserInformationQueryHandler interface {
		Execute(context.Context, GetUserInformationQuery) (*GetUserInformationOutput, error)
	}

	GetUserInformationQuery struct {
		UserID shared.UUID[user.User]
	}

	GetUserInformationOutput struct {
		User user.User
	}

	getUserInformationQueryHandler struct {
		*db.DBManager
		user.UserRepository
	}
)

func NewGetUserInformationQueryHandler(
	dbManager *db.DBManager,
	userRepo user.UserRepository,
) GetUserInformationQueryHandler {
	return &getUserInformationQueryHandler{
		DBManager:      dbManager,
		UserRepository: userRepo,
	}
}

// Execute implements GetUserInformationQueryHandler.
func (g *getUserInformationQueryHandler) Execute(ctx context.Context, input GetUserInformationQuery) (*GetUserInformationOutput, error) {
	u, err := g.UserRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &GetUserInformationOutput{
		User: *u,
	}, nil
}

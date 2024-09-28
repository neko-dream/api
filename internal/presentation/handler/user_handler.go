package handler

import (
	"context"
	"log"

	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
)

type userHandler struct {
}

// EditUserProfile implements oas.UserHandler.
func (u *userHandler) EditUserProfile(ctx context.Context) (*oas.EditUserProfileOK, error) {
	panic("unimplemented")
}

// GetUserProfile implements oas.UserHandler.
func (u *userHandler) GetUserProfile(ctx context.Context) (*oas.GetUserProfileOK, error) {
	panic("unimplemented")
}

// RegisterUser implements oas.UserHandler.
func (u *userHandler) RegisterUser(ctx context.Context, params oas.RegisterUserParams) (oas.RegisterUserRes, error) {
	claim := session.GetSession(ctx)
	log.Println(claim)

	return &oas.RegisterUserOK{
		DisplayID:   "string",
		DisplayName: "string",
	}, nil

}

func NewUserHandler() oas.UserHandler {
	return &userHandler{}
}

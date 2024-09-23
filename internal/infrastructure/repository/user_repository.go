package repository

import (
	"context"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type userRepository struct {
}

// Create implements user.UserRepository.
func (u *userRepository) Create(context.Context, user.User) (user.User, error) {
	panic("unimplemented")
}

// FindByID implements user.UserRepository.
func (u *userRepository) FindByID(context.Context, shared.UUID[user.User]) (*user.User, error) {
	panic("unimplemented")
}

// FindBySubject implements user.UserRepository.
func (u *userRepository) FindBySubject(context.Context, user.UserSubject) (*user.User, error) {
	panic("unimplemented")
}

func NewUserRepository() user.UserRepository {
	return &userRepository{}
}

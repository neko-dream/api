package service

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/user"
)

type userService struct {
	userRepo user.UserRepository
}

// DisplayIDCheckDuplicate ユーザーの表示用IDが重複していないかチェック
func (s *userService) DisplayIDCheckDuplicate(ctx context.Context, displayID string) (bool, error) {
	foundUser, err := s.userRepo.FindByDisplayID(ctx, displayID)
	if err != nil {
		return true, err
	}

	return foundUser != nil, nil
}

func NewUserService(userRepo user.UserRepository) user.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

package service

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type userService struct {
	userRepo user.UserRepository
}

// DisplayIDCheckDuplicate ユーザーの表示用IDが重複していないかチェック
func (s *userService) DisplayIDCheckDuplicate(ctx context.Context, displayID string) (bool, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "userService.DisplayIDCheckDuplicate")
	defer span.End()

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

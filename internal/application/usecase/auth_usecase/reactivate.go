package auth_usecase

import (
	"context"
	"errors"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	Reactivate interface {
		Execute(context.Context, ReactivateInput) (*ReactivateOutput, error)
	}

	ReactivateInput struct {
		UserID shared.UUID[user.User]
	}

	ReactivateOutput struct {
		Message string
		User    dto.User
	}

	reactivate struct {
		userRepository                user.UserRepository
		tokenManager                  session.TokenManager
		userStatusChangeLogRepository user.UserStatusChangeLogRepository
	}
)

func NewReactivate(
	userRepository user.UserRepository,
	tokenManager session.TokenManager,
	userStatusChangeLogRepository user.UserStatusChangeLogRepository,
) Reactivate {
	return &reactivate{
		userRepository:                userRepository,
		tokenManager:                  tokenManager,
		userStatusChangeLogRepository: userStatusChangeLogRepository,
	}
}

func (r *reactivate) Execute(ctx context.Context, input ReactivateInput) (*ReactivateOutput, error) {
	ctx, span := otel.Tracer("usecase").Start(ctx, "ReactivateUseCase.Execute")
	defer span.End()

	// ユーザーを取得
	u, err := r.userRepository.FindByID(ctx, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "userRepository.FindByID")
		return nil, errtrace.Wrap(err)
	}
	if u == nil {
		return nil, messages.UserNotFound
	}

	// ユーザーを復活させる
	now := time.Now()
	if err := u.Reactivate(now); err != nil {
		// ドメインエラーをアプリケーション層のエラーに変換
		if errors.Is(err, user.ErrNotWithdrawn) {
			return nil, messages.UserNotWithdrawn
		} else if errors.Is(err, user.ErrReactivationPeriodExpired) {
			return nil, messages.UserReactivationPeriodExpired
		}
		utils.HandleError(ctx, err, "u.Reactivate")
		return nil, errtrace.Wrap(err)
	}

	// ユーザー情報を更新
	if err := r.userRepository.Update(ctx, *u); err != nil {
		utils.HandleError(ctx, err, "userRepository.Update")
		return nil, errtrace.Wrap(err)
	}

	// UserStatusChangeLogへの記録
	statusChangeLog := user.NewUserStatusChangeLog(
		input.UserID,
		user.UserStatusReactivated,
		nil, // 復活理由は特に記録しない
		time.Now(),
		user.ChangedByUser,
		nil, // IPアドレスは今後コンテキストから取得
		nil, // UserAgentは今後コンテキストから取得
	)
	if err := r.userStatusChangeLogRepository.Create(ctx, statusChangeLog); err != nil {
		// ログへの記録失敗は復活処理自体は成功とする
		utils.HandleError(ctx, err, "userStatusChangeLogRepository.Create")
	}

	output := &ReactivateOutput{
		Message: "アカウントが復活しました",
		User: dto.User{
			DisplayID:   *u.DisplayID(),
			DisplayName: *u.DisplayName(),
			IconURL:     u.IconURL(),
		},
	}

	return output, nil
}

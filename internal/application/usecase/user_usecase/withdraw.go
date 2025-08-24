package user_usecase

import (
	"context"
	"errors"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	Withdraw interface {
		Execute(context.Context, WithdrawInput) (*WithdrawOutput, error)
	}

	WithdrawInput struct {
		UserID shared.UUID[user.User]
	}

	WithdrawOutput struct {
		Message        string
		WithdrawalDate time.Time
	}

	withdraw struct {
		userRepository                user.UserRepository
		sessionRepository             session.SessionRepository
		userStatusChangeLogRepository user.UserStatusChangeLogRepository
	}
)

func NewWithdraw(
	userRepository user.UserRepository,
	sessionRepository session.SessionRepository,
	userStatusChangeLogRepository user.UserStatusChangeLogRepository,
) Withdraw {
	return &withdraw{
		userRepository:                userRepository,
		sessionRepository:             sessionRepository,
		userStatusChangeLogRepository: userStatusChangeLogRepository,
	}
}

func (w *withdraw) Execute(ctx context.Context, input WithdrawInput) (*WithdrawOutput, error) {
	ctx, span := otel.Tracer("usecase").Start(ctx, "WithdrawUseCase.Execute")
	defer span.End()

	// ユーザーを取得
	u, err := w.userRepository.FindByID(ctx, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "userRepository.FindByID")
		return nil, errtrace.Wrap(err)
	}
	if u == nil {
		return nil, messages.UserNotFound
	}

	// 退会処理
	now := time.Now()
	if err := u.Withdraw(now); err != nil {
		if errors.Is(err, user.ErrAlreadyWithdrawn) {
			return nil, messages.UserAlreadyWithdrawn
		}
		utils.HandleError(ctx, err, "u.Withdraw")
		return nil, errtrace.Wrap(err)
	}

	// ユーザー情報を更新
	if err := w.userRepository.Update(ctx, *u); err != nil {
		utils.HandleError(ctx, err, "userRepository.Update")
		return nil, errtrace.Wrap(err)
	}

	// すべてのセッションを無効化
	if err := w.sessionRepository.DeactivateAllByUserID(ctx, input.UserID); err != nil {
		utils.HandleError(ctx, err, "sessionRepository.DeactivateAllByUserID")
		return nil, errtrace.Wrap(err)
	}

	// UserStatusChangeLogへの記録
	statusChangeLog := user.NewUserStatusChangeLog(
		input.UserID,
		user.UserStatusWithdrawn,
		nil, // reasonは削除
		now,
		user.ChangedByUser,
		nil, // IPアドレスは今後コンテキストから取得
		nil, // UserAgentは今後コンテキストから取得
	)
	if err := w.userStatusChangeLogRepository.Create(ctx, statusChangeLog); err != nil {
		// ログへの記録失敗は退会処理自体は成功とする
		utils.HandleError(ctx, err, "userStatusChangeLogRepository.Create")
	}

	return &WithdrawOutput{
		Message:        "退会処理が完了しました。30日以内であれば再ログインで復活できます。",
		WithdrawalDate: now,
	}, nil
}

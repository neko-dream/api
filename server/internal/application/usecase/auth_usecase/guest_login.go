package auth_usecase

import (
	"context"
	"log"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/consent"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	GuestLogin interface {
		Execute(context.Context, GuestLoginInput) (GuestLoginOutput, error)
	}

	GuestLoginInput struct {
		// 将来的な拡張用（現時点では空）
	}

	GuestLoginOutput struct {
		Token string
	}

	guestLoginInteractor struct {
		*db.DBManager
		session.SessionRepository
		session.SessionService
		session.TokenManager
		user.UserRepository
		consent.ConsentService
		consent.PolicyRepository
	}
)

func NewGuestLogin(
	tm *db.DBManager,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
	userRepository user.UserRepository,
	consentService consent.ConsentService,
	policyRepository consent.PolicyRepository,
) GuestLogin {
	return &guestLoginInteractor{
		DBManager:         tm,
		SessionRepository: sessionRepository,
		SessionService:    sessionService,
		TokenManager:      tokenManager,
		UserRepository:    userRepository,
		ConsentService:    consentService,
		PolicyRepository:  policyRepository,
	}
}

func (g *guestLoginInteractor) Execute(ctx context.Context, input GuestLoginInput) (GuestLoginOutput, error) {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "guestLoginInteractor.Execute")
	defer span.End()

	var tok string

	if err := g.ExecTx(ctx, func(ctx context.Context) error {
		// ゲストユーザー用のランダム情報を生成
		displayID := user.GenerateGuestDisplayID()
		displayName := user.GenerateGuestDisplayName()
		subject := "guest_" + shared.NewUUID[user.User]().String() // 一意なsubject

		// ゲストユーザーを作成
		newUser := user.NewUser(
			shared.NewUUID[user.User](),
			&displayID,
			&displayName,
			subject,
			shared.ProviderGuest,
			nil,
		)
		newUser.SetDemographics(user.NewUserDemographic(
			ctx,
			shared.NewUUID[user.UserDemographic](),
			nil, nil, nil, nil,
		))
		log.Println("Created guest user:", newUser.UserID().String())

		// ポリシー同意を記録
		version, err := g.PolicyRepository.FetchLatestPolicy(ctx)
		if err != nil {
			utils.HandleError(ctx, err, "PolicyRepository.FetchLatestPolicy")
			return errtrace.Wrap(err)
		}
		if _, err := g.ConsentService.RecordConsent(ctx, newUser.UserID(), version.Version, "", ""); err != nil {
			utils.HandleError(ctx, err, "ConsentService.RecordConsent")
			return errtrace.Wrap(err)
		}

		// ユーザーを永続化
		if err := g.UserRepository.Create(ctx, newUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Create")
			return errtrace.Wrap(err)
		}

		// セッション作成（ゲストは30日間有効）
		sess := session.NewSession(
			shared.NewUUID[session.Session](),
			newUser.UserID(),
			newUser.Provider(),
			session.SESSION_ACTIVE,
			*session.NewGuestExpiresAt(ctx),
			clock.Now(ctx),
		)

		if _, err := g.SessionRepository.Create(ctx, *sess); err != nil {
			utils.HandleError(ctx, err, "SessionRepository.Create")
			return errtrace.Wrap(err)
		}

		// トークン生成
		token, err := g.TokenManager.Generate(ctx, newUser, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "TokenManager.Generate")
			return errtrace.Wrap(err)
		}

		tok = token
		return nil
	}); err != nil {
		return GuestLoginOutput{}, errtrace.Wrap(err)
	}

	return GuestLoginOutput{Token: tok}, nil
}

package auth_command

import (
	"context"
	"log"
	"regexp"

	"github.com/neko-dream/server/internal/domain/messages"
	password_auth "github.com/neko-dream/server/internal/domain/model/auth/password"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/hash"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type PasswordLogin interface {
	Execute(ctx context.Context, input PasswordLoginInput) (*PasswordLoginOutput, error)
}

type PasswordLoginInput struct {
	IDorEmail string
	Password  string
}

type PasswordLoginOutput struct {
	Token string
}

type passwordLoginInteractor struct {
	user.UserRepository
	session.SessionService
	session.SessionRepository
	*db.DBManager
	*config.Config
	password_auth.PasswordAuthManager
	session.TokenManager
}

func NewPasswordLogin(
	userRep user.UserRepository,
	sessService session.SessionService,
	sessionRepository session.SessionRepository,
	tm *db.DBManager,
	config *config.Config,
	passwordAuthManager password_auth.PasswordAuthManager,
	tokenManager session.TokenManager,
) PasswordLogin {
	return &passwordLoginInteractor{
		UserRepository:      userRep,
		SessionService:      sessService,
		SessionRepository:   sessionRepository,
		DBManager:           tm,
		Config:              config,
		PasswordAuthManager: passwordAuthManager,
		TokenManager:        tokenManager,
	}
}

func IsEmail(s string) bool {
	// メールアドレスの正規表現
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(s)
}

func (p *passwordLoginInteractor) Execute(ctx context.Context, input PasswordLoginInput) (*PasswordLoginOutput, error) {
	ctx, span := otel.Tracer("auth_command").Start(ctx, "passwordLoginInteractor.Execute")
	defer span.End()

	var tokenRes string
	if err := p.ExecTx(ctx, func(ctx context.Context) error {
		var usr *user.User
		// IDorEmailがメールアドレスの場合は、ユーザーをメールアドレスで取得
		if IsEmail(input.IDorEmail) {
			emailHash, err := hash.HashEmail(input.IDorEmail, p.Config.HASH_PEPPER)
			if err != nil {
				utils.HandleError(ctx, err, "failed to hash email")
				return messages.InvalidPasswordOrEmailError
			}
			foundUser, err := p.UserRepository.FindBySubject(ctx, user.UserSubject(emailHash))
			if err != nil {
				log.Println("failed to find user by subject", err)
				return messages.InvalidPasswordOrEmailError
			}
			usr = foundUser
		} else {
			// DisplayIDからユーザーを取得
			foundUser, err := p.UserRepository.FindByDisplayID(ctx, input.IDorEmail)
			if err != nil {
				return messages.InvalidPasswordOrEmailError
			}
			usr = foundUser
		}
		if usr == nil {
			log.Println("user not found")
			return messages.InvalidPasswordOrEmailError
		}

		// ユーザーからユーザーパスワードを取得
		userPassword, err := p.PasswordAuthManager.VerifyPassword(ctx, usr.UserID(), input.Password)
		if err != nil {
			return messages.InvalidPasswordOrEmailError
		}
		// パスワードが一致しない場合はエラーを返す
		if !userPassword {
			log.Println("password mismatch")
			return messages.InvalidPasswordOrEmailError
		}

		if err := p.SessionService.DeactivateUserSessions(ctx, usr.UserID()); err != nil {
			utils.HandleError(ctx, err, "failed to deactivate user sessions")
			return err
		}

		sess := session.NewSession(
			shared.NewUUID[session.Session](),
			usr.UserID(),
			usr.Provider(),
			session.SESSION_ACTIVE,
			*session.NewExpiresAt(ctx),
			clock.Now(ctx),
		)
		if _, err := p.SessionRepository.Create(ctx, *sess); err != nil {
			utils.HandleError(ctx, err, "failed to create session")
			return err
		}

		token, err := p.TokenManager.Generate(ctx, *usr, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return err
		}

		tokenRes = token
		return nil
	}); err != nil {
		return nil, err
	}

	// トークンを生成
	return &PasswordLoginOutput{
		Token: tokenRes,
	}, nil
}

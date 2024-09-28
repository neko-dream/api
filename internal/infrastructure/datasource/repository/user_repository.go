package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"braces.dev/errtrace"
	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/oauth"
)

type userRepository struct {
	*db.DBManager
}

func (u *userRepository) Create(ctx context.Context, user user.User) (user.User, error) {

	if err := u.GetQueries(ctx).CreateUser(ctx, model.CreateUserParams{
		UserID:      user.UserID().UUID(),
		DisplayName: user.DisplayName(),
		DisplayID:   user.DisplayID(),
		CreatedAt:   time.Now(),
	}); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return user, errtrace.Wrap(err)
		}
	}

	if err := u.GetQueries(ctx).CreateUserAuth(ctx, model.CreateUserAuthParams{
		UserID:    uuid.NullUUID{UUID: user.UserID().UUID(), Valid: true},
		Provider:  strings.ToUpper(user.Provider().String()),
		Subject:   user.Subject(),
		CreatedAt: time.Now(),
	}); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return user, errtrace.Wrap(err)
		}
	}

	return user, nil
}

// FindByID implements user.UserRepository.
func (u *userRepository) FindByID(ctx context.Context, userID shared.UUID[user.User]) (*user.User, error) {
	panic("implement me")
}

func (u *userRepository) FindBySubject(ctx context.Context, subject user.UserSubject) (*user.User, error) {
	row, err := u.GetQueries(ctx).GetUserBySubject(ctx, subject.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errtrace.Wrap(err)
	}

	providerName, err := oauth.NewAuthProviderName(row.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var picture *string
	if row.Picture.Valid {
		picture = &row.Picture.String
	}

	user := user.NewUser(
		shared.MustParseUUID[user.User](row.UserID.String()),
		row.DisplayID,
		row.DisplayName,
		row.Subject,
		providerName,
		picture,
	)

	return &user, nil
}

func NewUserRepository(
	DBManager *db.DBManager,
) user.UserRepository {
	return &userRepository{
		DBManager: DBManager,
	}
}

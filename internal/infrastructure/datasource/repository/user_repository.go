package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	um "github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/oauth"
	"github.com/samber/lo"
)

type userRepository struct {
	*db.DBManager
}

// FindByDisplayID ユーザーのDisplayIDを元にユーザーを取得する
func (u *userRepository) FindByDisplayID(ctx context.Context, displayID string) (*um.User, error) {
	userRow, err := u.
		GetQueries(ctx).
		UserFindByDisplayID(ctx, sql.NullString{String: displayID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errtrace.Wrap(err)
	}

	userAuthRow, err := u.
		GetQueries(ctx).
		GetUserAuthByUserID(ctx, userRow.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	providerName, err := oauth.
		NewAuthProviderName(userAuthRow.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var picture *string
	if userRow.Picture.Valid {
		picture = &userRow.Picture.String
	}
	user := user.NewUser(
		shared.MustParseUUID[user.User](userRow.UserID.String()),
		lo.ToPtr(displayID),
		lo.ToPtr(userRow.DisplayName.String),
		userAuthRow.Subject,
		providerName,
		picture,
	)

	return &user, nil
}

// Update ユーザー情報を更新する
func (u *userRepository) Update(ctx context.Context, user um.User) error {
	var picture, displayID, displayName sql.NullString
	if user.Picture() != nil {
		picture = sql.NullString{String: *user.Picture(), Valid: true}
	}
	if user.DisplayID() != nil {
		displayID = sql.NullString{String: *user.DisplayID(), Valid: true}
	}
	if user.DisplayName() != nil {
		displayName = sql.NullString{String: *user.DisplayName(), Valid: true}
	}

	if err := u.DBManager.GetQueries(ctx).UpdateUser(ctx, model.UpdateUserParams{
		UserID:      user.UserID().UUID(),
		DisplayID:   displayID,
		DisplayName: displayName,
		Picture:     picture,
	}); err != nil {
		return errtrace.Wrap(err)
	}

	if user.Demographics() != nil {
		var (
			yearOfBirth   sql.NullInt32
			occupation    sql.NullInt16
			municipality  sql.NullString
			householdSize sql.NullInt16
			gender        int16
		)
		if user.Demographics().YearOfBirth() != nil {
			yearOfBirth = sql.NullInt32{Int32: int32(*user.Demographics().YearOfBirth()), Valid: true}
		}
		if user.Demographics().Occupation() != nil {
			occupation = sql.NullInt16{Int16: int16(*user.Demographics().Occupation()), Valid: true}
		}
		if user.Demographics().Municipality() != nil {
			municipality = sql.NullString{String: user.Demographics().Municipality().String(), Valid: true}
		}
		if user.Demographics().HouseholdSize() != nil {
			householdSize = sql.NullInt16{Int16: int16(*user.Demographics().HouseholdSize()), Valid: true}
		}
		if user.Demographics().Gender() == nil {
			gender = int16(um.GenderPreferNotToSay)
		}

		if err := u.DBManager.GetQueries(ctx).
			UpdateOrCreateUserDemographics(ctx, model.UpdateOrCreateUserDemographicsParams{
				UserDemographicsID: shared.NewUUID[um.UserDemographics]().UUID(),
				UserID:             user.UserID().UUID(),
				YearOfBirth:        yearOfBirth,
				Occupation:         occupation,
				Municipality:       municipality,
				HouseholdSize:      householdSize,
				Gender:             gender,
			}); err != nil {
			return errtrace.Wrap(err)
		}
	}

	return nil
}

// Create 初回登録時は必ずDisplayID, DisplayName, Pictureが空文字列で登録される
// また、UserAuthはIsVerifyがfalseで登録される
func (u *userRepository) Create(ctx context.Context, usr user.User) error {
	if err := u.GetQueries(ctx).CreateUser(ctx, model.CreateUserParams{
		UserID:    usr.UserID().UUID(),
		CreatedAt: time.Now(),
	}); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errtrace.Wrap(err)
		}
	}

	if err := u.GetQueries(ctx).CreateUserAuth(ctx, model.CreateUserAuthParams{
		UserAuthID: shared.NewUUID[user.User]().UUID(),
		UserID:     usr.UserID().UUID(),
		Provider:   strings.ToUpper(usr.Provider().String()),
		Subject:    usr.Subject(),
		CreatedAt:  time.Now(),
	}); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

// FindByID ユーザーIDを元にユーザーを取得する
func (u *userRepository) FindByID(ctx context.Context, userID shared.UUID[user.User]) (*user.User, error) {
	userRow, err := u.GetQueries(ctx).GetUserByID(ctx, userID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errtrace.Wrap(err)
	}
	userAuthRow, err := u.GetQueries(ctx).GetUserAuthByUserID(ctx, userID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	providerName, err := oauth.NewAuthProviderName(userAuthRow.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var picture, displayID, displayName *string
	if userRow.Picture.Valid {
		picture = &userRow.Picture.String
	}
	if userRow.DisplayID.Valid {
		displayID = &userRow.DisplayID.String
	}
	if userRow.DisplayName.Valid {
		displayName = &userRow.DisplayName.String
	}

	user := user.NewUser(
		shared.MustParseUUID[user.User](userRow.UserID.String()),
		displayID,
		displayName,
		userAuthRow.Subject,
		providerName,
		picture,
	)

	return &user, nil
}

func (u *userRepository) FindBySubject(ctx context.Context, subject user.UserSubject) (*user.User, error) {
	row, err := u.GetQueries(ctx).GetUserBySubject(ctx, subject.String())
	if err != nil {
		return nil, nil
	}

	providerName, err := oauth.NewAuthProviderName(row.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var picture, displayID, displayName *string
	if row.Picture.Valid {
		picture = &row.Picture.String
	}
	if row.DisplayID.Valid {
		displayID = &row.DisplayID.String
	}
	if row.DisplayName.Valid {
		displayName = &row.DisplayName.String
	}

	user := user.NewUser(
		shared.MustParseUUID[user.User](row.UserID.String()),
		displayID,
		displayName,
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

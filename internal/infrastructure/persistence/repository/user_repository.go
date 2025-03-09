package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/crypto"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	um "github.com/neko-dream/server/internal/domain/model/user"
	ci "github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type userRepository struct {
	*db.DBManager
	imageStorage image.ImageStorage
	encryptor    crypto.Encryptor
}

func NewUserRepository(
	DBManager *db.DBManager,
	imageStorage image.ImageStorage,
	encryptor crypto.Encryptor,
) user.UserRepository {
	return &userRepository{
		DBManager:    DBManager,
		imageStorage: imageStorage,
		encryptor:    encryptor,
	}
}

// FindByDisplayID ユーザーのDisplayIDを元にユーザーを取得する
func (u *userRepository) FindByDisplayID(ctx context.Context, displayID string) (*um.User, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "userRepository.FindByDisplayID")
	defer span.End()

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

	providerName, err := auth.NewAuthProviderName(userAuthRow.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var iconURL *string
	if userRow.IconUrl.Valid {
		iconURL = &userRow.IconUrl.String
	}

	userID, err := shared.ParseUUID[user.User](userRow.UserID.String())
	if err != nil {
		return nil, err
	}

	user := user.NewUser(
		userID,
		lo.ToPtr(displayID),
		lo.ToPtr(userRow.DisplayName.String),
		userAuthRow.Subject,
		providerName,
		iconURL,
	)

	return &user, nil
}

// Update ユーザー情報を更新する
func (u *userRepository) Update(ctx context.Context, user um.User) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "userRepository.Update")
	defer span.End()

	var displayID, displayName, iconURL sql.NullString
	if user.DisplayID() != nil {
		displayID = sql.NullString{String: *user.DisplayID(), Valid: true}
	}
	if user.DisplayName() != nil {
		displayName = sql.NullString{String: *user.DisplayName(), Valid: true}
	}
	if user.IconURL() != nil {
		iconURL = sql.NullString{String: *user.IconURL(), Valid: true}
	}

	if err := u.GetQueries(ctx).UpdateUser(ctx, model.UpdateUserParams{
		UserID:      user.UserID().UUID(),
		DisplayName: displayName,
		DisplayID:   displayID,
		IconUrl:     iconURL,
	}); err != nil {
		utils.HandleError(ctx, err, "UpdateUser")
		return errtrace.Wrap(err)
	}

	if user.Demographics() != nil {
		encryptedDemo, err := ci.EncryptUserDemographics(ctx, u.encryptor, user.UserID(), user.Demographics())
		if err != nil {
			return errtrace.Wrap(err)
		}

		if err := u.GetQueries(ctx).UpdateOrCreateUserDemographic(ctx, model.UpdateOrCreateUserDemographicParams{
			UserDemographicsID: encryptedDemo.UserDemographicsID,
			UserID:             encryptedDemo.UserID,
			YearOfBirth:        encryptedDemo.YearOfBirth,
			City:               encryptedDemo.City,
			Gender:             encryptedDemo.Gender,
			Prefecture:         encryptedDemo.Prefecture,
		}); err != nil {
			utils.HandleError(ctx, err, "UpdateOrCreateUserDemographic")
			return errtrace.Wrap(err)
		}
	}

	if err := u.GetQueries(ctx).VerifyUser(ctx, user.UserID().UUID()); err != nil {
		utils.HandleError(ctx, err, "VerifyUser")
		return errtrace.Wrap(err)
	}

	return nil
}

// Create 初回登録時は必ずDisplayID, DisplayName, Pictureが空文字列で登録される
// また、UserAuthはIsVerifyがfalseで登録される
func (u *userRepository) Create(ctx context.Context, usr user.User) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "userRepository.Create")
	defer span.End()

	if err := u.GetQueries(ctx).CreateUser(ctx, model.CreateUserParams{
		UserID:    usr.UserID().UUID(),
		CreatedAt: clock.Now(ctx),
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
		CreatedAt:  clock.Now(ctx),
	}); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

// FindByID ユーザーIDを元にユーザーを取得する
func (u *userRepository) FindByID(ctx context.Context, userID shared.UUID[user.User]) (*user.User, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "userRepository.FindByID")
	defer span.End()

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
	userDemographic, err := u.findUserDemographic(ctx, userID)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	providerName, err := auth.NewAuthProviderName(userAuthRow.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var displayID, displayName, iconURL *string
	if userRow.DisplayID.Valid {
		displayID = &userRow.DisplayID.String
	}
	if userRow.DisplayName.Valid {
		displayName = &userRow.DisplayName.String
	}
	if userRow.IconUrl.Valid {
		iconURL = &userRow.IconUrl.String
	}

	user := user.NewUser(
		userID,
		displayID,
		displayName,
		userAuthRow.Subject,
		providerName,
		iconURL,
	)
	if userDemographic != nil {
		user.SetDemographics(*userDemographic)
	}

	return &user, nil
}

func (u *userRepository) findUserDemographic(ctx context.Context, userID shared.UUID[user.User]) (*user.UserDemographic, error) {
	userDemoRow, err := u.GetQueries(ctx).GetUserDemographicByUserID(ctx, userID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.HandleError(ctx, err, "GetUserDemographicByUserID")
		return nil, errtrace.Wrap(err)
	}

	decrypted, err := ci.DecryptUserDemographics(ctx, u.encryptor, &userDemoRow)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return decrypted, nil
}

func (u *userRepository) FindBySubject(ctx context.Context, subject user.UserSubject) (*user.User, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "userRepository.FindBySubject")
	defer span.End()

	row, err := u.GetQueries(ctx).GetUserBySubject(ctx, subject.String())
	if err != nil {
		return nil, nil
	}

	providerName, err := auth.NewAuthProviderName(row.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var displayID, displayName, iconURL *string
	if row.DisplayID.Valid {
		displayID = &row.DisplayID.String
	}
	if row.DisplayName.Valid {
		displayName = &row.DisplayName.String
	}
	if row.IconUrl.Valid {
		iconURL = &row.IconUrl.String
	}

	user := user.NewUser(
		shared.MustParseUUID[user.User](row.UserID.String()),
		displayID,
		displayName,
		row.Subject,
		providerName,
		iconURL,
	)

	return &user, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	um "github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type userRepository struct {
	*db.DBManager
	imageStorage image.ImageStorage
}

func NewUserRepository(
	DBManager *db.DBManager,
	imageStorage image.ImageStorage,
) user.UserRepository {
	return &userRepository{
		DBManager:    DBManager,
		imageStorage: imageStorage,
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
	if err := u.GetQueries(ctx).UpdateUser(ctx, model.UpdateUserParams{
		UserID: user.UserID().UUID(),
		DisplayName: utils.IfThenElse(
			user.DisplayName() != nil,
			sql.NullString{String: *user.DisplayName(), Valid: true},
			sql.NullString{},
		),
		DisplayID: utils.IfThenElse(
			user.DisplayID() != nil,
			sql.NullString{String: *user.DisplayID(), Valid: true},
			sql.NullString{},
		),
		IconUrl: utils.IfThenElse(
			user.IconURL() != nil,
			sql.NullString{String: *user.IconURL(), Valid: true},
			sql.NullString{},
		),
	}); err != nil {
		utils.HandleError(ctx, err, "UpdateUser")
		return errtrace.Wrap(err)
	}

	if user.Demographics() != nil {
		UserDemographic := *user.Demographics()

		var city sql.NullString
		if UserDemographic.City() != nil {
			city = sql.NullString{String: (*UserDemographic.City()).String(), Valid: true}
		}
		var yearOfBirth sql.NullInt32
		if UserDemographic.YearOfBirth() != nil {
			yearOfBirth = sql.NullInt32{Int32: int32(*UserDemographic.YearOfBirth()), Valid: true}
		}
		var householdSize sql.NullInt16
		if UserDemographic.HouseholdSize() != nil {
			householdSize = sql.NullInt16{Int16: int16(*UserDemographic.HouseholdSize()), Valid: true}
		}
		var prefecture sql.NullString
		if UserDemographic.Prefecture() != nil {
			prefecture = sql.NullString{String: *UserDemographic.Prefecture(), Valid: true}
		}
		var occupation sql.NullInt16
		if UserDemographic.Occupation() != nil {
			occupation = sql.NullInt16{Int16: int16(*UserDemographic.Occupation()), Valid: true}
		}
		var gender sql.NullInt16
		if UserDemographic.Gender() != nil {
			gender = sql.NullInt16{Int16: int16(*UserDemographic.Gender()), Valid: true}
		}

		if err := u.GetQueries(ctx).
			UpdateOrCreateUserDemographic(ctx, model.UpdateOrCreateUserDemographicParams{
				UserDemographicsID: UserDemographic.ID().UUID(),
				UserID:             user.UserID().UUID(),
				YearOfBirth:        yearOfBirth,
				Occupation:         occupation,
				City:               city,
				HouseholdSize:      householdSize,
				Gender:             gender,
				Prefecture:         prefecture,
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
	UserDemographic, err := u.findUserDemographic(ctx, userID)
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
	if UserDemographic != nil {
		user.SetDemographics(*UserDemographic)
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

	var (
		yearOfBirth   *int
		occupation    *string
		city          *string
		householdSize *int
		gender        *string
		prefecture    *string
	)
	UserDemographicID, err := shared.ParseUUID[user.UserDemographic](userDemoRow.UserDemographicsID.String())
	if err != nil {
		return nil, err
	}

	if userDemoRow.YearOfBirth.Valid {
		yearOfBirth = lo.ToPtr(int(userDemoRow.YearOfBirth.Int32))
	}
	if userDemoRow.Occupation.Valid {
		occupation = lo.ToPtr(um.Occupation(int(userDemoRow.Occupation.Int16)).String())
	}
	if userDemoRow.Gender.Valid {
		gender = lo.ToPtr(um.Gender(int(userDemoRow.Gender.Int16)).String())
	}
	if userDemoRow.City.Valid {
		city = lo.ToPtr(userDemoRow.City.String)
	}
	if userDemoRow.HouseholdSize.Valid {
		householdSize = lo.ToPtr(int(userDemoRow.HouseholdSize.Int16))
	}
	if userDemoRow.Prefecture.Valid {
		prefecture = lo.ToPtr(userDemoRow.Prefecture.String)
	}

	ud := user.NewUserDemographic(
		ctx,
		UserDemographicID,
		yearOfBirth,
		occupation,
		gender,
		city,
		householdSize,
		prefecture,
	)
	return &ud, nil
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

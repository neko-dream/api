package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	um "github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/oauth"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type userRepository struct {
	*db.DBManager
	imageRepository image.ImageRepository
}

func NewUserRepository(
	DBManager *db.DBManager,
	imageRepository image.ImageRepository,
) user.UserRepository {
	return &userRepository{
		DBManager:       DBManager,
		imageRepository: imageRepository,
	}
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
	var profileIcon *user.ProfileIcon
	if userRow.IconUrl.Valid {
		profileIcon = user.NewProfileIcon(lo.ToPtr(userRow.IconUrl.String))
	}

	user := user.NewUser(
		shared.MustParseUUID[user.User](userRow.UserID.String()),
		lo.ToPtr(displayID),
		lo.ToPtr(userRow.DisplayName.String),
		userAuthRow.Subject,
		providerName,
		profileIcon,
	)

	return &user, nil
}

// Update ユーザー情報を更新する
func (u *userRepository) Update(ctx context.Context, user um.User) error {
	var iconURL sql.NullString
	if user.IsIconUpdateRequired() {
		url, err := u.imageRepository.Create(ctx, *user.ProfileIcon().ImageInfo())
		if err != nil {
			utils.HandleError(ctx, err, "ImageRepository.Create")
			return errtrace.Wrap(err)
		}
		iconURL = sql.NullString{String: *url, Valid: true}
	}

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
		IconUrl: iconURL,
	}); err != nil {
		utils.HandleError(ctx, err, "UpdateUser")
		return errtrace.Wrap(err)
	}

	if user.Demographics() != nil {
		userDemographics := *user.Demographics()

		var city sql.NullString
		if userDemographics.City() != nil {
			city = sql.NullString{String: (*userDemographics.City()).String(), Valid: true}
		}
		var yearOfBirth sql.NullInt32
		if userDemographics.YearOfBirth() != nil {
			yearOfBirth = sql.NullInt32{Int32: int32(*userDemographics.YearOfBirth()), Valid: true}
		}
		var householdSize sql.NullInt16
		if userDemographics.HouseholdSize() != nil {
			householdSize = sql.NullInt16{Int16: int16(*userDemographics.HouseholdSize()), Valid: true}
		}
		var prefecture sql.NullString
		if userDemographics.Prefecture() != nil {
			prefecture = sql.NullString{String: *userDemographics.Prefecture(), Valid: true}
		}

		if err := u.GetQueries(ctx).
			UpdateOrCreateUserDemographics(ctx, model.UpdateOrCreateUserDemographicsParams{
				UserDemographicsID: userDemographics.UserDemographicsID().UUID(),
				UserID:             user.UserID().UUID(),
				YearOfBirth:        yearOfBirth,
				Occupation:         sql.NullInt16{Int16: int16(userDemographics.Occupation()), Valid: true},
				City:               city,
				HouseholdSize:      householdSize,
				Gender:             int16(userDemographics.Gender()),
				Prefecture:         prefecture,
			}); err != nil {
			utils.HandleError(ctx, err, "UpdateOrCreateUserDemographics")
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
	userDemographics, err := u.findUserDemographics(ctx, userID)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	providerName, err := oauth.NewAuthProviderName(userAuthRow.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var displayID, displayName *string
	if userRow.DisplayID.Valid {
		displayID = &userRow.DisplayID.String
	}
	if userRow.DisplayName.Valid {
		displayName = &userRow.DisplayName.String
	}
	var profIcon *user.ProfileIcon
	if userRow.IconUrl.Valid {
		profIcon = user.NewProfileIcon(lo.ToPtr(userRow.IconUrl.String))
	}

	user := user.NewUser(
		shared.MustParseUUID[user.User](userRow.UserID.String()),
		displayID,
		displayName,
		userAuthRow.Subject,
		providerName,
		profIcon,
	)
	if userDemographics != nil {
		user.SetDemographics(*userDemographics)
	}

	return &user, nil
}

func (u *userRepository) findUserDemographics(ctx context.Context, userID shared.UUID[user.User]) (*user.UserDemographics, error) {
	userDemoRow, err := u.GetQueries(ctx).GetUserDemographicsByUserID(ctx, userID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.HandleError(ctx, err, "GetUserDemographicsByUserID")
		return nil, errtrace.Wrap(err)
	}

	var (
		yearOfBirth   *um.YearOfBirth
		occupation    *um.Occupation
		city          *um.City
		householdSize *um.HouseholdSize
		gender        *um.Gender
		prefecture    *string
	)
	userDemographicsID := shared.MustParseUUID[user.UserDemographics](userDemoRow.UserDemographicsID.String())

	if userDemoRow.YearOfBirth.Valid {
		yearOfBirth = um.NewYearOfBirth(lo.ToPtr(int(userDemoRow.YearOfBirth.Int32)))
	}
	if userDemoRow.Occupation.Valid {
		occupation = lo.ToPtr(um.Occupation(int(userDemoRow.Occupation.Int16)))
	}
	gender = lo.ToPtr(um.Gender(int(userDemoRow.Gender)))
	if userDemoRow.City.Valid {
		city = um.NewCity(lo.ToPtr(userDemoRow.City.String))
	}
	if userDemoRow.HouseholdSize.Valid {
		householdSize = um.NewHouseholdSize(lo.ToPtr(int(userDemoRow.HouseholdSize.Int16)))
	}
	if userDemoRow.Prefecture.Valid {
		prefecture = lo.ToPtr(userDemoRow.Prefecture.String)
	}

	ud := user.NewUserDemographics(
		userDemographicsID,
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
	row, err := u.GetQueries(ctx).GetUserBySubject(ctx, subject.String())
	if err != nil {
		return nil, nil
	}

	providerName, err := oauth.NewAuthProviderName(row.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var displayID, displayName *string
	if row.DisplayID.Valid {
		displayID = &row.DisplayID.String
	}
	if row.DisplayName.Valid {
		displayName = &row.DisplayName.String
	}
	var profileIcon *user.ProfileIcon
	if row.IconUrl.Valid {
		profileIcon = user.NewProfileIcon(lo.ToPtr(row.IconUrl.String))
	}

	user := user.NewUser(
		shared.MustParseUUID[user.User](row.UserID.String()),
		displayID,
		displayName,
		row.Subject,
		providerName,
		profileIcon,
	)

	return &user, nil
}

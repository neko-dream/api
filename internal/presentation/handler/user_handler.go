package handler

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/internal/presentation/oas"
	user_usecase "github.com/neko-dream/server/internal/usecase/user"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type userHandler struct {
	*db.DBManager
	user_usecase.RegisterUserUseCase
}

// EditUserProfile implements oas.UserHandler.
func (u *userHandler) EditUserProfile(ctx context.Context) (*oas.EditUserProfileOK, error) {
	panic("unimplemented")
}

// GetUserProfile implements oas.UserHandler.
func (u *userHandler) GetUserProfile(ctx context.Context) (*oas.GetUserProfileOK, error) {

	return &oas.GetUserProfileOK{
		DisplayID:   "string",
		DisplayName: "string",
	}, nil
}

// RegisterUser ユーザー登録
func (u *userHandler) RegisterUser(ctx context.Context, params oas.OptRegisterUserReq) (oas.RegisterUserRes, error) {
	claim := session.GetSession(ctx)
	if !params.IsSet() {
		return nil, messages.RequiredParameterError
	}
	value := params.Value
	if err := value.Validate(); err != nil {
		utils.HandleError(ctx, err, "value.Validate")
		return nil, messages.RequiredParameterError
	}
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.InternalServerError
	}
	genderByte, err := value.Gender.MarshalText()
	if err != nil {
		utils.HandleError(ctx, err, "value.Gender")
		return nil, messages.InternalServerError
	}
	gender := string(genderByte)
	occupationByte, err := value.Occupation.Value.MarshalText()
	if err != nil {
		utils.HandleError(ctx, err, "value.Occupation")
		return nil, messages.InternalServerError
	}
	occupation := string(occupationByte)

	input := user_usecase.RegisterUserInput{
		UserID:        userID,
		DisplayID:     value.DisplayID,
		DisplayName:   value.DisplayName,
		PictureURL:    claim.Picture,
		YearOfBirth:   &value.YearOfBirth.Value,
		Gender:        lo.ToPtr(gender),
		Municipality:  lo.ToPtr(value.Municipality.Value),
		Occupation:    lo.ToPtr(occupation),
		HouseholdSize: &value.HouseholdSize.Value,
	}
	out, err := u.RegisterUserUseCase.Execute(ctx, input)
	if err != nil {
		utils.HandleError(ctx, err, "RegisterUserUseCase.Execute")
		return nil, messages.InternalServerError
	}

	return &oas.RegisterUserOK{
		DisplayID:   out.DisplayID,
		DisplayName: out.DisplayName,
	}, nil
}

func NewUserHandler(
	DBManager *db.DBManager,
) oas.UserHandler {
	return &userHandler{
		DBManager: DBManager,
	}
}

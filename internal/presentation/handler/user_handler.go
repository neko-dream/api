package handler

import (
	"context"
	"net/http"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/internal/presentation/oas"
	user_usecase "github.com/neko-dream/server/internal/usecase/user"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
)

type userHandler struct {
	*db.DBManager
	user_usecase.RegisterUserUseCase
}

// EditUserProfile implements oas.UserHandler.
func (u *userHandler) EditUserProfile(ctx context.Context) (oas.EditUserProfileRes, error) {
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
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.InternalServerError
	}
	value := params.Value
	if err := value.Validate(); err != nil {
		utils.HandleError(ctx, err, "value.Validate")
		return nil, messages.RequiredParameterError
	}
	gender := utils.ToPtrIfNotNullFunc(value.Gender.Null, func() string {
		txt, err := value.Gender.Value.MarshalText()
		if err != nil {
			utils.HandleError(ctx, err, "value.Gender")
			return ""
		}
		return string(txt)
	})
	occupation := utils.ToPtrIfNotNullFunc(value.Occupation.Null, func() string {
		txt, err := value.Occupation.Value.MarshalText()
		if err != nil {
			utils.HandleError(ctx, err, "value.Occupation")
			return ""
		}
		return string(txt)
	})
	municipality := utils.ToPtrIfNotNullValue(value.Municipality.Null, value.Municipality.Value)
	yearOfBirth := utils.ToPtrIfNotNullValue(value.YearOfBirth.Null, value.YearOfBirth.Value)

	input := user_usecase.RegisterUserInput{
		UserID:        userID,
		DisplayID:     value.DisplayID,
		DisplayName:   value.DisplayName,
		PictureURL:    claim.Picture,
		YearOfBirth:   yearOfBirth,
		Gender:        gender,
		Municipality:  municipality,
		Occupation:    occupation,
		HouseholdSize: &value.HouseholdSize.Value,
	}
	out, err := u.RegisterUserUseCase.Execute(ctx, input)
	if err != nil {
		utils.HandleError(ctx, err, "RegisterUserUseCase.Execute")
		return nil, err
	}

	w := http_utils.GetHTTPResponse(ctx)
	http.SetCookie(w, out.Cookie)

	return &oas.RegisterUserOK{
		DisplayID:   out.DisplayID,
		DisplayName: out.DisplayName,
		PictureURL:  utils.StringToOptString(claim.Picture),
	}, nil
}

func NewUserHandler(
	DBManager *db.DBManager,
	registerUserUsecase user_usecase.RegisterUserUseCase,
) oas.UserHandler {
	return &userHandler{
		DBManager:           DBManager,
		RegisterUserUseCase: registerUserUsecase,
	}
}

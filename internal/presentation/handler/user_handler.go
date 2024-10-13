package handler

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
	user_usecase "github.com/neko-dream/server/internal/usecase/user"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type userHandler struct {
	user_usecase.RegisterUserUseCase
	user_usecase.EditUserUseCase
	user_usecase.GetUserInformationQueryHandler
}

func NewUserHandler(
	registerUserUsecase user_usecase.RegisterUserUseCase,
	editUserUsecase user_usecase.EditUserUseCase,
	getUserInformationQueryHandler user_usecase.GetUserInformationQueryHandler,
) oas.UserHandler {
	return &userHandler{
		RegisterUserUseCase:            registerUserUsecase,
		EditUserUseCase:                editUserUsecase,
		GetUserInformationQueryHandler: getUserInformationQueryHandler,
	}
}

// GetUserInfo implements ユーザーの情報取得
func (u *userHandler) GetUserInfo(ctx context.Context) (oas.GetUserInfoRes, error) {
	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	res, err := u.GetUserInformationQueryHandler.Execute(ctx, user_usecase.GetUserInformationQuery{
		UserID: userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetUserInformationQueryHandler.Execute")
		return nil, messages.InternalServerError
	}

	userResp := oas.GetUserInfoOKUser{
		DisplayID:   *res.User.DisplayID(),
		DisplayName: *res.User.DisplayName(),
		IconURL:     utils.ToOptNil[oas.OptNilString](res.User.ProfileIconURL()),
	}
	var demographicsResp oas.GetUserInfoOKDemographics
	if res.User.Demographics() != nil {
		demographics := res.User.Demographics()
		var municipality oas.OptNilString
		if demographics.Municipality() != nil {
			municipality = oas.OptNilString{
				Set:   true,
				Value: demographics.Municipality().String(),
			}
		}
		var yearOfBirth oas.OptNilInt
		if demographics.YearOfBirth() != nil {
			yearOfBirth = oas.OptNilInt{
				Set:   true,
				Value: int(*demographics.YearOfBirth()),
			}
		}
		var householdSize oas.OptNilInt
		if demographics.HouseholdSize() != nil {
			householdSize = oas.OptNilInt{
				Set:   true,
				Value: int(*demographics.HouseholdSize()),
			}
		}
		demographicsResp = oas.GetUserInfoOKDemographics{
			YearOfBirth:   yearOfBirth,
			Municipality:  municipality,
			Occupation:    demographics.Occupation().String(),
			Gender:        demographics.Gender().String(),
			HouseholdSize: householdSize,
		}
	}

	return &oas.GetUserInfoOK{
		User:         userResp,
		Demographics: demographicsResp,
	}, nil
}

// EditUserProfile ユーザープロフィールの編集
func (u *userHandler) EditUserProfile(ctx context.Context, params oas.OptEditUserProfileReq) (oas.EditUserProfileRes, error) {
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

	var file *multipart.FileHeader
	if value.Icon.IsSet() {
		content, err := io.ReadAll(value.Icon.Value.File)
		if err != nil {
			utils.HandleError(ctx, err, "io.ReadAll")
			return nil, messages.InternalServerError
		}
		file, err = http_utils.MakeFileHeader(value.Icon.Value.Name, content)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}
	var yearOfBirth *int
	if !value.YearOfBirth.Null {
		yearOfBirth = &value.YearOfBirth.Value
	}

	var municipality *string
	if !value.Municipality.Null {
		municipality = &value.Municipality.Value
	}

	var householdSize *int
	if !value.HouseholdSize.Null {
		householdSize = &value.HouseholdSize.Value
	}

	out, err := u.EditUserUseCase.Execute(ctx, user_usecase.EditUserInput{
		UserID:       userID,
		DisplayName:  utils.ToPtrIfNotNullValue(value.DisplayName.Null, value.DisplayName.Value),
		Icon:         file,
		YearOfBirth:  yearOfBirth,
		Municipality: municipality,
		Occupation: utils.ToPtrIfNotNullFunc(value.Occupation.Null, func() *string {
			txt, err := value.Occupation.Value.MarshalText()
			if err != nil {
				return nil
			}
			return lo.ToPtr(string(txt))
		}),
		Gender: utils.ToPtrIfNotNullFunc(value.Gender.Null, func() *string {
			txt, err := value.Gender.Value.MarshalText()
			if err != nil {
				return nil
			}
			return lo.ToPtr(string(txt))
		}),
		HouseholdSize: householdSize,
	})
	if err != nil {
		utils.HandleError(ctx, err, "EditUserUseCase.Execute")
		return nil, err
	}

	w := http_utils.GetHTTPResponse(ctx)
	http.SetCookie(w, out.Cookie)

	return &oas.EditUserProfileOK{
		DisplayID:   out.DisplayID,
		DisplayName: out.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](claim.IconURL),
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

	var file *multipart.FileHeader
	if value.Icon.IsSet() {
		content, err := io.ReadAll(value.Icon.Value.File)
		if err != nil {
			utils.HandleError(ctx, err, "io.ReadAll")
			return nil, messages.InternalServerError
		}
		file, err = http_utils.MakeFileHeader(value.Icon.Value.Name, content)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}

	input := user_usecase.RegisterUserInput{
		UserID:       userID,
		DisplayID:    value.DisplayID,
		DisplayName:  value.DisplayName,
		Icon:         file,
		YearOfBirth:  utils.ToPtrIfNotNullValue(value.YearOfBirth.Null, value.YearOfBirth.Value),
		Municipality: utils.ToPtrIfNotNullValue(value.Municipality.Null, value.Municipality.Value),
		Occupation: utils.ToPtrIfNotNullFunc(value.Occupation.Null, func() *string {
			txt, err := value.Occupation.Value.MarshalText()
			if err != nil {
				return nil
			}
			return lo.ToPtr(string(txt))
		}),
		Gender: utils.ToPtrIfNotNullFunc(value.Gender.Null, func() *string {
			txt, err := value.Gender.Value.MarshalText()
			if err != nil {
				return nil
			}
			return lo.ToPtr(string(txt))
		}),
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
		IconURL:     utils.ToOptNil[oas.OptNilString](out.IconURL),
	}, nil
}

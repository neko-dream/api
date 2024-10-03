package handler

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
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

	var file *multipart.FileHeader
	if value.Icon.IsSet() {
		content, err := io.ReadAll(value.Icon.Value.File)
		if err != nil {
			utils.HandleError(ctx, err, "io.ReadAll")
			return nil, messages.InternalServerError
		}
		file, err = MakeFileHeader(value.Icon.Value.Name, content)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}

	input := user_usecase.RegisterUserInput{
		UserID:        userID,
		DisplayID:     value.DisplayID,
		DisplayName:   value.DisplayName,
		Icon:          file,
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
		IconURL:     utils.StringToOptString(claim.IconURL),
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
func MakeFileHeader(name string, dateBytes []byte) (*multipart.FileHeader, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, bytes.NewReader(dateBytes)); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// ダミーファイルをパースして、multipart.FileHeaderを取得する
	reader := multipart.NewReader(body, writer.Boundary())
	// 2M
	form, err := reader.ReadForm(2 * 1_000_000)
	if err != nil {
		return nil, err
	}
	fh := form.File["file"]

	return fh[0], nil
}

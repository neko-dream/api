package handler

import (
	"context"
	"mime/multipart"
	"net/http"

	opinion_query "github.com/neko-dream/api/internal/application/query/opinion"
	talksession_query "github.com/neko-dream/api/internal/application/query/talksession"
	user_query "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/application/usecase/user_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/infrastructure/http/cookie"
	"github.com/neko-dream/api/internal/presentation/oas"
	http_utils "github.com/neko-dream/api/pkg/http"
	"github.com/neko-dream/api/pkg/sort"
	"github.com/neko-dream/api/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type userHandler struct {
	getMyOpinionsQuery           opinion_query.GetMyOpinionsQuery
	browseJoinedTalkSessionQuery talksession_query.BrowseJoinedTalkSessionsQuery

	editUser     user_usecase.Edit
	registerUser user_usecase.Register
	withdrawUser user_usecase.Withdraw

	userDetail           user_query.Detail
	getByDisplayID       user_query.GetByDisplayID
	authorizationService service.AuthorizationService
	cookie.CookieManager
}

func NewUserHandler(
	getMyOpinionsQuery opinion_query.GetMyOpinionsQuery,
	browseJoinedTalkSessionQuery talksession_query.BrowseJoinedTalkSessionsQuery,

	editUser user_usecase.Edit,
	registerUser user_usecase.Register,
	withdrawUser user_usecase.Withdraw,

	userDetail user_query.Detail,
	getByDisplayID user_query.GetByDisplayID,
	authorizationService service.AuthorizationService,
	cookieManager cookie.CookieManager,
) oas.UserHandler {
	return &userHandler{
		getMyOpinionsQuery:           getMyOpinionsQuery,
		browseJoinedTalkSessionQuery: browseJoinedTalkSessionQuery,
		editUser:                     editUser,
		registerUser:                 registerUser,
		withdrawUser:                 withdrawUser,
		userDetail:                   userDetail,
		getByDisplayID:               getByDisplayID,
		authorizationService:         authorizationService,
		CookieManager:                cookieManager,
	}
}

// OpinionsHistory implements oas.UserHandler.
func (u *userHandler) OpinionsHistory(ctx context.Context, params oas.OpinionsHistoryParams) (oas.OpinionsHistoryRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.OpinionsHistory")
	defer span.End()

	authCtx, err := u.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID

	var sortKey sort.SortKey
	if params.Sort.IsSet() {
		txt, err := params.Sort.Value.MarshalText()
		if err != nil {
			utils.HandleError(ctx, err, "params.Sort.Value.MarshalText")
			return nil, messages.InternalServerError
		}
		sortKey = sort.SortKey(txt)
	}
	var limit, offset *int
	if params.Limit.IsSet() {
		limit = &params.Limit.Value
	}
	if params.Offset.IsSet() {
		offset = &params.Offset.Value
	}

	out, err := u.getMyOpinionsQuery.Execute(ctx, opinion_query.GetMyOpinionsQueryInput{
		UserID:  userID,
		SortKey: sortKey,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetUserOpinionListQueryHandler.Execute")
		return nil, messages.InternalServerError
	}

	opinions := make([]oas.OpinionWithReplyCount, 0, len(out.Opinions))
	for _, opinion := range out.Opinions {
		opinions = append(opinions, oas.OpinionWithReplyCount{
			Opinion:    opinion.Opinion.ToResponse(),
			User:       opinion.User.ToResponse(),
			ReplyCount: opinion.ReplyCount,
		})
	}

	return &oas.OpinionsHistoryOK{
		Opinions: opinions,
		Pagination: oas.OpinionsHistoryOKPagination{
			TotalCount: out.TotalCount,
		},
	}, nil
}

// SessionsHistory implements oas.UserHandler.
func (u *userHandler) SessionsHistory(ctx context.Context, params oas.SessionsHistoryParams) (oas.SessionsHistoryRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.SessionsHistory")
	defer span.End()

	authCtx, err := u.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID

	var status string
	if params.Status.IsSet() {
		txt, err := params.Status.Value.MarshalText()
		if err != nil {
			utils.HandleError(ctx, err, "params.Status.Value.MarshalText")
			return nil, messages.InternalServerError
		}
		status = string(txt)
	}
	var limit, offset *int
	if params.Limit.IsSet() {
		limit = &params.Limit.Value
	}
	if params.Offset.IsSet() {
		offset = &params.Offset.Value
	}

	out, err := u.browseJoinedTalkSessionQuery.Execute(ctx, talksession_query.BrowseJoinedTalkSessionsQueryInput{
		UserID: userID,
		Status: talksession_query.Status(status),
		Theme:  utils.ToPtrIfNotNullValue(!params.Theme.IsSet(), params.Theme.Value),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		utils.HandleError(ctx, err, "SearchTalkSessionsQuery.Execute")
		return nil, messages.InternalServerError
	}

	talkSessions := make([]oas.SessionsHistoryOKTalkSessionsItem, 0, len(out.TalkSessions))
	for _, talkSession := range out.TalkSessions {
		talkSessions = append(talkSessions, oas.SessionsHistoryOKTalkSessionsItem{
			TalkSession:  talkSession.ToResponse(),
			OpinionCount: talkSession.OpinionCount,
		})
	}

	return &oas.SessionsHistoryOK{
		TalkSessions: talkSessions,
		Pagination: oas.OffsetPagination{
			TotalCount: out.TotalCount,
		},
	}, nil
}

// GetUserInfo implements ユーザーの情報取得
func (u *userHandler) GetUserInfo(ctx context.Context) (oas.GetUserInfoRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.GetUserInfo")
	defer span.End()

	authCtx, err := u.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID

	res, err := u.userDetail.Execute(ctx, user_query.DetailInput{
		UserID: userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetUserInformationQueryHandler.Execute")
		return nil, messages.InternalServerError
	}

	userResp := res.User.User.ToResponse()

	var demographicsResp oas.UserDemographics
	if res.User.UserDemographic != nil {
		demographicsResp = res.User.UserDemographic.ToResponse()
	}

	email := res.User.UserAuth.ToEmailResponse()

	return &oas.GetUserInfoOK{
		User:         userResp,
		Demographics: demographicsResp,
		Email:        email,
	}, nil
}

// UpdateUserProfile ユーザープロフィールの編集
func (u *userHandler) UpdateUserProfile(ctx context.Context, params *oas.UpdateUserProfileReq) (oas.UpdateUserProfileRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.EditUserProfile")
	defer span.End()

	if params == nil {
		return nil, messages.RequiredParameterError
	}

	authCtx, err := u.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID
	deleteIcon := false
	if params.DeleteIcon.IsSet() {
		deleteIcon = params.DeleteIcon.Value
	}
	var email *string
	if params.Email.IsSet() {
		email = &params.Email.Value
	}

	var file *multipart.FileHeader
	if params.Icon.IsSet() {
		file, err = http_utils.CreateFileHeader(ctx, params.Icon.Value.File, params.Icon.Value.Name)
		if err != nil {
			utils.HandleError(ctx, err, "CreateFileHeader")
			return nil, messages.InternalServerError
		}
	}

	var dateOfBirth *int
	if birthStr, ok := params.DateOfBirth.Get(); ok {
		dateOfBirth = lo.ToPtr(int(birthStr))
	}

	var city *string
	if params.City.IsSet() && params.City.Value != "" {
		city = &params.City.Value
	}

	var prefecture *string
	if params.Prefecture.IsSet() && params.Prefecture.Value != "" {
		prefecture = &params.Prefecture.Value
	}
	var displayName *string
	if params.DisplayName.IsSet() && params.DisplayName.Value != "" {
		if params.DisplayName.Value == "" {
			return nil, messages.UserDisplayNameTooShort
		}
		displayName = &params.DisplayName.Value
	}

	var gender *string
	if params.Gender.IsSet() && params.Gender.Value != "" {
		gender = lo.ToPtr(string(params.Gender.Value))
	}

	out, err := u.editUser.Execute(ctx, user_usecase.EditInput{
		UserID:      userID,
		DisplayName: displayName,
		Icon:        file,
		Email:       email,
		DateOfBirth: dateOfBirth,
		City:        city,
		Gender:      gender,
		Prefecture:  prefecture,
		DeleteIcon:  deleteIcon,
	})
	if err != nil {
		utils.HandleError(ctx, err, "EditUserUseCase.Execute")
		return nil, err
	}

	w := http_utils.GetHTTPResponse(ctx)
	http.SetCookie(w, u.CookieManager.CreateSessionCookie(out.Token))

	resp := out.ToResponse()
	return &resp, nil
}

// EstablishUser ユーザー登録
func (u *userHandler) EstablishUser(ctx context.Context, params *oas.EstablishUserReq) (oas.EstablishUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.RegisterUser")
	defer span.End()

	if params == nil {
		return nil, messages.RequiredParameterError
	}

	authCtx, err := u.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	value := params
	if err := value.Validate(); err != nil {
		utils.HandleError(ctx, err, "value.Validate")
		return nil, messages.RequiredParameterError
	}

	var file *multipart.FileHeader
	if value.Icon.IsSet() {
		file, err = http_utils.CreateFileHeader(ctx, value.Icon.Value.File, value.Icon.Value.Name)
		if err != nil {
			utils.HandleError(ctx, err, "CreateFileHeader")
			return nil, messages.InternalServerError
		}
	}

	var prefecture *string
	if value.Prefecture.IsSet() {
		prefecture = &value.Prefecture.Value
	}
	var dateOfBirth *int
	if birthStr, ok := value.DateOfBirth.Get(); ok {
		dateOfBirth = lo.ToPtr(int(birthStr))
	}
	var gender *string
	if value.Gender.IsSet() && value.Gender.Value != "" {
		gender = lo.ToPtr(string(value.Gender.Value))
	}

	if value.DisplayID == "" {
		return nil, messages.UserDisplayIDTooShort
	}
	if value.DisplayName == "" {
		return nil, messages.UserDisplayNameTooShort
	}

	input := user_usecase.RegisterInput{
		SessionID:   authCtx.SessionID,
		UserID:      authCtx.UserID,
		DisplayID:   value.DisplayID,
		DisplayName: value.DisplayName,
		Icon:        file,
		DateOfBirth: dateOfBirth,
		City:        utils.ToPtrIfNotNullValue(!value.City.IsSet(), value.City.Value),
		Gender:      gender,
		Prefecture:  prefecture,
	}
	out, err := u.registerUser.Execute(ctx, input)
	if err != nil {
		utils.HandleError(ctx, err, "RegisterUserUseCase.Execute")
		return nil, err
	}

	w := http_utils.GetHTTPResponse(ctx)
	http.SetCookie(w, u.CookieManager.CreateSessionCookie(out.Token))

	resp := out.ToResponse()
	return &resp, nil
}

func (u *userHandler) GetUserByDisplayID(ctx context.Context, params oas.GetUserByDisplayIDParams) (oas.GetUserByDisplayIDRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.GetUserByDisplayID")
	defer span.End()

	result, err := u.getByDisplayID.Execute(ctx, user_query.GetByDisplayIDInput{
		DisplayID: params.DisplayID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetByDisplayID.Execute")
		return &oas.GetUserByDisplayIDInternalServerError{}, nil
	}

	if result.User == nil {
		return nil, messages.UserNotFound
	}

	userResp := result.User.ToResponse()
	return &userResp, nil
}

// WithdrawUser implements oas.UserHandler.
func (u *userHandler) WithdrawUser(ctx context.Context) (oas.WithdrawUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.WithdrawUser")
	defer span.End()

	authCtx, err := u.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// 退会処理を実行
	output, err := u.withdrawUser.Execute(ctx, user_usecase.WithdrawInput{
		UserID: authCtx.UserID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "withdrawUser.Execute")
		return nil, err
	}

	return &oas.WithdrawUserOK{
		Message:        output.Message,
		WithdrawalDate: output.WithdrawalDate,
	}, nil
}

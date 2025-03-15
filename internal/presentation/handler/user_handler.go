package handler

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/infrastructure/http/cookie"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/user_command"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
	talksession_query "github.com/neko-dream/server/internal/usecase/query/talksession"
	user_query "github.com/neko-dream/server/internal/usecase/query/user"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/sort"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type userHandler struct {
	getMyOpinionsQuery           opinion_query.GetMyOpinionsQuery
	browseJoinedTalkSessionQuery talksession_query.BrowseJoinedTalkSessionsQuery

	editUser     user_command.Edit
	registerUser user_command.Register

	userDetail user_query.Detail
	cookie.CookieManager
}

func NewUserHandler(
	getMyOpinionsQuery opinion_query.GetMyOpinionsQuery,
	browseJoinedTalkSessionQuery talksession_query.BrowseJoinedTalkSessionsQuery,

	editUser user_command.Edit,
	registerUser user_command.Register,

	userDetail user_query.Detail,
	cookieManager cookie.CookieManager,
) oas.UserHandler {
	return &userHandler{
		getMyOpinionsQuery:           getMyOpinionsQuery,
		browseJoinedTalkSessionQuery: browseJoinedTalkSessionQuery,
		editUser:                     editUser,
		registerUser:                 registerUser,
		userDetail:                   userDetail,
		CookieManager:                cookieManager,
	}
}

// OpinionsHistory implements oas.UserHandler.
func (u *userHandler) OpinionsHistory(ctx context.Context, params oas.OpinionsHistoryParams) (oas.OpinionsHistoryRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.OpinionsHistory")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

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

	opinions := make([]oas.OpinionsHistoryOKOpinionsItem, 0, len(out.Opinions))
	for _, opinion := range out.Opinions {
		vt := oas.OptOpinionsHistoryOKOpinionsItemOpinionVoteType{}
		if opinion.GetParentVoteType() != nil {
			vt = oas.OptOpinionsHistoryOKOpinionsItemOpinionVoteType{
				Set:   true,
				Value: oas.OpinionsHistoryOKOpinionsItemOpinionVoteType(*opinion.GetParentVoteType()),
			}
		}

		opinions = append(opinions, oas.OpinionsHistoryOKOpinionsItem{
			Opinion: oas.OpinionsHistoryOKOpinionsItemOpinion{
				ID:           opinion.Opinion.OpinionID.String(),
				Title:        utils.ToOpt[oas.OptString](opinion.Opinion.Title),
				Content:      opinion.Opinion.Content,
				VoteType:     vt,
				ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
				PictureURL:   utils.ToOpt[oas.OptString](opinion.Opinion.PictureURL),
			},
			User: oas.OpinionsHistoryOKOpinionsItemUser{
				DisplayID:   opinion.User.DisplayID,
				DisplayName: opinion.User.DisplayName,
				IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.IconURL),
			},
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

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

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
			TalkSession: oas.SessionsHistoryOKTalkSessionsItemTalkSession{
				ID:    talkSession.TalkSessionID.String(),
				Theme: talkSession.Theme,
				Owner: oas.SessionsHistoryOKTalkSessionsItemTalkSessionOwner{
					DisplayID:   talkSession.User.DisplayID,
					DisplayName: talkSession.User.DisplayName,
					IconURL:     utils.ToOptNil[oas.OptNilString](talkSession.User.IconURL),
				},
				CreatedAt:        talkSession.CreatedAt.Format(time.RFC3339),
				ScheduledEndTime: talkSession.ScheduledEndTime.Format(time.RFC3339),
			},
			OpinionCount: talkSession.OpinionCount,
		})
	}

	return &oas.SessionsHistoryOK{
		TalkSessions: talkSessions,
		Pagination: oas.SessionsHistoryOKPagination{
			TotalCount: out.TotalCount,
		},
	}, nil
}

// GetUserInfo implements ユーザーの情報取得
func (u *userHandler) GetUserInfo(ctx context.Context) (oas.GetUserInfoRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.GetUserInfo")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	res, err := u.userDetail.Execute(ctx, user_query.DetailInput{
		UserID: userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetUserInformationQueryHandler.Execute")
		return nil, messages.InternalServerError
	}

	userResp := oas.GetUserInfoOKUser{
		DisplayID:   res.User.DisplayID,
		DisplayName: res.User.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](res.User.IconURL),
	}

	var demographicsResp oas.GetUserInfoOKDemographics
	if res.User.UserDemographic != nil {
		demographics := res.User.UserDemographic
		var city oas.OptNilString
		if demographics.City != nil {
			city = oas.OptNilString{
				Set:   true,
				Value: *demographics.City,
			}
		}
		var yearOfBirth oas.OptNilInt
		if demographics.YearOfBirth != nil {
			yearOfBirth = oas.OptNilInt{
				Set:   true,
				Value: *demographics.YearOfBirth,
			}
		}
		var prefecture oas.OptNilString
		if demographics.Prefecture != nil {
			prefecture = oas.OptNilString{
				Set:   true,
				Value: *demographics.Prefecture,
			}
		}
		var gender oas.OptNilString
		if demographics.GenderString() != nil {
			gender = oas.OptNilString{
				Set:   true,
				Value: *demographics.GenderString(),
			}
		}

		demographicsResp = oas.GetUserInfoOKDemographics{
			YearOfBirth: yearOfBirth,
			Gender:      gender,
			Prefecture:  prefecture,
			City:        city,
		}
	}

	var email oas.OptNilString
	if res.User.UserAuth.Email != nil {
		email = oas.OptNilString{
			Set:   true,
			Value: *res.User.UserAuth.Email,
		}
	}

	return &oas.GetUserInfoOK{
		User:         userResp,
		Demographics: demographicsResp,
		Email:        email,
	}, nil
}

// EditUserProfile ユーザープロフィールの編集
func (u *userHandler) EditUserProfile(ctx context.Context, params oas.OptEditUserProfileReq) (oas.EditUserProfileRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.EditUserProfile")
	defer span.End()

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
	deleteIcon := false
	if value.DeleteIcon.IsSet() {
		deleteIcon = value.DeleteIcon.Value
	}
	var email *string
	if value.Email.IsSet() {
		email = &value.Email.Value
	}

	var file *multipart.FileHeader
	if value.Icon.IsSet() {
		file, err = http_utils.CreateFileHeader(ctx, value.Icon.Value.File, value.Icon.Value.Name)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}
	var yearOfBirth *int
	if !value.YearOfBirth.Null && value.YearOfBirth.Value != 0 {
		yearOfBirth = &value.YearOfBirth.Value
	}
	var city *string
	if !value.City.Null && value.City.Value != "" {
		city = &value.City.Value
	}

	var prefecture *string
	if !value.Prefecture.Null && value.Prefecture.Value != "" {
		prefecture = &value.Prefecture.Value
	}
	var displayName *string
	if !value.DisplayName.Null && value.DisplayName.Value != "" {
		if value.DisplayName.Value == "" {
			return nil, messages.UserDisplayIDTooShort
		}
		if utf8.RuneCountInString(value.DisplayName.Value) > 20 || utf8.RuneCountInString(value.DisplayName.Value) < 4 {
			return nil, messages.UserDisplayIDTooShort
		}
		displayName = &value.DisplayName.Value
	}
	var gender *string
	if value.Gender.IsSet() && !value.Gender.IsNull() {
		txt, err := value.Gender.Value.MarshalText()
		if err != nil {
			return nil, messages.InternalServerError
		}
		if string(txt) != "" {
			gender = lo.ToPtr(string(txt))
		}
	}

	out, err := u.editUser.Execute(ctx, user_command.EditInput{
		UserID:      userID,
		DisplayName: displayName,
		Icon:        file,
		Email:       email,
		YearOfBirth: yearOfBirth,
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

	return &oas.EditUserProfileOK{
		DisplayID:   out.DisplayID,
		DisplayName: out.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](claim.IconURL),
	}, nil
}

// RegisterUser ユーザー登録
func (u *userHandler) RegisterUser(ctx context.Context, params oas.OptRegisterUserReq) (oas.RegisterUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.RegisterUser")
	defer span.End()

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
		file, err = http_utils.CreateFileHeader(ctx, value.Icon.Value.File, value.Icon.Value.Name)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}

	var prefecture *string
	if value.Prefecture.IsSet() {
		prefecture = &value.Prefecture.Value
	}

	input := user_command.RegisterInput{
		UserID:      userID,
		DisplayID:   value.DisplayID,
		DisplayName: value.DisplayName,
		Icon:        file,
		YearOfBirth: utils.ToPtrIfNotNullValue(value.YearOfBirth.Null, value.YearOfBirth.Value),
		City:        utils.ToPtrIfNotNullValue(value.City.Null, value.City.Value),
		Gender: utils.ToPtrIfNotNullFunc(value.Gender.Null, func() *string {
			txt, err := value.Gender.Value.MarshalText()
			if err != nil {
				return nil
			}
			return lo.ToPtr(string(txt))
		}),
		Prefecture: prefecture,
	}
	out, err := u.registerUser.Execute(ctx, input)
	if err != nil {
		utils.HandleError(ctx, err, "RegisterUserUseCase.Execute")
		return nil, err
	}

	w := http_utils.GetHTTPResponse(ctx)
	http.SetCookie(w, u.CookieManager.CreateSessionCookie(out.Token))

	return &oas.RegisterUserOK{
		DisplayID:   out.DisplayID,
		DisplayName: out.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](out.IconURL),
	}, nil
}

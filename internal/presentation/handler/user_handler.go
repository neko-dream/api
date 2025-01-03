package handler

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"unicode/utf8"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
	opinion_usecase "github.com/neko-dream/server/internal/usecase/opinion"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	user_usecase "github.com/neko-dream/server/internal/usecase/user"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type userHandler struct {
	user_usecase.RegisterUserUseCase
	user_usecase.EditUserUseCase
	user_usecase.GetUserInformationQueryHandler
	opinion_usecase.GetUserOpinionListQueryHandler
	talk_session_usecase.SearchTalkSessionsQuery
}

func NewUserHandler(
	registerUserUsecase user_usecase.RegisterUserUseCase,
	editUserUsecase user_usecase.EditUserUseCase,
	getUserInformationQueryHandler user_usecase.GetUserInformationQueryHandler,
	getUserOpinionListQueryHandler opinion_usecase.GetUserOpinionListQueryHandler,
	getTalkSessoinHistoriesQuery talk_session_usecase.SearchTalkSessionsQuery,
) oas.UserHandler {
	return &userHandler{
		RegisterUserUseCase:            registerUserUsecase,
		EditUserUseCase:                editUserUsecase,
		GetUserInformationQueryHandler: getUserInformationQueryHandler,
		GetUserOpinionListQueryHandler: getUserOpinionListQueryHandler,
		SearchTalkSessionsQuery:        getTalkSessoinHistoriesQuery,
	}
}

// OpinionsHistory implements oas.UserHandler.
func (u *userHandler) OpinionsHistory(ctx context.Context, params oas.OpinionsHistoryParams) (oas.OpinionsHistoryRes, error) {

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	var sortKey *string
	if params.Sort.IsSet() {
		txt, err := params.Sort.Value.MarshalText()
		if err != nil {
			utils.HandleError(ctx, err, "params.Sort.Value.MarshalText")
			return nil, messages.InternalServerError
		}
		sortKey = lo.ToPtr(string(txt))
	}
	var limit, offset *int
	if params.Limit.IsSet() {
		limit = &params.Limit.Value
	}
	if params.Offset.IsSet() {
		offset = &params.Offset.Value
	}

	out, err := u.GetUserOpinionListQueryHandler.Execute(ctx, opinion_usecase.GetUserOpinionListQuery{
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
		opinions = append(opinions, oas.OpinionsHistoryOKOpinionsItem{
			Opinion: oas.OpinionsHistoryOKOpinionsItemOpinion{
				ID:      opinion.Opinion.OpinionID,
				Title:   utils.ToOpt[oas.OptString](opinion.Opinion.Title),
				Content: opinion.Opinion.Content,
				VoteType: oas.OptOpinionsHistoryOKOpinionsItemOpinionVoteType{
					Set:   true,
					Value: oas.OpinionsHistoryOKOpinionsItemOpinionVoteType(opinion.Opinion.VoteType),
				},
				ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
				PictureURL:   utils.ToOpt[oas.OptString](opinion.Opinion.PictureURL),
			},
			User: oas.OpinionsHistoryOKOpinionsItemUser{
				DisplayID:   opinion.User.ID,
				DisplayName: opinion.User.Name,
				IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.Icon),
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

	out, err := u.SearchTalkSessionsQuery.Execute(ctx, talk_session_usecase.SearchTalkSessionsInput{
		UserID: userID,
		Status: status,
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
				ID:    talkSession.ID,
				Theme: talkSession.Theme,
				Owner: oas.SessionsHistoryOKTalkSessionsItemTalkSessionOwner{
					DisplayID:   talkSession.Owner.DisplayID,
					DisplayName: talkSession.Owner.DisplayName,
					IconURL:     utils.ToOptNil[oas.OptNilString](talkSession.Owner.IconURL),
				},
				CreatedAt:        talkSession.CreatedAt,
				ScheduledEndTime: talkSession.ScheduledEndTime,
			},
			OpinionCount: talkSession.OpinionCount,
		})
	}

	return &oas.SessionsHistoryOK{
		TalkSessions: talkSessions,
		Pagination: oas.SessionsHistoryOKPagination{
			TotalCount: out.TotalCount,
			Limit:      out.Limit,
			Offset:     out.Offset,
		},
	}, nil
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
		var city oas.OptNilString
		if demographics.City() != nil {
			city = oas.OptNilString{
				Set:   true,
				Value: demographics.City().String(),
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
		var prefecture oas.OptNilString
		if demographics.Prefecture() != nil {
			prefecture = oas.OptNilString{
				Set:   true,
				Value: *demographics.Prefecture(),
			}
		}

		demographicsResp = oas.GetUserInfoOKDemographics{
			YearOfBirth:   yearOfBirth,
			City:          city,
			Occupation:    demographics.Occupation().String(),
			Gender:        demographics.Gender().String(),
			HouseholdSize: householdSize,
			Prefecture:    prefecture,
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
	deleteIcon := false
	if value.DeleteIcon.IsSet() {
		deleteIcon = value.DeleteIcon.Value
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
	var city *string
	if !value.City.Null {
		city = &value.City.Value
	}
	var householdSize *int
	if !value.HouseholdSize.Null {
		householdSize = &value.HouseholdSize.Value
	}
	var prefecture *string
	if !value.Prefecture.Null {
		prefecture = &value.Prefecture.Value
	}
	var displayName *string
	if !value.DisplayName.Null {
		if value.DisplayName.Value == "" {
			return nil, messages.UserDisplayIDTooShort
		}

		if utf8.RuneCountInString(value.DisplayName.Value) > 20 || utf8.RuneCountInString(value.DisplayName.Value) < 4 {
			return nil, messages.UserDisplayIDTooShort
		}

		displayName = &value.DisplayName.Value
	}
	var occupation *string
	if value.Occupation.IsSet() && !value.Occupation.Null {
		txt, err := value.Occupation.Value.MarshalText()
		if err != nil {
			utils.HandleError(ctx, err, "value.Occupation.Value.MarshalText")
			return nil, messages.InternalServerError
		}
		occupation = lo.ToPtr(string(txt))
	}

	out, err := u.EditUserUseCase.Execute(ctx, user_usecase.EditUserInput{
		UserID:      userID,
		DisplayName: displayName,
		Icon:        file,
		YearOfBirth: yearOfBirth,
		City:        city,
		Occupation:  occupation,
		Gender: utils.ToPtrIfNotNullFunc(value.Gender.Null, func() *string {
			txt, err := value.Gender.Value.MarshalText()
			if err != nil {
				return nil
			}
			return lo.ToPtr(string(txt))
		}),
		HouseholdSize: householdSize,
		Prefecture:    prefecture,
		DeleteIcon:    deleteIcon,
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

	var prefecture *string
	if value.Prefecture.IsSet() {
		prefecture = &value.Prefecture.Value
	}

	input := user_usecase.RegisterUserInput{
		UserID:      userID,
		DisplayID:   value.DisplayID,
		DisplayName: value.DisplayName,
		Icon:        file,
		YearOfBirth: utils.ToPtrIfNotNullValue(value.YearOfBirth.Null, value.YearOfBirth.Value),
		City:        utils.ToPtrIfNotNullValue(value.City.Null, value.City.Value),
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
		Prefecture:    prefecture,
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

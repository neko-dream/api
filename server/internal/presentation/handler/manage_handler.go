package handler

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	opinionquery "github.com/neko-dream/api/internal/application/query/opinion"
	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/analysis"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/internal/presentation/oas"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type manageHandler struct {
	userManagement        userquery.Management
	setCreationPermission admin_usecase.SetSessionCreationPermission
	userInvitations       admin_usecase.ManageUserInvitations
	editTalkSession       admin_usecase.EditTalkSession
	manageOpinions        opinionquery.ManageOpinions
	analysis.AnalysisService
	analysis.AnalysisRepository
	analysisScheduler analysis_usecase.AnalysisScheduler
	*db.DBManager
	authorizationService service.AuthorizationService
	session.TokenManager
}

// GetUserListManage implements oas.ManageHandler.
func (m *manageHandler) GetUserListManage(ctx context.Context, params oas.GetUserListManageParams) (*oas.ManagedUserList, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetUserListManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}
	limit, offset := params.Limit.Or(20), params.Offset.Or(0)
	search := strings.TrimSpace(params.Search.Or(""))
	if limit < 1 || limit > 100 || offset < 0 || utf8.RuneCountInString(search) > 200 {
		return nil, messages.BadRequestError
	}
	out, err := m.userManagement.List(ctx, userquery.ManagementListInput{Offset: offset, Limit: limit, Search: search})
	if err != nil {
		return nil, err
	}
	items := make([]oas.ManagedUser, 0, len(out.Users))
	for _, u := range out.Users {
		items = append(items, managedUserResponse(u))
	}
	return &oas.ManagedUserList{Users: items, TotalCount: out.TotalCount}, nil
}

func (m *manageHandler) GetUserDetailManage(ctx context.Context, params oas.GetUserDetailManageParams) (oas.GetUserDetailManageRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetUserDetailManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}
	id, err := shared.ParseUUID[user.User](params.UserID)
	if err != nil {
		return nil, messages.BadRequestError
	}
	u, err := m.userManagement.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return &oas.GetUserDetailManageNotFound{}, nil
	}
	out := managedUserResponse(*u)
	return &out, nil
}

func managedUserResponse(u userquery.ManagedUser) oas.ManagedUser {
	out := oas.ManagedUser{UserID: u.UserID.String(), DisplayID: u.DisplayID, DisplayName: u.DisplayName, EmailVerified: u.EmailVerified, Providers: u.Providers, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
	if u.Email != nil {
		out.Email = oas.NewOptString(*u.Email)
	}
	out.Permissions = oas.UserPermissions{CanCreateTalkSession: u.CanCreateTalkSession}
	out.PasswordSetupRequired = u.PasswordSetupRequired
	if u.WithdrawalDate != nil {
		out.WithdrawalDate = oas.NewOptDateTime(*u.WithdrawalDate)
	}
	return out
}

// GetUserStatsTotalManage implements oas.ManageHandler.
func (m *manageHandler) GetUserStatsTotalManage(ctx context.Context) (*oas.UserStatsResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetUserStatsTotalManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}

	res, err := m.GetQueries(ctx).GetUserStats(ctx)
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.GetUserStats")
		return nil, err
	}

	totalTalkSessionCount, err := m.GetQueries(ctx).GetAllTalkSessionCount(ctx)
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.GetAllTalkSessionCount")
		return nil, err
	}

	return &oas.UserStatsResponse{
		UserCount:             int32(res.TotalUsers),
		UniqueActionUserCount: int32(res.ActiveUsers),
		TalkSessionCount:      int32(totalTalkSessionCount),
	}, nil
}

// GetUserStatsListManage implements oas.ManageHandler.
func (m *manageHandler) GetUserStatsListManage(ctx context.Context, params oas.GetUserStatsListManageParams) ([]oas.UserStatsResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetUserStatsListManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}

	var page, limit int32
	if params.Offset.IsSet() {
		page = params.Offset.Value
	} else {
		page = 1
	}
	if params.Limit.IsSet() {
		limit = params.Limit.Value
	} else {
		limit = 10
	}
	var stats []oas.UserStatsResponse
	if params.Range == "daily" {
		rows, err := m.GetQueries(ctx).GetDailyUserStats(ctx, model.GetDailyUserStatsParams{
			Offset: page,
			Limit:  limit,
		})
		if err != nil {
			utils.HandleError(ctx, err, "GetQueries.GetDailyUserStats")
			return []oas.UserStatsResponse{}, err
		}

		for _, row := range rows {
			stats = append(stats, oas.UserStatsResponse{
				Date:                  row.ActivityDate,
				UserCount:             int32(row.TotalUsers),
				UniqueActionUserCount: int32(row.ActiveUsers),
			})
		}
	} else if params.Range == "weekly" {
		rows, err := m.GetQueries(ctx).GetWeeklyUserStats(ctx, model.GetWeeklyUserStatsParams{
			Offset: page,
			Limit:  limit,
		})
		if err != nil {
			utils.HandleError(ctx, err, "GetQueries.GetWeeklyUserStats")
			return []oas.UserStatsResponse{}, err
		}
		for _, row := range rows {
			stats = append(stats, oas.UserStatsResponse{
				Date:                  row.ActivityDate,
				UserCount:             int32(row.TotalUsers),
				UniqueActionUserCount: int32(row.ActiveUsers),
			})
		}
	}

	return stats, nil
}

func NewManageHandler(
	dbm *db.DBManager,
	ansv analysis.AnalysisService,
	arep analysis.AnalysisRepository,
	scheduler analysis_usecase.AnalysisScheduler,
	authorizationService service.AuthorizationService,
	tokenManager session.TokenManager,
	userManagement userquery.Management,
	setCreationPermission admin_usecase.SetSessionCreationPermission,
	userInvitations admin_usecase.ManageUserInvitations,
	editTalkSession admin_usecase.EditTalkSession,
	manageOpinions opinionquery.ManageOpinions,
) oas.ManageHandler {
	return &manageHandler{
		DBManager:             dbm,
		AnalysisService:       ansv,
		AnalysisRepository:    arep,
		analysisScheduler:     scheduler,
		authorizationService:  authorizationService,
		TokenManager:          tokenManager,
		userManagement:        userManagement,
		setCreationPermission: setCreationPermission,
		userInvitations:       userInvitations,
		editTalkSession:       editTalkSession,
		manageOpinions:        manageOpinions,
	}
}

// GetTalkSessionListManage implements oas.ManageHandler.
func (m *manageHandler) GetTalkSessionListManage(ctx context.Context, params oas.GetTalkSessionListManageParams) (*oas.TalkSessionListResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetTalkSessionListManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}

	limit, ok := params.Limit.Get()
	if !ok {
		limit = 10
	}

	offset, ok := params.Offset.Get()
	if !ok {
		offset = 0
	}

	rows, err := m.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:   limit,
		Offset:  offset,
		SortKey: sql.NullString{String: "latest", Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.ListTalkSessions")
		return &oas.TalkSessionListResponse{
			TotalCount: 0,
		}, err
	}

	totalCount, err := m.GetQueries(ctx).GetAllTalkSessionCount(ctx)
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.GetAllTalkSessionCount")
		return &oas.TalkSessionListResponse{
			TotalCount: 0,
		}, err
	}

	talkSessionStats := make([]oas.TalkSessionStats, 0, len(rows))
	for _, row := range rows {
		var city oas.OptString
		if row.TalkSession.City.Valid {
			city = oas.OptString{
				Value: row.TalkSession.City.String,
				Set:   true,
			}
		}
		var description string
		if row.TalkSession.Description.Valid {
			description = row.TalkSession.Description.String
		}

		var prefecture oas.OptString
		if row.TalkSession.Prefecture.Valid {
			prefecture = oas.OptString{
				Value: row.TalkSession.Prefecture.String,
				Set:   true,
			}
		}

		var iconURL string
		if row.User.IconUrl.Valid {
			iconURL = row.User.IconUrl.String
		} else {
			iconURL = ""
		}
		owner := oas.UserForManage{
			DisplayID:   row.User.DisplayID.String,
			DisplayName: row.User.DisplayName.String,
			IconURL:     iconURL,
		}
		var thumbnailURL string
		if row.TalkSession.ThumbnailUrl.Valid {
			thumbnailURL = row.TalkSession.ThumbnailUrl.String
		} else {
			thumbnailURL = ""
		}

		talkSessionStats = append(talkSessionStats, oas.TalkSessionStats{
			TalkSessionID:    row.TalkSession.TalkSessionID.String(),
			Theme:            row.TalkSession.Theme,
			Description:      description,
			City:             city,
			Prefecture:       prefecture,
			ThumbnailURL:     thumbnailURL,
			Hidden:           row.TalkSession.HideReport.Bool,
			HideTop:          row.TalkSession.HideTop,
			Owner:            owner,
			ScheduledEndTime: row.TalkSession.ScheduledEndTime,
			CreatedAt:        row.TalkSession.CreatedAt.Format(time.RFC3339),
			OpinionCount:     int32(int(row.OpinionCount)),
			VoteCount:        int32(int(row.VoteCount)),
			VoteUserCount:    int32(int(row.VoteUserCount)),
		})
	}

	return &oas.TalkSessionListResponse{
		TotalCount:       int32(totalCount),
		TalkSessionStats: talkSessionStats,
	}, nil
}

// GetTalkSessionManage implements oas.ManageHandler.
func (m *manageHandler) GetTalkSessionManage(ctx context.Context, params oas.GetTalkSessionManageParams) (*oas.TalkSessionForManage, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetTalkSessionManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}
	return m.talkSessionForManage(ctx, talkSessionID)
}

// EditTalkSessionManage implements oas.ManageHandler.
func (m *manageHandler) EditTalkSessionManage(ctx context.Context, req *oas.TalkSessionEditForManage, params oas.EditTalkSessionManageParams) (*oas.TalkSessionForManage, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.EditTalkSessionManage")
	defer span.End()

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}
	if err := m.editTalkSession.Execute(ctx, admin_usecase.EditTalkSessionInput{
		TalkSessionID:    talkSessionID,
		Theme:            req.Theme,
		Description:      req.Description,
		ScheduledEndTime: req.ScheduledEndTime,
		Prefecture:       req.Prefecture,
		City:             req.City,
		HideTop:          req.HideTop,
		HideReport:       req.HideReport,
		DisableAnalysis:  req.DisableAnalysis,
		Restrictions:     req.Restrictions,
	}); err != nil {
		return nil, err
	}
	return m.talkSessionForManage(ctx, talkSessionID)
}

// talkSessionForManage は管理画面の詳細・編集画面で使うセッション情報を組み立てる。
func (m *manageHandler) talkSessionForManage(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*oas.TalkSessionForManage, error) {
	q := m.GetQueries(ctx)
	row, err := q.GetTalkSessionByID(ctx, talkSessionID.UUID())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, messages.TalkSessionNotFound
	}
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.GetTalkSessionByID")
		return nil, err
	}
	participantCount, err := q.CountTalkSessionParticipants(ctx, talkSessionID.UUID())
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.CountTalkSessionParticipants")
		return nil, err
	}
	hasSurvey := true
	if _, err := q.GetSurveyByTalkSessionID(ctx, talkSessionID.UUID()); errors.Is(err, sql.ErrNoRows) {
		hasSurvey = false
	} else if err != nil {
		utils.HandleError(ctx, err, "GetQueries.GetSurveyByTalkSessionID")
		return nil, err
	}

	ts := row.TalkSession
	restrictions := make([]string, 0, len(ts.Restrictions))
	restrictions = append(restrictions, ts.Restrictions...)
	res := &oas.TalkSessionForManage{
		TalkSessionID: ts.TalkSessionID.String(),
		Theme:         ts.Theme,
		Description:   ts.Description.String,
		Owner: oas.UserForManage{
			UserID:      row.User.UserID.String(),
			DisplayID:   row.User.DisplayID.String,
			DisplayName: row.User.DisplayName.String,
			IconURL:     row.User.IconUrl.String,
		},
		ScheduledEndTime: ts.ScheduledEndTime,
		ThumbnailURL:     ts.ThumbnailUrl.String,
		Hidden:           ts.HideReport.Bool,
		UpdatedAt:        ts.UpdatedAt.Format(time.RFC3339),
		CreatedAt:        ts.CreatedAt.Format(time.RFC3339),
		HideTop:          ts.HideTop,
		HideReport:       ts.HideReport.Bool,
		DisableAnalysis:  ts.DisableAnalysis,
		Restrictions:     restrictions,
		HasSurvey:        hasSurvey,
		ParticipantCount: participantCount,
	}
	if ts.City.Valid {
		res.City = oas.NewOptString(ts.City.String)
	}
	if ts.Prefecture.Valid {
		res.Prefecture = oas.NewOptString(ts.Prefecture.String)
	}
	return res, nil
}

// ManageRegenerateManage implements oas.ManageHandler.
func (m *manageHandler) ManageRegenerateManage(ctx context.Context, req *oas.RegenerateRequest, params oas.ManageRegenerateManageParams) (*oas.RegenerateResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.ManageRegenerate")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}

	tp := req.Type
	tpb, err := tp.MarshalText()
	if err != nil {
		return nil, err
	}
	tpt := string(tpb)

	talkSessionIDStr := params.TalkSessionID
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](talkSessionIDStr)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, err
	}

	switch tpt {
	case "group":
		m.analysisScheduler.ScheduleAnalysis(ctx, talkSessionID, analysis_usecase.ScheduleImmediate)
	case "report":
		if err := m.AnalysisService.GenerateReport(ctx, talkSessionID); err != nil {
			utils.HandleError(ctx, err, "AnalysisService.GenerateReport")
			return nil, err
		}
	case "image":
		// 非同期で画像生成
		go func() {
			if _, err := m.AnalysisService.GenerateImage(ctx, talkSessionID); err != nil {
				utils.HandleError(ctx, err, "AnalysisService.GenerateImage")
				return
			}
		}()
	}

	return &oas.RegenerateResponse{
		Message: "success",
		Status:  "success",
	}, nil
}

// ToggleReportVisibilityManage implements oas.ManageHandler.
func (m *manageHandler) ToggleReportVisibilityManage(ctx context.Context, req *oas.ToggleReportVisibilityRequest, params oas.ToggleReportVisibilityManageParams) (*oas.ToggleReportVisibilityResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.ToggleReportVisibilityManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if err := m.GetQueries(ctx).UpdateTalkSessionHideReport(ctx, model.UpdateTalkSessionHideReportParams{
		TalkSessionID: talkSessionID.UUID(),
		HideReport:    sql.NullBool{Bool: req.Hidden, Valid: true},
	}); err != nil {

		return nil, err
	}

	return &oas.ToggleReportVisibilityResponse{
		Status:  "success",
		Message: "success",
	}, nil
}

// GetReportBySessionId implements oas.ManageHandler.
func (m *manageHandler) GetAnalysisReportManage(ctx context.Context, params oas.GetAnalysisReportManageParams) (*oas.AnalysisReportResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetReportBySessionId")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}

	talkSessionIDStr := params.TalkSessionID
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](talkSessionIDStr)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, err
	}

	res, err := m.AnalysisRepository.FindByTalkSessionID(ctx, talkSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &messages.ReportNotFound
		}
		utils.HandleError(ctx, err, "GetQueries.GetReportByTalkSessionId")
		return nil, err
	}
	if res.Report == nil {
		return nil, &messages.ReportNotFound
	}

	return &oas.AnalysisReportResponse{
		Report: oas.OptString{Value: *res.Report, Set: true},
	}, nil
}

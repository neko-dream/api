package handler

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type manageHandler struct {
	analysis.AnalysisService
	analysis.AnalysisRepository
	*db.DBManager
	authorizationService service.AuthorizationService
	session.TokenManager
}

// GetUserListManage implements oas.ManageHandler.
func (m *manageHandler) GetUserListManage(ctx context.Context, params oas.GetUserListManageParams) ([]oas.UserForManage, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetUserListManage")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}

// GetUserStatsTotalManage implements oas.ManageHandler.
func (m *manageHandler) GetUserStatsTotalManage(ctx context.Context) (*oas.UserStatsResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetUserStatsTotalManage")
	defer span.End()

	if !m.authorizationService.IsKotohiro(m.SetSession(ctx)) {
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

	if !m.authorizationService.IsKotohiro(m.SetSession(ctx)) {
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
	authorizationService service.AuthorizationService,
	tokenManager session.TokenManager,
) oas.ManageHandler {
	return &manageHandler{
		DBManager:            dbm,
		AnalysisService:      ansv,
		AnalysisRepository:   arep,
		authorizationService: authorizationService,
		TokenManager:         tokenManager,
	}
}

// GetTalkSessionListManage implements oas.ManageHandler.
func (m *manageHandler) GetTalkSessionListManage(ctx context.Context, params oas.GetTalkSessionListManageParams) (*oas.TalkSessionListResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetTalkSessionListManage")
	defer span.End()

	if !m.authorizationService.IsKotohiro(m.SetSession(ctx)) {
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

	_ = ctx

	panic("unimplemented")
}

// ManageRegenerateManage implements oas.ManageHandler.
func (m *manageHandler) ManageRegenerateManage(ctx context.Context, req *oas.RegenerateRequest, params oas.ManageRegenerateManageParams) (*oas.RegenerateResponse, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.ManageRegenerate")
	defer span.End()

	if !m.authorizationService.IsKotohiro(m.SetSession(ctx)) {
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
		if err := m.AnalysisService.StartAnalysis(ctx, talkSessionID); err != nil {
			utils.HandleError(ctx, err, "AnalysisService.StartAnalysis")
			return nil, err
		}
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

	if !m.authorizationService.IsKotohiro(m.SetSession(ctx)) {
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

	authCtx, err := m.authorizationService.RequireAuth(m.SetSession(ctx))
	if err != nil {
		return nil, err
	}
	// org所属のユーザであることを確認
	if !authCtx.IsInOrganization() {
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

package handler

import (
	"context"
	"database/sql"
	"text/template"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type manageHandler struct {
	templates *template.Template
	analysis.AnalysisService
	analysis.AnalysisRepository
	*db.DBManager
	session.TokenManager
}

func NewManageHandler(
	dbm *db.DBManager,
	ansv analysis.AnalysisService,
	arep analysis.AnalysisRepository,
	tokenManager session.TokenManager,
) oas.ManageHandler {
	return &manageHandler{
		DBManager:          dbm,
		AnalysisService:    ansv,
		AnalysisRepository: arep,
		TokenManager:       tokenManager,
	}
}

// GetTalkSessionListManage implements oas.ManageHandler.
func (m *manageHandler) GetTalkSessionListManage(ctx context.Context, params oas.GetTalkSessionListManageParams) ([]oas.TalkSessionStats, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetTalkSessionListManage")
	defer span.End()

	claim := session.GetSession(m.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
	}
	if userID == nil {
		return []oas.TalkSessionStats{}, messages.ForbiddenError
	}
	// org所属のユーザであることを確認
	if claim.OrgType == nil {
		return []oas.TalkSessionStats{}, messages.ForbiddenError
	}

	rows, err := m.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:   1000,
		Offset:  0,
		SortKey: sql.NullString{String: "latest", Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.ListTalkSessions")
		return []oas.TalkSessionStats{}, err
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

	return talkSessionStats, nil
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

	claim := session.GetSession(m.SetSession(ctx))
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	if claim.OrgType == nil {
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

	claim := session.GetSession(m.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
	}
	if userID == nil {
		return nil, messages.ForbiddenError
	}
	// org所属のユーザであることを確認
	if claim.OrgType == nil {
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

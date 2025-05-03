package handler

import (
	"context"
	"database/sql"
	"strings"
	"text/template"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/http/templates"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type manageHandler struct {
	templates *template.Template
	analysis.AnalysisService
	*db.DBManager
	session.TokenManager
}

func NewManageHandler(
	dbm *db.DBManager,
	ansv analysis.AnalysisService,
	tokenManager session.TokenManager,
) oas.ManageHandler {
	tmpl, err := template.ParseFS(templates.TemplateFS, "*html")
	if err != nil {
		panic(err)
	}

	return &manageHandler{
		templates:       tmpl,
		DBManager:       dbm,
		AnalysisService: ansv,
		TokenManager:    tokenManager,
	}
}

// ManageRegenerate implements oas.ManageHandler.
func (m *manageHandler) ManageRegenerate(ctx context.Context, req oas.OptManageRegenerateReq) (*oas.ManageRegenerateOK, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.ManageRegenerate")
	defer span.End()

	tp := req.Value.Type
	tpb, err := tp.MarshalText()
	if err != nil {
		return nil, err
	}
	tpt := string(tpb)

	talkSessionIDStr := req.Value.TalkSessionID
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

	return &oas.ManageRegenerateOK{
		Message: oas.OptString{Value: "success", Set: true},
	}, nil
}

// ManageIndex implements oas.ManageHandler.
func (m *manageHandler) ManageIndex(ctx context.Context) (oas.ManageIndexOK, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.ManageIndex")
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
		return oas.ManageIndexOK{}, messages.ForbiddenError
	}
	// org所属のユーザであることを確認
	if claim.OrgType == nil {
		return oas.ManageIndexOK{}, messages.ForbiddenError
	}

	rows, err := m.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:   1000,
		Offset:  0,
		SortKey: sql.NullString{String: "latest", Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.ListTalkSessions")
		return oas.ManageIndexOK{}, err
	}

	var sessions []map[string]any
	for _, row := range rows {
		res := map[string]any{
			"ID":            row.TalkSession.TalkSessionID,
			"Theme":         row.TalkSession.Theme,
			"HideReport":    row.TalkSession.HideReport.Bool,
			"CreatedAt":     row.TalkSession.CreatedAt.Format(time.RFC3339),
			"EndTime":       row.TalkSession.ScheduledEndTime.Format(time.RFC3339),
			"OpinionCount":  row.OpinionCount,
			"DisplayName":   row.User.DisplayName.String,
			"IsOwner":       row.TalkSession.OwnerID == userID.UUID(),
			"VoteCount":     row.VoteCount,
			"VoteUserCount": row.VoteUserCount,
		}

		rr, err := m.GetQueries(ctx).GetGeneratedImages(ctx, row.TalkSession.TalkSessionID)
		if err == nil {
			res["WordCloud"] = rr.WordmapUrl
			res["Tsnc"] = rr.TsncUrl
		}

		sessions = append(sessions, res)
	}

	var html strings.Builder
	data := map[string]any{
		"Sessions": sessions,
	}

	if err := m.templates.ExecuteTemplate(&html, "index.gohtml", data); err != nil {
		utils.HandleError(ctx, err, "templates.ExecuteTemplate")
		return oas.ManageIndexOK{}, err
	}

	return oas.ManageIndexOK{
		Data: strings.NewReader(html.String()),
	}, nil
}

// TalkSessionHideToggle implements oas.ManageHandler.
func (m *manageHandler) TalkSessionHideToggle(ctx context.Context, req oas.OptTalkSessionHideToggleReq) (*oas.TalkSessionHideToggleOK, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.TalkSessionHideToggle")
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

	talkSessionIDStr := req.Value.TalkSessionID
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](talkSessionIDStr)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, err
	}

	hide := req.Value.Hide
	if err := m.GetQueries(ctx).UpdateTalkSessionHideReport(ctx, model.UpdateTalkSessionHideReportParams{
		TalkSessionID: talkSessionID.UUID(),
		HideReport:    sql.NullBool{Bool: hide, Valid: true},
	}); err != nil {
		utils.HandleError(ctx, err, "GetQueries.UpdateTalkSessionHideReport")
		return nil, err
	}

	return &oas.TalkSessionHideToggleOK{
		Status: oas.OptString{Value: "success", Set: true},
	}, nil
}

// GetReportBySessionId implements oas.ManageHandler.
func (m *manageHandler) GetReportBySessionId(ctx context.Context, params oas.GetReportBySessionIdParams) (*oas.GetReportBySessionIdOK, error) {
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

	talkSessionIDStr := params.TalkSessionId
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](talkSessionIDStr)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, err
	}

	res, err := m.GetQueries(ctx).GetReportByTalkSessionId(ctx, talkSessionID.UUID())
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.GetReportByTalkSessionId")
		return nil, err
	}

	return &oas.GetReportBySessionIdOK{
		Report: oas.OptString{Value: res.Report, Set: true},
	}, nil
}

package handler

import (
	"context"
	"database/sql"
	"strings"
	"text/template"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/internal/infrastructure/web/templates"
	"github.com/neko-dream/server/internal/presentation/oas"
)

type manageHandler struct {
	templates *template.Template
	analysis.AnalysisService
	*db.DBManager
}

func NewManageHandler(
	dbm *db.DBManager,
	ansv analysis.AnalysisService,
) oas.ManageHandler {
	tmpl, err := template.ParseFS(templates.TemplateFS, "*.html")
	if err != nil {
		panic(err)
	}

	return &manageHandler{
		templates:       tmpl,
		DBManager:       dbm,
		AnalysisService: ansv,
	}
}

// ManageRegenerate implements oas.ManageHandler.
func (m *manageHandler) ManageRegenerate(ctx context.Context, req oas.OptManageRegenerateReq) (*oas.ManageRegenerateOK, error) {
	tp := req.Value.Type
	tpb, err := tp.MarshalText()
	if err != nil {
		return nil, err
	}
	tpt := string(tpb)

	talkSessionIDStr := req.Value.TalkSessionID
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](talkSessionIDStr)
	if err != nil {
		return nil, err
	}

	switch tpt {
	case "group":
		if err := m.AnalysisService.StartAnalysis(ctx, talkSessionID); err != nil {
			return nil, err
		}
	case "report":
		if err := m.AnalysisService.GenerateReport(ctx, talkSessionID); err != nil {
			return nil, err
		}
	}

	return &oas.ManageRegenerateOK{
		Message: oas.OptString{Value: "success", Set: true},
	}, nil
}

// ManageIndex implements oas.ManageHandler.
func (m *manageHandler) ManageIndex(ctx context.Context) (oas.ManageIndexOK, error) {

	rows, err := m.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:   1000,
		Offset:  0,
		SortKey: sql.NullString{String: "latest", Valid: true},
	})
	if err != nil {
		return oas.ManageIndexOK{}, err
	}
	var sessions []map[string]interface{}
	for _, row := range rows {
		sessions = append(sessions, map[string]interface{}{
			"ID":    row.TalkSessionID,
			"Theme": row.Theme,
		})
	}

	var html strings.Builder
	data := map[string]interface{}{
		"Sessions": sessions,
	}

	if err := m.templates.ExecuteTemplate(&html, "index.html", data); err != nil {
		return oas.ManageIndexOK{}, err
	}

	return oas.ManageIndexOK{
		Data: strings.NewReader(html.String()),
	}, nil
}

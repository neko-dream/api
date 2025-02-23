package handler

import (
	"context"
	"database/sql"
	"strings"
	"text/template"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
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

	rows, err := m.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:   1000,
		Offset:  0,
		SortKey: sql.NullString{String: "latest", Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.ListTalkSessions")
		return oas.ManageIndexOK{}, err
	}
	var sessions []map[string]interface{}
	for _, row := range rows {
		res := map[string]interface{}{
			"ID":    row.TalkSession.TalkSessionID,
			"Theme": row.TalkSession.Theme,
		}

		rr, err := m.GetQueries(ctx).GetGeneratedImages(ctx, row.TalkSession.TalkSessionID)
		if err == nil {
			res["WordCloud"] = rr.WordmapUrl
			res["Tsnc"] = rr.TsncUrl
		}

		sessions = append(sessions, res)
	}

	var html strings.Builder
	data := map[string]interface{}{
		"Sessions": sessions,
	}

	if err := m.templates.ExecuteTemplate(&html, "index.html", data); err != nil {
		utils.HandleError(ctx, err, "templates.ExecuteTemplate")
		return oas.ManageIndexOK{}, err
	}

	return oas.ManageIndexOK{
		Data: strings.NewReader(html.String()),
	}, nil
}

package handler

import (
	"context"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/internal/presentation/oas"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

const (
	manageOpinionDefaultLimit = 20
	manageOpinionMaxLimit     = 100
)

// GetTalkSessionOpinionsManage implements oas.ManageHandler.
func (m *manageHandler) GetTalkSessionOpinionsManage(ctx context.Context, params oas.GetTalkSessionOpinionsManageParams) (*oas.TalkSessionOpinionsForManage, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "manageHandler.GetTalkSessionOpinionsManage")
	defer span.End()

	if admin.GetAdminSession(ctx) == nil {
		return nil, messages.ForbiddenError
	}
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}
	limit, offset := manageOpinionPage(params.Limit, params.Offset)

	q := m.GetQueries(ctx)
	total, err := q.CountOpinionsForManage(ctx, talkSessionID.UUID())
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.CountOpinionsForManage")
		return nil, err
	}
	rows, err := q.ListOpinionsForManage(ctx, model.ListOpinionsForManageParams{
		TalkSessionID: talkSessionID.UUID(),
		RowLimit:      limit,
		RowOffset:     offset,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetQueries.ListOpinionsForManage")
		return nil, err
	}

	opinions := make([]oas.OpinionForManage, 0, len(rows))
	for _, row := range rows {
		opinions = append(opinions, opinionForManage(row))
	}
	return &oas.TalkSessionOpinionsForManage{TotalCount: total, Opinions: opinions}, nil
}

func manageOpinionPage(limitParam, offsetParam oas.OptInt32) (limit, offset int32) {
	limit = manageOpinionDefaultLimit
	if v, ok := limitParam.Get(); ok && v > 0 {
		limit = min(v, manageOpinionMaxLimit)
	}
	if v, ok := offsetParam.Get(); ok && v > 0 {
		offset = v
	}
	return limit, offset
}

func opinionForManage(row model.ListOpinionsForManageRow) oas.OpinionForManage {
	o := row.Opinion
	res := oas.OpinionForManage{
		OpinionID: o.OpinionID.String(),
		Content:   o.Content,
		Author: oas.UserForManage{
			UserID:      row.User.UserID.String(),
			DisplayID:   row.User.DisplayID.String,
			DisplayName: row.User.DisplayName.String,
			IconURL:     row.User.IconUrl.String,
		},
		CreatedAt:           o.CreatedAt,
		AgreeCount:          row.AgreeCount,
		DisagreeCount:       row.DisagreeCount,
		PassCount:           row.PassCount,
		ReportCount:         row.ReportCount,
		UnsolvedReportCount: row.UnsolvedReportCount,
		Deleted:             row.Deleted,
	}
	if o.ParentOpinionID.Valid {
		res.ParentOpinionID = oas.NewOptString(o.ParentOpinionID.UUID.String())
	}
	if o.PictureUrl.Valid {
		res.PictureURL = oas.NewOptString(o.PictureUrl.String)
	}
	if o.ReferenceUrl.Valid {
		res.ReferenceURL = oas.NewOptString(o.ReferenceUrl.String)
	}
	return res
}

package handler

import (
	"context"

	opinionquery "github.com/neko-dream/api/internal/application/query/opinion"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/presentation/oas"
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

	out, err := m.manageOpinions.List(ctx, opinionquery.ManageOpinionsInput{
		TalkSessionID: talkSessionID,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		return nil, err
	}

	opinions := make([]oas.OpinionForManage, 0, len(out.Opinions))
	for _, o := range out.Opinions {
		opinions = append(opinions, opinionForManage(o))
	}
	return &oas.TalkSessionOpinionsForManage{TotalCount: out.TotalCount, Opinions: opinions}, nil
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

func opinionForManage(o opinionquery.ManagedOpinion) oas.OpinionForManage {
	res := oas.OpinionForManage{
		OpinionID: o.OpinionID.String(),
		Content:   o.Content,
		Author: oas.UserForManage{
			DisplayID:   o.Author.DisplayID,
			DisplayName: o.Author.DisplayName,
			IconURL:     o.Author.IconURL,
		},
		CreatedAt:           o.CreatedAt,
		AgreeCount:          o.AgreeCount,
		DisagreeCount:       o.DisagreeCount,
		PassCount:           o.PassCount,
		ReportCount:         o.ReportCount,
		UnsolvedReportCount: o.UnsolvedReportCount,
		Deleted:             o.Deleted,
	}
	if o.Author.UserID != nil {
		res.Author.UserID = o.Author.UserID.String()
	}
	if o.ParentOpinionID != nil {
		res.ParentOpinionID = oas.NewOptString(o.ParentOpinionID.String())
	}
	if o.PictureURL != nil {
		res.PictureURL = oas.NewOptString(*o.PictureURL)
	}
	if o.ReferenceURL != nil {
		res.ReferenceURL = oas.NewOptString(*o.ReferenceURL)
	}
	return res
}

package report_query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/application/query/report_query"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type getOpinionReportQueryInteractor struct {
	*db.DBManager
	talkSessionRep talksession.TalkSessionRepository
}

func NewGetOpinionReportQueryInteractor(
	dbManager *db.DBManager,
	talkSessionRep talksession.TalkSessionRepository,
) report_query.GetOpinionReportQuery {
	return &getOpinionReportQueryInteractor{
		DBManager:      dbManager,
		talkSessionRep: talkSessionRep,
	}
}

// Execute implements report_query.GetOpinionReportQuery.
func (g *getOpinionReportQueryInteractor) Execute(ctx context.Context, input report_query.GetOpinionReportInput) (*report_query.GetOpinionReportOutput, error) {
	ctx, span := otel.Tracer("report_query").Start(ctx, "getOpinionReportQueryInteractor.Execute")
	defer span.End()

	// opinion取得
	opr, err := g.DBManager.GetQueries(ctx).GetOpinionByID(ctx, model.GetOpinionByIDParams{
		OpinionID: input.OpinionID.UUID(),
	})
	if err != nil {
		return nil, err
	}

	// talkSession取得
	talkSession, err := g.talkSessionRep.FindByID(ctx, shared.UUID[talksession.TalkSession](opr.Opinion.TalkSessionID))
	if err != nil {
		return nil, err
	}

	// talkSessionのオーナーと一致するか確認
	if talkSession.OwnerUserID() != input.UserID {
		return nil, messages.TalkSessionNotFound
	}

	// report取得
	reports, err := g.DBManager.GetQueries(ctx).FindReportByOpinionID(ctx, uuid.NullUUID{UUID: opr.Opinion.OpinionID, Valid: true})
	if err != nil {
		return nil, err
	}
	if len(reports) <= 0 {
		return nil, &messages.ReportNotFound
	}

	var op dto.Opinion
	if err := copier.CopyWithOption(&op, opr.Opinion, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utils.HandleError(ctx, err, "copier.CopyWithOption for Opinion")
		return nil, err
	}
	var usr dto.User
	if err := copier.CopyWithOption(&usr, opr.User, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utils.HandleError(ctx, err, "copier.CopyWithOption for User")
		return nil, err
	}

	// reportMapより、reportの情報を取得
	var reportDetailReasons []dto.ReportDetailReason
	reportedUser := make(map[shared.UUID[user.User]]any)
	for _, reportDetail := range reports {
		detailDTO := dto.ReportDetailReason{
			ReportID: shared.UUID[opinion.Report](reportDetail.OpinionReport.OpinionReportID),
			Reason:   opinion.Reason(reportDetail.OpinionReport.Reason).StringJP(),
		}
		if reportDetail.OpinionReport.ReasonText.Valid {
			detailDTO.Content = &reportDetail.OpinionReport.ReasonText.String
		}
		reportedUser[shared.UUID[user.User](reportDetail.OpinionReport.ReporterID)] = struct{}{}
		reportDetailReasons = append(reportDetailReasons, detailDTO)
	}

	return &report_query.GetOpinionReportOutput{
		Report: dto.ReportDetail{
			Opinion:     op,
			User:        usr,
			Reasons:     reportDetailReasons,
			ReportCount: len(reportedUser),
			Status:      reports[0].OpinionReport.Status,
		},
	}, nil
}

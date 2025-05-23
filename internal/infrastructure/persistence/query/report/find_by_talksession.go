package report_query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/application/query/report_query"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type getByTalkSessionQueryInteractor struct {
	*db.DBManager
	talkSessionRep talksession.TalkSessionRepository
}

func NewGetByTalkSessionQueryInteractor(
	dbManager *db.DBManager,
	talkSessionRep talksession.TalkSessionRepository,
) report_query.GetByTalkSessionQuery {
	return &getByTalkSessionQueryInteractor{
		DBManager:      dbManager,
		talkSessionRep: talkSessionRep,
	}
}

func (i *getByTalkSessionQueryInteractor) Execute(ctx context.Context, input report_query.GetByTalkSessionInput) (*report_query.GetByTalkSessionOutput, error) {
	ctx, span := otel.Tracer("report").Start(ctx, "getByTalkSessionQueryInteractor.Execute")
	defer span.End()

	// 操作ユーザーがセッションの作成者かどうかを確認
	talkSession, err := i.talkSessionRep.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "talkSessionRep.FindByID")
		return nil, err
	}
	if talkSession.OwnerUserID() != input.UserID {
		utils.HandleError(ctx, err, "talkSession.UserID != input.UserID")
		return nil, messages.TalkSessionNotFound
	}

	reports, err := i.DBManager.GetQueries(ctx).FindReportsByTalkSession(ctx, model.FindReportsByTalkSessionParams{
		TalkSessionID: uuid.NullUUID{UUID: input.TalkSessionID.UUID(), Valid: true},
		Status:        sql.NullString{String: input.Status, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	if len(reports) == 0 {
		return &report_query.GetByTalkSessionOutput{
			Reports: []dto.ReportDetail{},
		}, nil
	}

	var reportIDs []uuid.UUID
	for _, report := range reports {
		reportIDs = append(reportIDs, report.OpinionReport.OpinionID)
	}

	var reportMap = make(map[shared.UUID[opinion.Opinion]][]model.OpinionReport)
	for _, report := range reports {
		reportMap[shared.UUID[opinion.Opinion](report.OpinionReport.OpinionID)] = append(reportMap[shared.UUID[opinion.Opinion](report.OpinionReport.OpinionID)], report.OpinionReport)
	}

	opinions, err := i.DBManager.GetQueries(ctx).FindOpinionsByOpinionIDs(ctx, reportIDs)
	if err != nil {
		return nil, err
	}

	var reportDetails []dto.ReportDetail
	for _, opRow := range opinions {
		var op dto.Opinion
		if err := copier.CopyWithOption(&op, opRow.Opinion, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
			utils.HandleError(ctx, err, "copier.CopyWithOption for Opinion")
			return nil, err
		}
		var usr dto.User
		if err := copier.CopyWithOption(&usr, opRow.User, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
			utils.HandleError(ctx, err, "copier.CopyWithOption for User")
			return nil, err
		}

		// reportMapより、reportの情報を取得
		var reportDetailReasons []dto.ReportDetailReason
		reportUsers := make(map[shared.UUID[user.User]]any)
		for _, reportDetail := range reportMap[shared.UUID[opinion.Opinion](opRow.Opinion.OpinionID)] {
			detailDTO := dto.ReportDetailReason{
				ReportID: shared.UUID[opinion.Report](reportDetail.OpinionReportID),
				Reason:   opinion.Reason(reportDetail.Reason).StringJP(),
			}
			if reportDetail.ReasonText.Valid {
				detailDTO.Content = lo.ToPtr(reportDetail.ReasonText.String)
			}
			reportUsers[shared.UUID[user.User](reportDetail.ReporterID)] = struct{}{}
			reportDetailReasons = append(reportDetailReasons, detailDTO)
		}

		reportDetails = append(reportDetails, dto.ReportDetail{
			Opinion:     op,
			User:        usr,
			Reasons:     reportDetailReasons,
			ReportCount: len(reportUsers),
			Status:      input.Status,
		})
	}

	return &report_query.GetByTalkSessionOutput{
		Reports: reportDetails,
	}, nil
}

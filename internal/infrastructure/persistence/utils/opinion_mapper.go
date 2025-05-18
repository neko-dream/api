package dto_mapper

import (
	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
)

// ProcessReportedOpinions 通報された意見の内容を置き換える関数
func ProcessReportedOpinions(opinions []dto.SwipeOpinion, reports []model.FindReportByOpinionIDsRow) []dto.SwipeOpinion {
	if len(reports) == 0 {
		return opinions
	}

	reportMap := makeReportMap(reports)

	// 意見を処理
	for i, opinion := range opinions {
		if reportedList, ok := reportMap[opinion.Opinion.OpinionID]; ok {
			opinions[i].Mask(reportedList)
		}
	}

	return opinions
}

// ProcessSingleReportedOpinion 単一の通報された意見を処理する
func ProcessSingleReportedOpinion(opinion *dto.SwipeOpinion, reports []model.FindReportByOpinionIDsRow) {
	if len(reports) == 0 {
		return
	}

	reportMap := makeReportMap(reports)

	if reportedList, ok := reportMap[opinion.Opinion.OpinionID]; ok {
		opinion.Mask(reportedList)
	}
}

// ProcessReportedOpinionsWithRepresentative 代表意見のある通報された意見を処理する
func ProcessReportedOpinionsWithRepresentative(opinions []dto.OpinionWithRepresentative, reports []model.FindReportByOpinionIDsRow) []dto.OpinionWithRepresentative {
	if len(reports) == 0 {
		return opinions
	}

	reportMap := makeReportMap(reports)

	// 意見を処理
	for i, opinion := range opinions {
		if reportedList, ok := reportMap[opinion.Opinion.OpinionID]; ok {
			opinions[i].Mask(reportedList)
		}
	}

	return opinions
}

// ExtractOpinionIDs 意見IDのリストを抽出する
func ExtractOpinionIDs(opinions []dto.SwipeOpinion) []uuid.UUID {
	opinionIDs := make([]uuid.UUID, len(opinions))
	for i, opinion := range opinions {
		opinionIDs[i] = opinion.Opinion.OpinionID.UUID()
	}
	return opinionIDs
}

// ExtractOpinionIDsWithRepresentative 代表意見からIDリストを抽出する
func ExtractOpinionIDsWithRepresentative(opinions []dto.OpinionWithRepresentative) []uuid.UUID {
	opinionIDs := make([]uuid.UUID, len(opinions))
	for i, opinion := range opinions {
		opinionIDs[i] = opinion.Opinion.OpinionID.UUID()
	}
	return opinionIDs
}

// makeReportMap 通報のマップを作成する
func makeReportMap(reports []model.FindReportByOpinionIDsRow) map[shared.UUID[opinion.Opinion]][]model.OpinionReport {
	reportMap := make(map[shared.UUID[opinion.Opinion]][]model.OpinionReport)
	for _, report := range reports {
		opinionID := shared.UUID[opinion.Opinion](report.OpinionReport.OpinionID)
		reportMap[opinionID] = append(reportMap[opinionID], report.OpinionReport)
	}
	return reportMap
}

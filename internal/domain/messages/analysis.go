package messages

var (
	AnalysisReportNotFound = &APIError{
		StatusCode: 404,
		Code:       "ANALYSIS-0000",
		Message:    "レポートが見つかりません。",
	}
	GenerateAnalysisReportFailed = &APIError{
		StatusCode: 500,
		Code:       "ANALYSIS-0001",
		Message:    "レポートの生成に失敗しました。",
	}
	AnalysisReportOpinionNotFound = &APIError{
		StatusCode: 400,
		Code:       "ANALYSIS-0002",
		Message:    "意見がないためレポートを生成できません。",
	}

	AnalysisReportAlreadyFeedbacked = &APIError{
		StatusCode: 400,
		Code:       "ANALYSIS-0003",
		Message:    "既にフィードバック済みのレポートです。",
	}
	AnalysisReportFeedbackFailed = &APIError{
		StatusCode: 500,
		Code:       "ANALYSIS-0004",
		Message:    "レポートのフィードバックに失敗しました。",
	}
)

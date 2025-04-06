package messages

var (
	// ReportNotFound レポートが見つからない
	ReportNotFound = APIError{
		Code:       "REPORT-001",
		Message:    "レポートが見つかりません",
		StatusCode: 404,
	}
)

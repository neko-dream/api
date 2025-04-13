package messages

import "net/http"

var (
	OrganizationAlreadyExists = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "OPINION-001",
		Message:    "この名前の組織はすでに存在します",
	}
	OrganizationForbidden = &APIError{
		StatusCode: http.StatusForbidden,
		Code:       "ORGANIZATION-002",
		Message:    "この操作は許可されていません",
	}
)

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
	UserAlreadyInOrganization = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "ORGANIZATION-003",
		Message:    "ユーザーはすでに組織に参加しています",
	}
	OrganizationNotFound = &APIError{
		StatusCode: http.StatusNotFound,
		Code:       "ORGANIZATION-004",
		Message:    "組織が見つかりません",
	}
	OrganizationInternalServerError = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "ORGANIZATION-005",
		Message:    "組織の操作中にエラーが発生しました。時間をおいて再度お試しください",
	}
)

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
	OrganizationPermissionDenied = &APIError{
		StatusCode: http.StatusForbidden,
		Code:       "ORGANIZATION-006",
		Message:    "操作に必要な権限がありません",
	}
	OrganizationCodeAlreadyExists = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "ORGANIZATION-007",
		Message:    "この組織コードはすでに使用されています",
	}
	OrganizationCodeTooShort = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "ORGANIZATION-008",
		Message:    "組織コードは4文字以上でなければなりません",
	}
	OrganizationCodeInvalid = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "ORGANIZATION-009",
		Message:    "組織コードは英数字で構成される必要があります",
	}
	OrganizationTypeInvalid = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "ORGANIZATION-010",
		Message:    "無効な組織種別です",
	}
)

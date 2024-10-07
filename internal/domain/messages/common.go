package messages

var (
	InternalServerError = &APIError{
		StatusCode: 500,
		Code:       "INTERNAL-0000",
		Message:    "Internal Server Error",
	}

	RequiredParameterError = &APIError{
		StatusCode: 400,
		Code:       "REQ-0001",
		Message:    "必須パラメータが不足しています。",
	}
	BadRequestError = &APIError{
		StatusCode: 400,
		Code:       "REQ-0002",
		Message:    "Bad Request",
	}
)

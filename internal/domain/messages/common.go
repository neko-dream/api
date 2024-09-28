package messages

var (
	InternalServerError = &APIError{
		StatusCode: 500,
		Code:       "INTERNAL-0000",
		Message:    "Internal Server Error",
	}
)

package errors

const (
	ValidationError     = "Bad request"
	InternalServerError = "Something went wrong"
	UnAuthorized        = "Unauthorized"
	UserNotFound        = "User not found"
)

type ApiError struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

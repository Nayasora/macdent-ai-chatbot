package utils

type UserErrorResponse struct {
	StatusCode int
	Message    string
	Details    string
}

func NewUserErrorResponse(statusCode int, message string, details string) *UserErrorResponse {
	return &UserErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

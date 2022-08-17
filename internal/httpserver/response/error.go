package response

const defaultMessage = "There is something wrong with the server."

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func NewErrorResponse(message, detail string) *ErrorResponse {
	return &ErrorResponse{message, detail}
}

func NewDefaultErrorResponse(detail string) *ErrorResponse {
	return &ErrorResponse{defaultMessage, detail}
}

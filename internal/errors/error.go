package errors

const defaultMessage = "There is something wrong with the server."

type handledError struct {
	HttpStatus int    `json:"-"`
	Message    string `json:"message"`
	Detail     string `json:"detail"`
}

func NewHandledError(status int, message, detail string) handledError {
	return handledError{status, message, detail}
}

func NewDefaultHandledError(status int, detail string) handledError {
	return handledError{status, defaultMessage, detail}
}

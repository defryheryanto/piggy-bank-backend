package errors

const defaultMessage = "There is something wrong with the server."

type HandledError struct {
	HttpStatus int    `json:"-"`
	Message    string `json:"message"`
	Detail     string `json:"detail"`
}

func NewHandledError(status int, message, detail string) HandledError {
	return HandledError{status, message, detail}
}

func NewDefaultHandledError(status int, detail string) HandledError {
	return HandledError{status, defaultMessage, detail}
}

func (e HandledError) Error() string {
	return e.Detail
}

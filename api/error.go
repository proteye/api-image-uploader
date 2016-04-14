package api

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(msg string, code int) *Error {
	return &Error{Message: msg, Code: code}
}

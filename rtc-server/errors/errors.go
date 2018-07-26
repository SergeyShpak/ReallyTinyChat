package errors

type ServerError struct {
	Code int
	Hint string
}

func (e ServerError) Error() string {
	return e.Hint
}

func NewServerError(code int, hint string) *ServerError {
	return &ServerError{
		Code: code,
		Hint: hint,
	}
}

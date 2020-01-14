package public

func New(code int, msg string) error {
	return &ErrorMsg{code, msg}
}

type ErrorMsg struct {
	code int
	str  string
}

func (e *ErrorMsg) Error() string {
	return e.str
}

func (e *ErrorMsg) Code() int {
	return e.code
}

package public

import "fmt"

import "errors"

func New(code int, msg string) *ErrorMsg {
	return &ErrorMsg{code: code, msg: msg}
}

type ErrorMsg struct {
	code int
	msg  string
	err  error
}

func (e *ErrorMsg) GetCode() int {
	return e.code
}

func (e *ErrorMsg) GetMsg() string {
	return e.msg
}

func (e *ErrorMsg) GetErr() error {
	return e.err
}

func (e *ErrorMsg) SetErr(err error) *ErrorMsg {
	e.err = err
	return e
}

func (e *ErrorMsg) SetErrText(text string) *ErrorMsg {
	e.err = errors.New(text)
	return e
}

func (e *ErrorMsg) Error() string {
	if e.err != nil {
		return fmt.Sprintf("code:%d, msg:%s, err:%s", e.code, e.msg, e.err.Error())
	}
	return fmt.Sprintf("code:%d, msg:%s", e.code, e.msg)
}

var (
	ErrNoCreateObj      = New(0, "uninitialized object")
	ErrUnKnown          = New(-1, "unknown error")
	ErrDirNoExist       = New(-2, "dir does not exist")
	ErrNewWatcher       = New(-3, "new watcher object failed")
	ErrWatcherAdd       = New(-4, "watcher add dir failed")
	ErrLoadPubKey       = New(-5, "failed to load public key")
	ErrReadAuthFile     = New(-6, "reading authorization file failed")
	ErrDecodeAuthFile   = New(-7, "decode authorization file failed")
	ErrVerifySign       = New(-8, "failed to verify signature")
	ErrUnmarshalLiObj   = New(-9, "unmarshal license object failed")
	ErrGetMachineCode   = New(-10, "failed to get machine code")
	ErrNoMatchProName   = New(-11, "product name does not match")
	ErrLicenseExpired   = New(-12, "license is expired")
	ErrBeforeIssued     = New(-13, "license used before issued")
	ErrNoMatchMachineID = New(-14, "machine id does not match")
)

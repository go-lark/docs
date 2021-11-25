package docs

import (
	"errors"
	"fmt"
)

func GetRawErr(err error) (*Err, bool) {
	e := &Err{}
	ok := errors.As(err, &e)
	return e, ok
}

func newErr(format string, params ...interface{}) *Err {
	return &Err{
		Code: -1,
		Msg:  fmt.Sprintf(format, params...),
	}
}

type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

package docs

import (
	"encoding/json"
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
	Meta Meta   `json:"meta"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Err) Error() string {
	en, err := json.Marshal(e)
	if err == nil {
		return string(en)
	}
	return fmt.Sprintf("code: %d, msg: %s, meta:%#+v", e.Code, e.Msg, e.Meta)
}

type Meta struct {
	RequestID string `json:"request_id"`
	TTLogID   string `json:"tt_log_id"`
	TraceHost string `json:"trace_host"`
	TraceTag  string `json:"trace_tag"`
}

package vo

import (
	"github.com/yuancore/go-zen/zen"
	"github.com/yuancore/zentool/response"
)

type Base struct{}

// Success writes a successful JSON response.
func (b *Base) Success(c zen.Context, code string, data ...interface{}) {
	msg := GetMessage(code)
	c.SecureJSON(200, response.Success(code, msg, data...))
}

// Fail writes an error JSON response.
// bizErr can be *vo.Error, a string error code, or a plain error.
func (b *Base) Fail(c zen.Context, bizErr any, extraErr ...error) {
	var code string
	var err error
	var isBusinessError bool

	switch v := bizErr.(type) {
	case *Error:
		code = v.Code
		err = v.Err
		isBusinessError = code != FAILED && code != ""
	case string:
		code = v
		isBusinessError = code != FAILED && code != ""
	case error:
		if ve, ok := v.(*Error); ok {
			code = ve.Code
			err = ve.Err
			isBusinessError = code != FAILED && code != ""
		} else {
			code = FAILED
			err = v
		}
	default:
		code = FAILED
	}

	if !isBusinessError && len(extraErr) > 0 && extraErr[0] != nil {
		err = extraErr[0]
	}

	msg := GetMessage(code)
	if msg == "" {
		code = FAILED
		msg = GetMessage(FAILED)
	}

	_ = err // err available for future debug-mode logging

	c.SecureJSON(200, response.Fail(code, msg))
}

// Page builds a pagination response object.
func (b *Base) Page(total int64, list interface{}) response.Page {
	return response.Page{
		Total: total,
		Items: list,
	}
}

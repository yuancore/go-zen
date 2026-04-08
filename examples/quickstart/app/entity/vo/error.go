package vo

import "fmt"

// Error 表示业务错误，包含业务码和可选原始错误
type Error struct {
	Code string // 业务错误码
	Err  error  // 原始错误，可为 nil
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return e.Code
}

// 实现错误解包
func (e *Error) Unwrap() error { return e.Err }

// Err 创建一个业务错误（简洁命名）
func Err(code string, err ...error) *Error {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &Error{
		Code: code,
		Err:  e,
	}
}

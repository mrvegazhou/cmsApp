package errorx

import (
	"github.com/pkg/errors"
)

// 定义postgresql错误类型
const (
	POSTGRESQL_FIND_ERR = iota + 1000
	POSTGRESQL_COUNT_ERR
)

// 定义业务错误类型
const (
	HTTP_UNKNOW_ERR = iota + 2000
	HTTP_BIND_PARAMS_ERR
)

type CustomError struct {
	ErrStatus int
	ErrMsg    string
	Err       error
}

// 生成包含另一个error的error
func NewCustomError(status int, message string) error {

	return errors.WithStack(&CustomError{status, message, nil})
}

// 生成一个新的error
func NewCustomErrorWrap(status int, message string, err error) error {

	return errors.WithStack(&CustomError{status, message, err})
}

func (err *CustomError) Error() string {

	if err.Err != nil {
		return err.Err.Error()
	} else {
		return err.ErrMsg
	}
}

package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm: 只支持指向结构体的一级指针")
	ErrNoRows      = errors.New("orm: No data")
)

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm: 非法标签值 %s", pair)
}

func NewErrUnknownField(field string) error {
	return fmt.Errorf("orm: 未知字段 %s", field)
}

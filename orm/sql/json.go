package sql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonCoumn[T any] struct {
	Val T
	// NULL 的问题
	Valid bool
}

func (j *JsonCoumn[T]) Scan(src any) error {
	var bs []byte
	switch data := src.(type) {
	case []byte:
		bs = data
	case string:
		bs = []byte(data)
	case nil:
		return nil
	default:
		return errors.New("不支持的类型")
	}
	return json.Unmarshal(bs, &j.Val)
}

func (j JsonCoumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}

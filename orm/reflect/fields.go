package reflect

import (
	"errors"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}

	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)

	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}

	switch typ.Kind() {
	case reflect.Pointer:
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持的类型")
	}

	numField := typ.NumField()
	m := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)
		if fieldType.IsExported() {
			m[fieldType.Name] = fieldValue.Interface()
		} else {
			m[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}

	return m, nil
}

func SetField(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)

	for val.Type().Kind() == reflect.Pointer {
		val = val.Elem()
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.CanSet() {
		return errors.New("不可修改字段")
	}
	fieldVal.Set(reflect.ValueOf(newValue))
	return nil
}

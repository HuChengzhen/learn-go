package orm

import (
	"learn_geektime_go/orm/internal/errs"
	"reflect"
	"strings"
)

type field struct {
	// 列名
	colName string
}

type model struct {
	tableName string
	fields    map[string]*field
}

// 限制只能用一级指针
func parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)

	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	numField := typ.NumField()
	fieldMap := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		structField := typ.Field(i)
		fieldMap[structField.Name] = &field{
			colName: CamelToSnake(structField.Name),
		}
	}
	return &model{
		tableName: CamelToSnake(typ.Name()),
		fields:    fieldMap,
	}, nil
}

func CamelToSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

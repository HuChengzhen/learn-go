package orm

import (
	"learn_geektime_go/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
)

const (
	tagColumn = "column"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOption) (*Model, error)
}

type Field struct {
	// 列名
	colName string
	typ     reflect.Type
}

type ModelOption func(*Model) error

type Model struct {
	tableName string
	fields    map[string]*Field
}

// var models = map[reflect.Type]*model{}

type registry struct {
	// lock   sync.RWMutex
	models sync.Map
}

func newRegistry() *registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}

	m, err := r.Register(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), nil
}

// func (r *registry) get1(val any) (*model, error) {
// 	typ := reflect.TypeOf(val)
// 	r.lock.RLock()
// 	m, ok := r.models[typ]
// 	r.lock.RUnlock()
// 	if ok {
// 		return m, nil
// 	}

// 	r.lock.Lock()
// 	defer r.lock.Unlock()
// 	m, ok = r.models[typ]
// 	if ok {
// 		return m, nil
// 	}

// 	m, err := r.parseModel(val)
// 	if err != nil {
// 		return nil, err
// 	}
// 	r.models[typ] = m

// 	return m, nil
// }

// 限制只能用一级指针
func (r *registry) Register(entity any, opts ...ModelOption) (*Model, error) {
	typ := reflect.TypeOf(entity)

	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	elemType := typ.Elem()
	numField := elemType.NumField()
	fieldMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		structField := elemType.Field(i)
		pair, err := r.parseTag(structField.Tag)

		if err != nil {
			return nil, err
		}

		columnName := pair[tagColumn]

		if columnName == "" {
			columnName = CamelToSnake(structField.Name)
		}

		fieldMap[structField.Name] = &Field{
			colName: columnName,
			typ:     structField.Type,
		}
	}
	var tableName string

	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}

	if tableName == "" {
		tableName = CamelToSnake(elemType.Name())
	}

	res := &Model{
		tableName: tableName,
		fields:    fieldMap,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, res)
	return res, nil
}

func ModelWithTableName(tableName string) ModelOption {
	return func(m *Model) error {
		m.tableName = tableName

		return nil
	}
}

func ModelWithColumnName(field string, colName string) ModelOption {
	return func(m *Model) error {
		f, ok := m.fields[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		f.colName = colName
		return nil
	}
}

type User struct {
	ID uint64 `orm:"column=id,xxx=bbb`
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}

	pairs := strings.Split(ormTag, ",")

	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}

	return res, nil
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

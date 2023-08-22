package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	model *Model
	where []Predicate
	sb    *strings.Builder
	args  []any
	db    *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: &strings.Builder{},
		db: db,
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	s.model, err = s.db.r.Register(new(T))
	if err != nil {
		return nil, err
	}
	var sb = s.sb
	sb.WriteString("SELECT * FROM ")

	if s.table == "" {
		sb.WriteString("`")
		sb.WriteString(s.model.tableName)
		sb.WriteString("`")
	} else {
		sb.WriteString(s.table)
	}
	if len(s.where) > 0 {
		sb.WriteString("WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		err := s.buildExpression(p)
		if err != nil {
			return nil, err
		}
	}

	sb.WriteString(";")
	return &Query{
		SQL:  sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	build, err := s.Build()
	if err != nil {
		return nil, err
	}
	var db *sql.DB
	db = s.db.db

	rows, err := db.QueryContext(ctx, build.SQL, build.Args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		// 返回error和sql包语义一致。

	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	tp := new(T)

	var vals []any
	vals = make([]any, len(columns))

	for _, column := range columns {
		for _, field := range s.model.fields {
			if field.colName == column {
				val := reflect.New(field.typ)
				vals = append(vals, val.Interface())
			}
		}
	}

	rows.Scan(vals...)

	tpValue := reflect.ValueOf(tp)
	for i, column := range columns {
		for key, field := range s.model.fields {
			if field.colName == column {
				tpValue.Elem().FieldByName(key).Set(reflect.ValueOf(vals[i]).Elem())
			}
		}
	}

	return tp, err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	build, err := s.Build()
	if err != nil {
		return nil, err
	}
	var db *sql.DB
	db = s.db.db
	_, err = db.QueryContext(ctx, build.SQL, build.Args...)
	if err != nil {
		return nil, err
	}
	//for _, r := range rows {

	//}
	return nil, nil
}

func (s *Selector[T]) Where(eq Predicate) QueryBuilder {
	s.where = []Predicate{eq}
	return s
}

func (s *Selector[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column:
		s.sb.WriteByte('`')

		f, ok := s.model.fields[exp.name]
		if !ok {
			return errors.New("orm: 未知字段")
		}
		s.sb.WriteString(f.colName)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.args = append(s.args, exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}
	default:
		return fmt.Errorf("orm: 不支持的表达式 %v", exp)
	}
	return nil
}

func (s *Selector[T]) addArg(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
}

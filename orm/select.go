package orm

import (
	"context"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	sb    *strings.Builder
	args  []any
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var sb = s.sb
	sb.WriteString("SELECT * FROM ")

	if s.table == "" {
		var t T
		typ := reflect.TypeOf(t)
		sb.WriteString("`")
		sb.WriteString(typ.Name())
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
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Where(eq Predicate) QueryBuilder {
	s.where = []Predicate{eq}
	return s
}

func (s *Selector[T]) buildExpression(expr Expression) error {
	// 在这里处理 p
	// p.left
	// p.op
	// p.right

	if expr == nil {
		return nil
	}

	switch expr := expr.(type) {
	case Predicate:
		s.sb.WriteByte('(')
		err := s.buildExpression(expr.left)
		if err != nil {
			return err
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(expr.op.String())
		s.sb.WriteByte(' ')
		err = s.buildExpression(expr.right)
		if err != nil {
			return err
		}
		s.sb.WriteByte(')')
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(expr.name)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.addArg(expr.val)
	}
	return nil
}

func (s *Selector[T]) addArg(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
}

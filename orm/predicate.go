package orm

type op string

const (
	opEq  = "="
	opNot = "NOT"
	opAnd = "AND"
	opOR  = "OR"
)

func (o op) String() string {
	return string(o)
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (p Predicate) expr() {
	//TODO implement me
	panic("implement me")
}

type Column struct {
	name string
}

func (c Column) expr() {
	//TODO implement me
	panic("implement me")
}

func C(name string) Column {
	return Column{
		name: name,
	}
}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: value{val: arg},
	}
}

func Not(predicate Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: predicate,
	}
}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOR,
		right: right,
	}
}

type Expression interface {
	expr()
}

type value struct {
	val any
}

func (v value) expr() {
	//TODO implement me
	panic("implement me")
}

package city

import (
	"fmt"

	"github.com/yoyo-project/yoyo/example/mysql/yoyo/repositories/query"
)

type Query struct {
	n query.Node
}

func (q Query) SQL() (string, []interface{}) {
	cs, ps := q.n.SQL()
	return fmt.Sprintf("WHERE %s", cs), ps
}

func (q Query) Or(in Query) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, in.n},
		Operator: query.Or,
	}}
}
func (q Query) Id(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Id(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdGreaterOrEqual(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdGreaterOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdGreaterThan(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdLessOrEqual(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdLessOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdLessThan(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdNot(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) Name(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Name(in).n},
		Operator: query.And,
	}}
}

func (q Query) NameContains(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NameContains(in).n},
		Operator: query.And,
	}}
}

func (q Query) NameContainsNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NameContainsNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) NameEndsWith(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NameEndsWith(in).n},
		Operator: query.And,
	}}
}

func (q Query) NameEndsWithNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NameEndsWithNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) NameNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NameNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) NameStartsWith(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NameStartsWith(in).n},
		Operator: query.And,
	}}
}

func (q Query) NameStartsWithNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NameStartsWithNot(in).n},
		Operator: query.And,
	}}
}
func Id(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func IdGreaterOrEqual(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.GreaterOrEqual,
			Value:    in,
		},
	}}
}

func IdGreaterThan(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func IdLessOrEqual(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.LessOrEqual,
			Value:    in,
		},
	}}
}

func IdLessThan(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func IdNot(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.NotEquals,
			Value:    in,
		},
	}}
}

func Name(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func NameContains(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%%%s%%'", in),
		},
	}}
}

func NameContainsNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%%%s%%'", in),
		},
	}}
}

func NameEndsWith(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%%%s'", in),
		},
	}}
}

func NameEndsWithNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%%%s'", in),
		},
	}}
}

func NameNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.NotEquals,
			Value:    in,
		},
	}}
}

func NameStartsWith(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%s%%'", in),
		},
	}}
}

func NameStartsWithNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%s%%'", in),
		},
	}}
}

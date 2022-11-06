package state

import (
	"fmt"

	"github.com/yoyo-project/yoyo/example/mysql/yoyo/repositories/query"
)

type Query struct {
	n query.Node
}

func (q Query) SQL() (string, []interface{}) {
	return q.n.SQL()
}

func (q Query) Or(in Query) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, in.n},
		Operator: query.Or,
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

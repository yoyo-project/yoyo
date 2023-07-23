package no_pk_table

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
func (q Query) Col(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Col(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColGreaterOrEqual(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColGreaterOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColGreaterThan(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColLessOrEqual(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColLessOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColLessThan(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColNot(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColNot(in).n},
		Operator: query.And,
	}}
}
func Col(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func ColGreaterOrEqual(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.GreaterOrEqual,
			Value:    in,
		},
	}}
}

func ColGreaterThan(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func ColLessOrEqual(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.LessOrEqual,
			Value:    in,
		},
	}}
}

func ColLessThan(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func ColNot(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.NotEquals,
			Value:    in,
		},
	}}
}

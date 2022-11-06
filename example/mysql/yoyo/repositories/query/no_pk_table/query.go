package no_pk_table

import (
	

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
func (q Query) Col(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Col(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColGreaterOrEqual(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColGreaterOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColGreaterThan(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColLessOrEqual(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColLessOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColLessThan(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) ColNot(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ColNot(in).n},
		Operator: query.And,
	}}
}
func Col(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func ColGreaterOrEqual(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.GreaterOrEqual,
			Value:    in,
		},
	}}
}

func ColGreaterThan(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func ColLessOrEqual(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.LessOrEqual,
			Value:    in,
		},
	}}
}

func ColLessThan(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func ColNot(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "col",
			Operator: query.NotEquals,
			Value:    in,
		},
	}}
}

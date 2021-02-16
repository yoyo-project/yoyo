package state

import (
	"github.com/yoyo-project/yoyo/example/yoyo/repositories/query"
)

type Query query.Node

func (s Query) SQL() (string, []interface{}) {
	return query.Node(s).SQL()
}

func ID(in uint) Query {
	return Query(query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.Equals,
			Value:    in,
		},
	})
}

func IDGreaterThan(in uint) Query {
	return Query(query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.GreaterThan,
			Value:    in,
		},
	})
}

func IDLessThan(in uint) Query {
	return Query(query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.LessThan,
			Value:    in,
		},
	})
}

func IDNot(in uint) Query {
	return Query(query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.NotEquals,
			Value:    in,
		},
	})
}

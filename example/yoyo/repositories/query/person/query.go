package person

import (
	"fmt"
	"time"

	"github.com/dotvezz/yoyo/example/yoyo/repositories/query"
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

func (q Query) ID(in uint) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ID(in).n},
		Operator: query.And,
	}}
}

func (q Query) IDGreaterThan(in uint) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IDGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) IDLessThan(in uint) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IDLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) IDNot(in uint) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IDNot(in).n},
		Operator: query.And,
	}}
}

func ID(in uint) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func IDGreaterThan(in uint) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func IDLessThan(in uint) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func IDNot(in uint) Query {
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
			Value:    fmt.Sprintf("%%%s%%", in),
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

func NameContainsNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "name",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("%%%s%%", in),
		},
	}}
}

func FavoriteColor(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func FavoriteColorContains(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.Like,
			Value:    fmt.Sprintf("%%%s%%", in),
		},
	}}
}

func FavoriteColorNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.NotEquals,
			Value:    in,
		},
	}}
}

func FavoriteColorContainsNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("%%%s%%", in),
		},
	}}
}

func Birthday(in time.Time) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "birthday",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

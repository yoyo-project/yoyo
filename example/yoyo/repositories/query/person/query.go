package person

import (
	"fmt"

	"github.com/yoyo-project/yoyo/example/yoyo/repositories/query"
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
func (q Query) Age(in ufloat64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Age(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeGreaterOrEqual(in ufloat64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeGreaterOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeGreaterThan(in ufloat64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeLessOrEqual(in ufloat64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeLessOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeLessThan(in ufloat64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeNot(in ufloat64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColor(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColor(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorContains(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorContains(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorContainsNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorContainsNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorEndsWith(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorEndsWith(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorEndsWithNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorEndsWithNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorStartsWith(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorStartsWith(in).n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorStartsWithNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorStartsWithNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) Id(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Id(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdGreaterOrEqual(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdGreaterOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdGreaterThan(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdLessOrEqual(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdLessOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdLessThan(in int32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, IdLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) IdNot(in int32) Query {
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
func Age(in ufloat64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func AgeGreaterOrEqual(in ufloat64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.GreaterOrEqual,
			Value:    in,
		},
	}}
}

func AgeGreaterThan(in ufloat64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func AgeLessOrEqual(in ufloat64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.LessOrEqual,
			Value:    in,
		},
	}}
}

func AgeLessThan(in ufloat64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func AgeNot(in ufloat64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.NotEquals,
			Value:    in,
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
			Value:    fmt.Sprintf("'%%%s%%'", in),
		},
	}}
}

func FavoriteColorContainsNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%%%s%%'", in),
		},
	}}
}

func FavoriteColorEndsWith(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%%%s'", in),
		},
	}}
}

func FavoriteColorEndsWithNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%%%s'", in),
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

func FavoriteColorStartsWith(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%s%%'", in),
		},
	}}
}

func FavoriteColorStartsWithNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%s%%'", in),
		},
	}}
}

func Id(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func IdGreaterOrEqual(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.GreaterOrEqual,
			Value:    in,
		},
	}}
}

func IdGreaterThan(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func IdLessOrEqual(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.LessOrEqual,
			Value:    in,
		},
	}}
}

func IdLessThan(in int32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "id",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func IdNot(in int32) Query {
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

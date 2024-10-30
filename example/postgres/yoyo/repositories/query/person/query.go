package person

import (
	"fmt"

	"github.com/yoyo-project/yoyo/example/postgres/yoyo/repositories/query"
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
func (q Query) Age(in float64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Age(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeGreaterOrEqual(in float64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeGreaterOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeGreaterThan(in float64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeLessOrEqual(in float64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeLessOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeLessThan(in float64) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, AgeLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) AgeNot(in float64) Query {
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

func (q Query) FavoriteColorIsNotNull() Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorIsNotNull().n},
		Operator: query.And,
	}}
}

func (q Query) FavoriteColorIsNull() Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, FavoriteColorIsNull().n},
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

func (q Query) HometownId(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, HometownId(in).n},
		Operator: query.And,
	}}
}

func (q Query) HometownIdGreaterOrEqual(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, HometownIdGreaterOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) HometownIdGreaterThan(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, HometownIdGreaterThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) HometownIdLessOrEqual(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, HometownIdLessOrEqual(in).n},
		Operator: query.And,
	}}
}

func (q Query) HometownIdLessThan(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, HometownIdLessThan(in).n},
		Operator: query.And,
	}}
}

func (q Query) HometownIdNot(in uint32) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, HometownIdNot(in).n},
		Operator: query.And,
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

func (q Query) Nickname(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, Nickname(in).n},
		Operator: query.And,
	}}
}

func (q Query) NicknameContains(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NicknameContains(in).n},
		Operator: query.And,
	}}
}

func (q Query) NicknameContainsNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NicknameContainsNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) NicknameEndsWith(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NicknameEndsWith(in).n},
		Operator: query.And,
	}}
}

func (q Query) NicknameEndsWithNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NicknameEndsWithNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) NicknameNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NicknameNot(in).n},
		Operator: query.And,
	}}
}

func (q Query) NicknameStartsWith(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NicknameStartsWith(in).n},
		Operator: query.And,
	}}
}

func (q Query) NicknameStartsWithNot(in string) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, NicknameStartsWithNot(in).n},
		Operator: query.And,
	}}
}
func Age(in float64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func AgeGreaterOrEqual(in float64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.GreaterOrEqual,
			Value:    in,
		},
	}}
}

func AgeGreaterThan(in float64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func AgeLessOrEqual(in float64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.LessOrEqual,
			Value:    in,
		},
	}}
}

func AgeLessThan(in float64) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "age",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func AgeNot(in float64) Query {
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

func FavoriteColorIsNotNull() Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.IsNotNull,
		},
	}}
}

func FavoriteColorIsNull() Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "favorite_color",
			Operator: query.IsNull,
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

func HometownId(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "fk_city_id",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func HometownIdGreaterOrEqual(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "fk_city_id",
			Operator: query.GreaterOrEqual,
			Value:    in,
		},
	}}
}

func HometownIdGreaterThan(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "fk_city_id",
			Operator: query.GreaterThan,
			Value:    in,
		},
	}}
}

func HometownIdLessOrEqual(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "fk_city_id",
			Operator: query.LessOrEqual,
			Value:    in,
		},
	}}
}

func HometownIdLessThan(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "fk_city_id",
			Operator: query.LessThan,
			Value:    in,
		},
	}}
}

func HometownIdNot(in uint32) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "fk_city_id",
			Operator: query.NotEquals,
			Value:    in,
		},
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

func Nickname(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.Equals,
			Value:    in,
		},
	}}
}

func NicknameContains(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%%%s%%'", in),
		},
	}}
}

func NicknameContainsNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%%%s%%'", in),
		},
	}}
}

func NicknameEndsWith(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%%%s'", in),
		},
	}}
}

func NicknameEndsWithNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%%%s'", in),
		},
	}}
}

func NicknameNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.NotEquals,
			Value:    in,
		},
	}}
}

func NicknameStartsWith(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.Like,
			Value:    fmt.Sprintf("'%s%%'", in),
		},
	}}
}

func NicknameStartsWithNot(in string) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "nickname",
			Operator: query.NotLike,
			Value:    fmt.Sprintf("'%s%%'", in),
		},
	}}
}

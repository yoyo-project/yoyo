package template

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yoyo-project/yoyo/internal/schema"
)

const (
	PackageName         = "$PACKAGE_NAME$"
	StdlibImports       = "$STDLIB_IMPORTS$"
	RepositoriesPackage = "$REPOSITORIES_PACKAGE$"
	Type                = "$TYPE$"
	ColumnName          = "$COLUMN_NAME$"
	FuncName            = "$FUNC_NAME$"
	Operator            = "$OPERATOR$"
	Value               = "$VALUE$"

	Equals         = "Equals"
	Not            = "Not"
	Contains       = "Contains"
	ContainsNot    = "ContainsNot"
	StartsWith     = "StartsWith"
	StartsWithNot  = "StartsWithNot"
	EndsWith       = "EndsWith"
	EndsWithNot    = "EndsWithNot"
	GreaterThan    = "GreaterThan"
	GreaterOrEqual = "GreaterOrEqual"
	LessThan       = "LessThan"
	LessOrEqual    = "LessOrEqual"
	Before         = "Before"
	BeforeOrEqual  = "BeforeOrEqual"
	After          = "After"
	AfterOrEqual   = "AfterOrEqual"

	IsNull = "IsNull"
	IsNotNull = "IsNotNull"
)

const (
	QueryFile = `package ` + PackageName + `

import (
	` + StdlibImports + `

	"` + RepositoriesPackage + `query"
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
`
	QueryMethod = `func (q Query) ` + FuncName + `(in ` + Type + `) Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ` + FuncName + `(in).n},
		Operator: query.And,
	}}
}
`
	QueryFunction = `func ` + FuncName + `(in ` + Type + `) Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "` + ColumnName + `",
			Operator: query.` + Operator + `,
			Value:    ` + Value + `,
		},
	}}
}
`
	NullQueryMethod = `func (q Query) ` + FuncName + `() Query {
	return Query{query.Node{
		Children: &[2]query.Node{q.n, ` + FuncName + `().n},
		Operator: query.And,
	}}
}
`
	NullQueryFunction = `func ` + FuncName + `() Query {
	return Query{query.Node{
		Condition: query.Condition{
			Column:   "` + ColumnName + `",
			Operator: query.` + Operator + `,
		},
	}}
}
`
)

func GenerateQueryLogic(col string, column schema.Column) (methods, functions, imports []string) {
	var (
		ops    []operation
		goType string
	)
	switch {
	case column.Datatype.IsTime():
		ops = []operation{
			{Equals},
			{Not},
			{Before},
			{After},
			{BeforeOrEqual},
			{AfterOrEqual},
		}
	case column.Datatype.IsNumeric():
		ops = []operation{
			{Equals},
			{Not},
			{GreaterThan},
			{LessThan},
			{GreaterOrEqual},
			{LessOrEqual},
		}
	case column.Datatype.IsString():
		ops = []operation{
			{Equals},
			{Not},
			{Contains},
			{ContainsNot},
			{StartsWith},
			{StartsWithNot},
			{EndsWith},
			{EndsWithNot},
		}
	case column.Datatype.IsBinary():
		ops = []operation{
			{Equals},
			{Not},
		}
	}

	if column.Nullable {
		ops = append(ops, operation{IsNull}, operation{IsNotNull})
	}

	goType = column.GoTypeString()
	goName := column.ExportedGoName()

	var imports2 []string
	methods, functions, imports2 = buildOperations(goName, col, goType, ops)
	sort.Strings(imports2)

	exists := make(map[string]bool)
	for _, s := range imports2 {
		if _, ok := exists[s]; !ok {
			exists[s] = true
			imports = append(imports, s)
		}
	}

	return methods, functions, imports
}

type operation struct {
	name string
}

func (o operation) funcName(fieldName string) string {
	if o.name == Equals {
		return fieldName
	}
	return fmt.Sprintf("%s%s", fieldName, o.name)
}

func (o operation) val() string {
	switch o.name {
	case Contains, ContainsNot:
		return `fmt.Sprintf("'%%%s%%'", in)`
	case StartsWith, StartsWithNot:
		return `fmt.Sprintf("'%s%%'", in)`
	case EndsWith, EndsWithNot:
		return `fmt.Sprintf("'%%%s'", in)`
	case IsNull, IsNotNull:
		return `nil`
	default:
		return "in"
	}
}

func (o operation) operator() (operator string) {
	switch o.name {
	case Contains:
		operator = "Like"
	case ContainsNot:
		operator = "NotLike"
	case StartsWith:
		operator = "Like"
	case StartsWithNot:
		operator = "NotLike"
	case EndsWith:
		operator = "Like"
	case EndsWithNot:
		operator = "NotLike"
	case Not:
		operator = "NotEquals"
	default:
		operator = o.name
	}
	return operator
}

func (o operation) imports() (imports []string) {
	switch o.name {
	case Contains, ContainsNot, StartsWith, StartsWithNot, EndsWith, EndsWithNot:
		imports = append(imports, `"fmt"`)
	case Before, After, BeforeOrEqual, AfterOrEqual:
		imports = append(imports, `"time"`)
	}
	return imports
}

func buildMethod(fnc, typ string, nullCheck bool) string {
	r := strings.NewReplacer(
		FuncName, fnc,
		Type, typ,
	)

	template := QueryMethod
	if nullCheck {
		template = NullQueryMethod
	}

	return r.Replace(template)
}

func buildFunc(fnc, col, typ, op, val string, nullCheck bool) string {
	r := strings.NewReplacer(
		FuncName, fnc,
		Type, typ,
		ColumnName, col,
		Operator, op,
		Value, val,
	)

	template := QueryFunction
	if nullCheck {
		template = NullQueryFunction
	}

	return r.Replace(template)
}

func buildOperations(field, col, typ string, ops []operation) (methods, functions, imports []string) {
	for _, op := range ops {
		funcName := op.funcName(field)
		val := op.val()
		nullCheck := op.name == IsNotNull || op.name == IsNull
		methods = append(methods, buildMethod(funcName, typ, nullCheck))
		functions = append(functions, buildFunc(funcName, col, typ, op.operator(), val, nullCheck))
		imports = append(imports, op.imports()...)
	}

	return methods, functions, imports
}

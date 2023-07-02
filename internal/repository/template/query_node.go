package template

const NodeFile = `package query

import (
	"fmt"
)

type ComparisonOperator string
type LogicalOperator string

const (
	Equals         ComparisonOperator = "="
	NotEquals      ComparisonOperator = "!="
	Like           ComparisonOperator = "LIKE"
	NotLike        ComparisonOperator = "NOT LIKE"
	GreaterThan    ComparisonOperator = ">"
	GreaterOrEqual ComparisonOperator = ">="
	LessThan       ComparisonOperator = "<"
	LessOrEqual    ComparisonOperator = "<="

	Before        = LessThan
	After         = GreaterThan
	BeforeOrEqual = LessOrEqual
	AfterOrEqual  = GreaterOrEqual

	IsNull    ComparisonOperator = "IS NULL"
	IsNotNull ComparisonOperator = "IS NOT NULL"

	And LogicalOperator = "AND"
	Or  LogicalOperator = "OR"
)

type Condition struct {
	Column   string
	Value    interface{}
	Operator ComparisonOperator
}

func (c Condition) SQL() (string, []interface{}) {
	switch c.Operator {
	case IsNull, IsNotNull:
		return fmt.Sprintf("%s %s", c.Column, c.Operator), []interface{}{}
	default:
		return fmt.Sprintf("%s %s ?", c.Column, c.Operator), []interface{}{c.Value}
	}
}

type Node struct {
	Children  *[2]Node
	Operator  LogicalOperator
	Condition Condition
}

func (n Node) SQL() (s string, args []interface{}) {
	if n.Children != nil {
		sql1, args1 := n.Children[0].SQL()
		sql2, args2 := n.Children[1].SQL()
		if n.Operator == Or {
			sql1 = fmt.Sprintf("(%s)", sql1)
			sql2 = fmt.Sprintf("(%s)", sql2)
		}
		s, args = fmt.Sprintf("%s %s %s", sql1, n.Operator, sql2), append(args1, args2...)
		return s, args
	}

	return n.Condition.SQL()
}
`

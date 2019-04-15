package main

import (
	"fmt"
	"strings"
)

type Row interface {
	GetValue(key string) string
}

type Expression interface {
	Eval(r Row) bool
	String() string
}

/*
 * Query Expression
 */

type Query struct {
	expression Expression
}

func (q *Query) Eval(r Row) bool {
	return q.expression.Eval(r)
}

func (q *Query) String() string {
	return q.expression.String()
}

/*
 * Or Expression
 */

type Or struct {
	left Expression
	right Expression
}

func (o *Or) Eval(r Row) bool {
	return o.left.Eval(r) || o.right.Eval(r)
}

func (o *Or) String() string {
	return fmt.Sprintf("(%v OR %v)", o.left, o.right)
}

/*
 * Group Expression
 */

type Group struct {
	expressions []Expression
}

func (g *Group) Eval(r Row) bool {
	for _, e := range g.expressions {
		if !e.Eval(r) {
			return false
		}
	}

	return true
}

func (g* Group) String() string {
	exprs := make([]string, len(g.expressions))
	for i, e := range g.expressions {
		exprs[i] = e.String()
	}

	return fmt.Sprintf("(%v)", strings.Join(exprs, " AND "))
}

/*
 * Filter Expression
 */

type Filter struct {
	negate bool
	key string
	op string
	value string
}

func (f *Filter) Eval(r Row) bool {
	var result bool
	switch f.op {
	case ":": // TODO
		result = r.GetValue(f.key) == f.value
	default:
		panic(fmt.Sprintf("Unknown op. <%v>", f.String()))
	}

	if f.negate {
		return !result
	} else {
		return result
	}
}

func (f *Filter) String() string {
	return fmt.Sprintf("{negate=%v, key=%v, op=%v, value=%v}", f.negate, f.key, f.op, f.value)
}

/*
 * Unbound Value Expression
 */

type UnboundValue struct {
	value string
}

func (v *UnboundValue) Eval(r Row) bool {
	return true // TODO
}

func (v *UnboundValue) String() string {
	return fmt.Sprintf("{value=%v}", v.value)
}

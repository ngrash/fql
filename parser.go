package main

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/ngrash/fql/parser"
)

func ParseQuery(str string) Expression {
	input := antlr.NewInputStream(str)
	lexer := parser.NewFQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewFQLParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	//p.BuildParseTrees = true
	tree := p.Query()
	listener := NewListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
	return listener.result
}

/*
 * Stack of Expressions
 */

type Stack struct {
	expressions []Expression
}

func NewStack() *Stack {
	return &Stack{expressions: make([]Expression, 0)}
}

func (s *Stack) Push(e Expression) {
	s.expressions = append(s.expressions, e)
}

func (s *Stack) Pop() Expression {
	l := len(s.expressions)
	if l <= 0 {
		panic("Cannot pop from empty expression stack")
	}

	e := s.expressions[l-1]
	s.expressions = s.expressions[:l-1]
	return e
}

/*
 * Listener Implementation
 */

type Listener struct {
	*parser.BaseFQLListener
	stack  *Stack
	debug  bool
	result Expression
}

func NewListener() *Listener {
	return &Listener{stack: NewStack()}
}

func (l *Listener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	if l.debug {
		fmt.Printf("[enter] %v", debugOutput(l, ctx))
	}
}

func (l *Listener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	if l.debug {
		fmt.Printf("[exit]  %v", debugOutput(l, ctx))
	}
}

func debugOutput(l *Listener, ctx antlr.ParserRuleContext) string {
	return fmt.Sprintf("%T: '%v', stack depth: %v\n", ctx, ctx.GetText(), len(l.stack.expressions))
}

func (l *Listener) EnterQuery(ctx *parser.QueryContext) {
	l.stack.Push(&Query{})
}

func (l *Listener) ExitQuery(ctx *parser.QueryContext) {
	expression := l.stack.Pop()
	query := l.stack.Pop().(*Query)
	query.expression = expression
	l.result = query
}

func (l *Listener) ExitExpression(ctx *parser.ExpressionContext) {
	if ctx.OR() != nil {
		l.stack.Push(&Or{
			right: l.stack.Pop(),
			left:  l.stack.Pop()})
	} else if ctx.UnboundValue() == nil && ctx.Filter() == nil {
		num_expressions := len(ctx.AllExpression())
		expressions := make([]Expression, num_expressions)
		for i := 0; i < num_expressions; i++ {
			expressions[num_expressions-1-i] = l.stack.Pop()
		}

		l.stack.Push(&Group{expressions: expressions})
	}
}

func (l *Listener) EnterUnboundValue(ctx *parser.UnboundValueContext) {
	l.stack.Push(&UnboundValue{value: ctx.GetText()})
}

func (l *Listener) EnterFilter(ctx *parser.FilterContext) {
	l.stack.Push(&Filter{
		negate: ctx.NOT() != nil,
		key:    ctx.Key().GetText(),
		op:     ctx.Op().GetText(),
		value:  ctx.Value().GetText()})
}

package main

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"fql/parser"
	"os"
	"fmt"
)

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
		panic("Cannot pop from empty Expression stack")
	}

	e := s.expressions[l-1]
	s.expressions = s.expressions[:l-1]
	return e
}

/*
 * Listener Implementation
 */

type TreeShapeListener struct {
	*parser.BaseFQLListener
	stack *Stack
	debug bool
	result Expression
}

func NewTreeShapeListener() *TreeShapeListener {
	return &TreeShapeListener{stack: NewStack()}
}

func (this *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	if this.debug {
		fmt.Printf("[enter] %v", debugOutput(this, ctx))
	}
}

func (this *TreeShapeListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	if this.debug {
		fmt.Printf("[exit]  %v", debugOutput(this, ctx))
	}
}

func debugOutput(listener *TreeShapeListener, ctx antlr.ParserRuleContext) string {
	return fmt.Sprintf("%T: '%v', stack depth: %v\n", ctx, ctx.GetText(), len(listener.stack.expressions))
}

func (this *TreeShapeListener) EnterQuery(ctx *parser.QueryContext) {
	this.stack.Push(&Query{})
}

func (this *TreeShapeListener) ExitQuery(ctx *parser.QueryContext) {
	expression := this.stack.Pop()
	query := this.stack.Pop().(*Query)
	query.expression = expression
	this.result = query
}

func (this *TreeShapeListener) ExitExpression(ctx *parser.ExpressionContext) {
	if ctx.OR() != nil {
		this.stack.Push(&Or{
			right: this.stack.Pop(),
			left: this.stack.Pop()})
	} else if ctx.UnboundValue() == nil && ctx.Filter() == nil {
		num_expressions := len(ctx.AllExpression())
		expressions := make([]Expression, num_expressions)
		for i := 0; i < num_expressions; i++ {
			expressions[num_expressions-1-i] = this.stack.Pop()
		}

		this.stack.Push(&Group{expressions: expressions})
	}
}

func (this *TreeShapeListener) EnterUnboundValue(ctx *parser.UnboundValueContext) {
	this.stack.Push(&UnboundValue{value: ctx.GetText()})
}

func (this *TreeShapeListener) EnterFilter(ctx *parser.FilterContext) {
	this.stack.Push(&Filter{
		negate: ctx.NOT() != nil,
		key: ctx.Key().GetText(),
		op: ctx.Op().GetText(),
		value: ctx.Value().GetText()})
}

/*
 * Main Method
 */

func main() {
	input := antlr.NewInputStream(os.Args[1])
	lexer := parser.NewFQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer,0)
	p := parser.NewFQLParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Query()
	listener := NewTreeShapeListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
	fmt.Println(listener.result)
}

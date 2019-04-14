package main

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"fql/parser"
	"os"
	"fmt"
)

type TreeShapeListener struct {
	*parser.BaseFQLListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (this *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func main() {
	input, _ := antlr.NewFileStream(os.Args[1])
	lexer := parser.NewFQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer,0)
	p := parser.NewFQLParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Query()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}

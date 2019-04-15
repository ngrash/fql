fql: fql.go ast.go compare.go parser/fql_parser.go
	go build

parser/fql_parser.go: FQL.g4
	antlr4 -Dlanguage=Go -o parser FQL.g4

install:
	go install

uninstall:
	rm -f $(shell go env GOPATH)/bin/fql

clean:
	rm -rf ./parser
	rm -f ./fql

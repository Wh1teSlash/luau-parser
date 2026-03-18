# Luau Parser

A high-performance, lossless Luau parser engineered for speed and robust error handling.

Installation:

```bash
go get github.com/Wh1teSlash/luau-parser
```

Usage:
```go
package main

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
	"github.com/Wh1teSlash/luau-parser/parser"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

func main() {
	input := `
	local x = 5
	local y = 6
	print(x + 12 - y)
	`

	lexer := lexer.New(input)
	factory := ast.NewFactory()
	parser := parser.New(lexer, factory)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Println("Parser error:", err)
		}
		return
	}

	// Print parsed AST
	treePrinter := visitors.NewTreePrinter()
	fmt.Println(treePrinter.Print(program))
	// Print parsed code reconstructed from AST
	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))

	// Free Memory
	factory.Reset()

	pos := ast.Position{Line: 1, Column: 1}
	valX := factory.Literal(pos, "number", int64(5))
	stmt1 := factory.LocalAssignment(
		pos,
		[]string{"x"},
		[]ast.TypeNode{nil},
		[]ast.Expr{valX},
	)

	valY := factory.Literal(pos, "number", int64(6))
	stmt2 := factory.LocalAssignment(
		pos,
		[]string{"y"},
		[]ast.TypeNode{nil},
		[]ast.Expr{valY},
	)

	identPrint := factory.Identifier(pos, "print")
	identX := factory.Identifier(pos, "x")
	identY := factory.Identifier(pos, "y")
	val12 := factory.Literal(pos, "number", int64(12))

	addExpr := factory.BinaryOp(pos, identX, "+", val12)
	subExpr := factory.BinaryOp(pos, addExpr, "-", identY)

	callPrint := factory.FunctionCall(pos, identPrint, []ast.Expr{subExpr})

	stmt3 := factory.ExpressionStatement(pos, callPrint)

	astProgram := factory.Program(pos, []ast.Stmt{stmt1, stmt2, stmt3})
	// Print AST made with factory
	fmt.Println(treePrinter.Print(astProgram))
	// Print reconstructed code from AST made with factory
	fmt.Println(printer.Print(astProgram))
}
```

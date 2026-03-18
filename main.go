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

	treePrinter := visitors.NewTreePrinter()
	fmt.Println(treePrinter.Print(program))
	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))

	factory.Reset()

	pos := ast.Position{Line: 1, Column: 1}
	valX := factory.Literal(pos, "number", int64(5))
	stmt1 := factory.LocalAssignment(
		pos,
		[]string{"x"},
		[]ast.Expr{valX},
		ast.WithTypes(factory.PrimitiveType(pos, "number")),
	)

	valY := factory.Literal(pos, "number", int64(6))
	stmt2 := factory.LocalAssignment(
		pos,
		[]string{"y"},
		[]ast.Expr{valY},
		ast.WithTypes(factory.PrimitiveType(pos, "number")),
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
	fmt.Println(treePrinter.Print(astProgram))
	fmt.Println(printer.Print(astProgram))
}

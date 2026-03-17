package main

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

func main() {
	program := &ast.Program{
		BaseNode: ast.BaseNode{Position: ast.Position{Line: 1, Column: 1}},
		Body: []ast.Stmt{
			&ast.Assignment{
				BaseNode: ast.BaseNode{Position: ast.Position{Line: 1, Column: 1}},
				Target: &ast.Identifier{
					BaseNode: ast.BaseNode{Position: ast.Position{Line: 1, Column: 1}},
					Name:     "x",
				},
				Value: &ast.BinaryOp{
					BaseNode: ast.BaseNode{Position: ast.Position{Line: 1, Column: 5}},
					Left: &ast.Literal{
						BaseNode: ast.BaseNode{Position: ast.Position{Line: 1, Column: 5}},
						Type:     "number",
						Value:    5,
					},
					Op: "+",
					Right: &ast.Literal{
						BaseNode: ast.BaseNode{Position: ast.Position{Line: 1, Column: 9}},
						Type:     "number",
						Value:    3,
					},
				},
			},
		},
	}

	printer := visitors.NewPrinter()
	output := printer.Print(program)

	fmt.Print(output)

	l := lexer.New(`x = 5 + 3`)
	for tok := l.NextToken(); tok.Type != lexer.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}

# Luau Parser

A high-performance, zero-allocation (Arena-backed), lossless Luau parser engineered for speed and robust error handling.

# Installation

```bash
go get github.com/Wh1teSlash/luau-parser
```

# Basic Usage: Parsing Luau Code

The standard workflow involves passing your source code to the Lexer, injecting an AST Factory into the Parser, and generating the tree.

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
	// 1. Initialize Lexer and Arena Factory
	l := lexer.New(input)
	factory := ast.NewFactory()

	// 2. Parse the program
	p := parser.New(l, factory)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Println("Parser error:", err)
		}
		return
	}

	// 3. View the AST
	treePrinter := visitors.NewTreePrinter()
	fmt.Println(treePrinter.Print(program))

	// 4. Reconstruct the Luau code
	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))

	// 5. Instantly free the memory pool when done
	factory.Reset()
}
```

# Advanced: Programmatic AST Construction

Because luau-parser uses a high-performance NodeFactory, you can rapidly construct, manipulate, or test AST nodes programmatically without needing to parse raw text.

```Go

package main

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

func main() {
	factory := ast.NewFactory()
	pos := ast.Position{Line: 1, Column: 1}

	// Build: local x = 5
	valX := factory.Literal(pos, "number", int64(5))
	stmt1 := factory.LocalAssignment(pos, []string{"x"}, []ast.TypeNode{nil}, []ast.Expr{valX})

	// Build: local y = 6
	valY := factory.Literal(pos, "number", int64(6))
	stmt2 := factory.LocalAssignment(pos, []string{"y"}, []ast.TypeNode{nil}, []ast.Expr{valY})

	// Build: print(x + 12 - y)
	identPrint := factory.Identifier(pos, "print")
	identX := factory.Identifier(pos, "x")
	identY := factory.Identifier(pos, "y")
	val12 := factory.Literal(pos, "number", int64(12))
	addExpr := factory.BinaryOp(pos, identX, "+", val12)
	subExpr := factory.BinaryOp(pos, addExpr, "-", identY)
	callPrint := factory.FunctionCall(pos, identPrint, []ast.Expr{subExpr})
	stmt3 := factory.ExpressionStatement(pos, callPrint)

	// Combine into a Program Node
	program := factory.Program(pos, []ast.Stmt{stmt1, stmt2, stmt3})

	// Print reconstructed code
	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))
}
```

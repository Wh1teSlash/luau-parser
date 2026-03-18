# Luau Parser
A high-performance, zero-allocation (Arena-backed), lossless Luau parser engineered for speed and robust error handling.

# Installation
```bash
go get github.com/Wh1teSlash/luau-parser
```

# Usage
## Parsing Luau Code
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

## Programmatic AST Construction
Because luau-parser uses a high-performance NodeFactory, you can rapidly construct, manipulate, or test AST nodes programmatically without needing to parse raw text.

For simple nodes, use the standard factory methods:
```go
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
	stmt1 := factory.LocalAssignment(pos, []string{"x"}, []ast.Expr{valX})

	// Build: print(x)
	identPrint := factory.Identifier(pos, "print")
	identX := factory.Identifier(pos, "x")
	callPrint := factory.FunctionCall(pos, identPrint, []ast.Expr{identX})
	stmt2 := factory.ExpressionStatement(pos, callPrint)

	program := factory.Program(pos, []ast.Stmt{stmt1, stmt2})

	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))
}
```

## Complex Nodes (Functional Options)
For nodes with optional fields — like type annotations, generics, return types, or else branches — the factory uses the idiomatic Go Functional Options pattern. Only pass the options you need; everything else defaults to a sensible zero value.

```go
package main

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

func main() {
	factory := ast.NewFactory()
	pos := ast.Position{Line: 1, Column: 1}

	// 1. Typed local assignment
	// local x: number = 5
	valX := factory.Literal(pos, "number", int64(5))
	stmt1 := factory.LocalAssignment(pos, []string{"x"}, []ast.Expr{valX},
		ast.WithTypes(factory.PrimitiveType(pos, "number")),
	)

	// 2. Generic exported type alias
	// export type Map<K, V> = table
	tableType := factory.TableType(pos, nil)
	typeAlias := factory.TypeAlias(pos, "Map", tableType,
		ast.AsExported(),
		ast.WithTypeGenerics("K", "V"),
	)

	// 3. If statement with an else block
	// if true then ... else ... end
	cond := factory.Literal(pos, "boolean", true)
	thenBlock := factory.Block(pos, []ast.Stmt{})
	elseBlock := factory.Block(pos, []ast.Stmt{})
	ifStmt := factory.IfStatement(pos, cond, thenBlock,
		ast.WithStmtElse(elseBlock),
	)

	// 4. Generic function definition
	// function doMath<T>(val: T): T ... end
	param := factory.Parameter("val", factory.PrimitiveType(pos, "T"))
	funcDef := factory.FunctionDef(pos, "doMath", thenBlock,
		ast.WithDefGenerics("T"),
		ast.WithDefParams(param),
		ast.WithDefReturnType(factory.PrimitiveType(pos, "T")),
	)

	// 5. Generic function expression
	// function<T>(val: T): T ... end
	funcExpr := factory.FunctionExpr(pos, []*ast.Parameter{param}, thenBlock,
		ast.WithExprGenerics("T"),
		ast.WithExprReturnType(factory.PrimitiveType(pos, "T")),
	)

	// 6. Local function with attribute and generic
	// @native local function doMath<T>(val: T): T ... end
	attr := factory.Attribute(pos, "native")
	localFunc := factory.LocalFunction(pos, "doMath", []*ast.Parameter{param}, thenBlock,
		ast.WithLocalGenerics("T"),
		ast.WithLocalReturnType(factory.PrimitiveType(pos, "T")),
		ast.WithLocalAttributes(attr),
	)

	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(factory.Program(pos, []ast.Stmt{
		stmt1, typeAlias, ifStmt, funcDef,
		factory.ExpressionStatement(pos, funcExpr),
		localFunc,
	})))
}
```

## Transforming the AST
luau-parser provides a `Transformer` interface and a `BaseTransformer` base implementation for walking and rewriting the AST. Embed `BaseTransformer` in your own struct and override only the node types you want to change — everything else is passed through and recursively walked by default.

```go
package main

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
	"github.com/Wh1teSlash/luau-parser/parser"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

// RenameTransformer renames all identifiers matching From to To.
type RenameTransformer struct {
	ast.BaseTransformer
	From, To string
}

func (r *RenameTransformer) TransformIdentifier(node *ast.Identifier) ast.Expr {
	if node.Name == r.From {
		node.Name = r.To
	}
	return node
}

func main() {
	input := `
	local x = 5
	print(x + 1)
	`
	l := lexer.New(input)
	factory := ast.NewFactory()
	p := parser.New(l, factory)
	program := p.ParseProgram()

	// Rename every identifier "x" to "value"
	t := &RenameTransformer{From: "x", To: "value"}
	t.TransformProgram(program)

	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))
	// local value = 5
	// print(value + 1)

	factory.Reset()
}
```

The transformer mutates nodes in place, which fits the arena model — no extra allocations. If you need to produce a new node instead of mutating (e.g. changing a node's type entirely), use the factory inside your override:

```go
type ZeroLiteralsTransformer struct {
	ast.BaseTransformer
	factory *ast.NodeFactory
}

func (t *ZeroLiteralsTransformer) TransformLiteral(node *ast.Literal) ast.Expr {
	if node.Type == "number" {
		// replace every number literal with zero
		return t.factory.Literal(node.Pos(), "number", int64(0))
	}
	return node
}
```

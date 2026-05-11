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

## Rewriting the AST
`ast.Rewrite` and `ast.RewriteProgram` provide a functional, copy-on-change tree rewriter. You supply an `EditFunc` that receives each node (after its children have already been rewritten) along with the full ancestor chain. Return `nil` to leave the node unchanged, or return a replacement to substitute it. Because only modified paths are copied, the original tree is never mutated.

```go
type EditFunc func(node Node, parents []Node) Node

func Rewrite(root Node, fn EditFunc) Node
func RewriteProgram(root *Program, fn EditFunc) *Program  // type-safe convenience
```

**Walk without rewriting** — `EditFunc` doubles as a read-only visitor when you always return `nil`:

```go
func CollectCalls(root ast.Node) []*ast.FunctionCall {
	var calls []*ast.FunctionCall
	ast.Rewrite(root, func(node ast.Node, _ []ast.Node) ast.Node {
		if fc, ok := node.(*ast.FunctionCall); ok {
			calls = append(calls, fc)
		}
		return nil
	})
	return calls
}
```

**Rename an identifier across a whole program:**

```go
func Rename(root *ast.Program, from, to string) *ast.Program {
	return ast.RewriteProgram(root, func(node ast.Node, _ []ast.Node) ast.Node {
		id, ok := node.(*ast.Identifier)
		if !ok || id.Name != from {
			return nil
		}
		cp := *id
		cp.Name = to
		return &cp
	})
}
```

**Use `parents` for context-sensitive rewrites** — the slice runs from the root (`parents[0]`) down to the immediate parent (`parents[len-1]`):

```go
// Wrap every number literal that is directly inside a return with tostring()
ast.RewriteProgram(root, func(node ast.Node, parents []ast.Node) ast.Node {
	lit, ok := node.(*ast.Literal)
	if !ok || lit.Type != "number" || len(parents) == 0 {
		return nil
	}
	if _, ok := parents[len(parents)-1].(*ast.ReturnStatement); !ok {
		return nil
	}
	return &ast.FunctionCall{
		Function: &ast.Identifier{Name: "tostring"},
		Args:     []ast.Expr{lit},
	}
})
```

**Constant folding** — because children are rewritten before the parent sees them, multi-level simplification happens in a single pass:

```go
ast.RewriteProgram(root, func(node ast.Node, _ []ast.Node) ast.Node {
	bin, ok := node.(*ast.BinaryOp)
	if !ok {
		return nil
	}
	l, lOk := bin.Left.(*ast.Literal)
	r, rOk := bin.Right.(*ast.Literal)
	if !lOk || !rOk || l.Type != "number" || r.Type != "number" {
		return nil
	}
	lv, rv := l.Value.(float64), r.Value.(float64)
	var result float64
	switch bin.Op {
	case "+":
		result = lv + rv
	case "-":
		result = lv - rv
	case "*":
		result = lv * rv
	case "/":
		if rv == 0 {
			return nil
		}
		result = lv / rv
	default:
		return nil
	}
	return &ast.Literal{Type: "number", Value: result}
})
```

## Transforming the AST
luau-parser provides a `Transformer` interface and a `BaseTransformer` base implementation for walking and **mutating** the AST in place. Unlike `Rewrite`, the transformer modifies nodes directly — no copies are made, which fits the arena model perfectly.

Because Go does not have virtual method dispatch through struct embedding, you must pass the outer concrete type into the traversal yourself. The recommended pattern is to add a `self Transformer` field and initialise it in a constructor:

```go
package main

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
	"github.com/Wh1teSlash/luau-parser/parser"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

// RenameTransformer renames every identifier matching From to To.
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

	t := &RenameTransformer{From: "x", To: "value"}
	for i, stmt := range program.Body {
		program.Body[i] = t.TransformStmt(stmt)
	}

	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))
	// local value = 5
	// print(value + 1)

	factory.Reset()
}
```

If you need to replace a node with one of a different type rather than mutating it, allocate via the factory inside your override:

```go
type ZeroLiteralsTransformer struct {
	ast.BaseTransformer
	factory *ast.NodeFactory
}

func (t *ZeroLiteralsTransformer) TransformLiteral(node *ast.Literal) ast.Expr {
	if node.Type == "number" {
		return t.factory.Literal(node.Pos(), "number", int64(0))
	}
	return node
}
```

# Example Projects

- [Luau minifier/beautifier (luau-squeeze)](https://github.com/Wh1teSlash/luau-squeeze)

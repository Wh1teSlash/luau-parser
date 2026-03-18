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

## Complex Nodes (Functional Builders)

For complex nodes with optional branches, parameters, or types (like Functions, If Statements, and Type Aliases), the factory uses the idiomatic Go Functional Options pattern. This keeps your construction logic clean and strict.

```Go
package main

import (
    "fmt"
    "github.com/Wh1teSlash/luau-parser/ast"
)

func main() {
    factory := ast.NewFactory()
    pos := ast.Position{Line: 1, Column: 1}

    // 1. Building a complex generic exported Type Alias
    // export type Map<K, V> = table
    tableType := factory.TableType(pos, nil)
    typeAlias := factory.TypeAlias(pos, "Map", tableType, 
        ast.AsExported(), 
        ast.WithTypeGenerics("K", "V"),
    )

    // 2. Building an If Statement with an Else block
    // if true then ... else ... end
    cond := factory.Literal(pos, "boolean", true)
    thenBlock := factory.Block(pos, []ast.Stmt{})
    elseBlock := factory.Block(pos, []ast.Stmt{})
    
    ifStmt := factory.IfStatement(pos, cond, thenBlock, 
        ast.WithStmtElse(elseBlock),
    )

    // 3. Building a Function Definition
    // function doMath<T>(val: T): T ... end
    param := factory.Parameter("val", factory.PrimitiveType(pos, "T"))
    funcDef := factory.FunctionDef(pos, "doMath", thenBlock,
        ast.WithDefGenerics("T"),
        ast.WithDefParams(param),
        ast.WithDefReturnType(factory.PrimitiveType(pos, "T")),
    )
}
```

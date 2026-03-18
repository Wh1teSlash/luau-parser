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

	"github.com/Wh1teSlash/luau-parser/lexer"
	"github.com/Wh1teSlash/luau-parser/parser"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

func main() {
	input := `print("Hello World!")`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	errors := parser.Errors()

	if len(errors)) > 0 {
		for _, err := range errors {
			fmt.Println("Parser error:", err)
		}
		return
	}

	// print AST
	treePrinter := visitors.NewTreePrinter()
	fmt.Println(treePrinter.Print(program))

	// Print code based of AST
	printer := visitors.NewPrinter()
	fmt.Println(printer.Print(program))
}
```

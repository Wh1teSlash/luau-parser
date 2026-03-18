package main

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/lexer"
	"github.com/Wh1teSlash/luau-parser/parser"
	"github.com/Wh1teSlash/luau-parser/visitors"
)

func main() {
	input := `
	print("PROMETHEUS Benchmark")
print("Based On IronBrew Benchmark")
local Iterations = 100000
print("Iterations: " .. tostring(Iterations))

print("CLOSURE testing.")
local Start = os.clock()
local TStart = Start
for Idx = 1, Iterations do
    (function()
        if not true then
            print("Hey gamer.")
        end
    end)()
end
print("Time:", os.clock() - Start .. "s")

print("SETTABLE testing.")
Start = os.clock()
local T = {}
for Idx = 1, Iterations do
    T[tostring(Idx)] = "EPIC GAMER " .. tostring(Idx)
end

print("Time:", os.clock() - Start .. "s")

print("GETTABLE testing.")
Start = os.clock()
for Idx = 1, Iterations do
    T[1] = T[tostring(Idx)]
end

print("Time:", os.clock() - Start .. "s")
print("Total Time:", os.clock() - TStart .. "s")
	`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
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
}

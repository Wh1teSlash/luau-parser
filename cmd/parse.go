/*
Copyright © 2026 WhiteSlash whiteslashdev@gmail.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
	"github.com/Wh1teSlash/luau-parser/parser"
	"github.com/Wh1teSlash/luau-parser/visitors"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	inputCode  string
)

var parseCmd = &cobra.Command{
	Use:   "parse [file]",
	Short: "Parse a Luau file or input and print the AST",
	Long: `Parse a Luau source file or inline code and print the resulting AST.

You can provide input in three ways:
  1. Pass a file path as an argument:        luau-parser parse script.luau
  2. Pass inline code with --input flag:     luau-parser parse --input "local x = 1"
  3. Pipe code via stdin:                    echo "local x = 1" | luau-parser parse

Optionally, write the output to a file using --output:
  luau-parser parse script.luau --output ast.txt`,

	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := resolveInput(cmd, args)
		if err != nil {
			return err
		}

		result, err := parseInput(input)
		if err != nil {
			return err
		}

		return writeOutput(result, outputFile)
	},
}

func resolveInput(_ *cobra.Command, args []string) (string, error) {
	if inputCode != "" {
		return inputCode, nil
	}

	if len(args) > 0 {
		filePath := args[0]

		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file %q: %w", filePath, err)
		}

		return string(data), nil
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat stdin: %w", err)
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var sb strings.Builder
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			sb.WriteString(scanner.Text())
			sb.WriteString("\n")
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}

		return sb.String(), nil
	}

	return "", fmt.Errorf(
		"no input provided, use a file argument, --input flag, or pipe code via stdin\n\n" +
			"Run 'luau-parser parse --help' for usage details",
	)
}

func parseInput(input string) (string, error) {
	l := lexer.New(input)
	factory := ast.NewFactory()
	p := parser.New(l, factory)
	node := p.ParseProgram()

	printer := visitors.NewTreePrinter()
	result := printer.Print(node)

	return result, nil
}

func writeOutput(result string, outputFile string) error {
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(result), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output to %q: %w", outputFile, err)
		}

		fmt.Fprintf(os.Stderr, "AST written to %s\n", outputFile)
		return nil
	}

	fmt.Print(result)
	return nil
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(
		&outputFile,
		"output", "o",
		"",
		"Write AST output to a file instead of stdout",
	)

	parseCmd.Flags().StringVarP(
		&inputCode,
		"input", "i",
		"",
		"Luau source code to parse inline (e.g. --input \"local x = 1\")",
	)
}

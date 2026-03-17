package parser

import (
	"testing"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser found %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %v", msg)
	}
	t.FailNow()
}

func TestLocalAssignment(t *testing.T) {
	input := `local x = 5`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Body) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Body))
	}

	stmt, ok := program.Body[0].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("Statement not equal to LocalAssignment. Got=%T", program.Body[0])
	}

	if len(stmt.Names) != 1 || stmt.Names[0] != "x" {
		t.Errorf("Expected name of variable 'x', got %v", stmt.Names)
	}

	if len(stmt.Values) != 1 {
		t.Fatalf("Expected 1 value, got %d", len(stmt.Values))
	}

	literal, ok := stmt.Values[0].(*ast.Literal)
	if !ok {
		t.Fatalf("Value is not equal to Literal. Got=%T", stmt.Values[0])
	}

	if literal.Value != int64(5) {
		t.Errorf("Expected value 5, got %v", literal.Value)
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"-15", "-", int64(15)},
		{"not true", "not", true},
		{"#array", "#", "array"},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		expr := p.parseExpression(LOWEST)
		checkParserErrors(t, p)

		if expr == nil {
			t.Fatalf("parseExpression returned nil for %q", tt.input)
		}

		unaryExp, ok := expr.(*ast.UnaryOp)
		if !ok {
			t.Fatalf("Expected UnaryOp. got=%T", expr)
		}

		if unaryExp.Op != tt.operator {
			t.Errorf("Expected operator %q, got %q", tt.operator, unaryExp.Op)
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", int64(5), "+", int64(5)},
		{"5 - 5", int64(5), "-", int64(5)},
		{"5 * 5", int64(5), "*", int64(5)},
		{"5 / 5", int64(5), "/", int64(5)},
		{"5 > 5", int64(5), ">", int64(5)},
		{"5 < 5", int64(5), "<", int64(5)},
		{"5 == 5", int64(5), "==", int64(5)},
		{"5 ~= 5", int64(5), "~=", int64(5)},
		{"true == true", true, "==", true},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		expr := p.parseExpression(LOWEST)
		checkParserErrors(t, p)

		binaryExp, ok := expr.(*ast.BinaryOp)
		if !ok {
			t.Fatalf("Expected BinaryOp. Got=%T", expr)
		}

		if binaryExp.Op != tt.operator {
			t.Errorf("Expected operator %q, got %q", tt.operator, binaryExp.Op)
		}

		testLiteralObject(t, binaryExp.Left, tt.leftValue)
		testLiteralObject(t, binaryExp.Right, tt.rightValue)
	}
}

func testLiteralObject(t *testing.T, exp ast.Expr, expected any) {
	literal, ok := exp.(*ast.Literal)
	if !ok {
		t.Errorf("Expected ast.Literal, got %T", exp)
		return
	}

	if literal.Value != expected {
		t.Errorf("Expected value %v, got %v", expected, literal.Value)
	}
}

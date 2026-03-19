package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser found %d error(s)", len(errors))
	for _, msg := range errors {
		t.Errorf("  parser error: %v", msg)
	}
	t.FailNow()
}

func newParser(t *testing.T, factory *ast.NodeFactory, input string) (*Parser, *ast.Program) {
	t.Helper()
	factory.Reset()
	l := lexer.New(input)
	p := New(l, factory)
	program := p.ParseProgram()
	return p, program
}

func testLiteralObject(t *testing.T, exp ast.Expr, expected any) {
	t.Helper()
	literal, ok := exp.(*ast.Literal)
	if !ok {
		t.Errorf("expected ast.Literal, got %T", exp)
		return
	}
	if literal.Value != expected {
		t.Errorf("expected literal value %v, got %v", expected, literal.Value)
	}
}

func testIdentifier(t *testing.T, exp ast.Expr, value string) {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expected ast.Identifier, got %T", exp)
		return
	}
	if ident.Name != value {
		t.Errorf("expected identifier name %q, got %q", value, ident.Name)
	}
}

func formatType(node ast.TypeNode) string {
	if node == nil {
		return "nil"
	}
	switch n := node.(type) {
	case *ast.PrimitiveType:
		return n.Name
	case *ast.UnionType:
		return fmt.Sprintf("%s | %s", formatType(n.Left), formatType(n.Right))
	case *ast.OptionalType:
		return formatType(n.BaseType) + "?"
	case *ast.GenericType:
		types := make([]string, len(n.Types))
		for i, t := range n.Types {
			types[i] = formatType(t)
		}
		return fmt.Sprintf("%s<%s>", formatType(n.BaseType), strings.Join(types, ", "))
	case *ast.TableType:
		fields := make([]string, 0, len(n.Fields))
		for _, f := range n.Fields {
			switch {
			case f.IsAccess:
				fields = append(fields, fmt.Sprintf("[%s]: %s", formatType(f.Key), formatType(f.Value)))
			case f.KeyName != "":
				fields = append(fields, fmt.Sprintf("%s: %s", f.KeyName, formatType(f.Value)))
			default:
				fields = append(fields, formatType(f.Value))
			}
		}
		return "{" + strings.Join(fields, ", ") + "}"
	default:
		return "unknown_type"
	}
}

func formatAstTree(node ast.Node) string {
	if node == nil {
		return "nil"
	}
	switch n := node.(type) {
	case *ast.Identifier:
		return n.Name
	case *ast.Literal:
		if n.Type == "string" {
			return fmt.Sprintf("%q", n.Value)
		}
		return fmt.Sprintf("%v", n.Value)
	case *ast.BinaryOp:
		return fmt.Sprintf("BinaryOp: %s (Left: %s, Right: %s)", n.Op, formatAstTree(n.Left), formatAstTree(n.Right))
	case *ast.UnaryOp:
		return fmt.Sprintf("UnaryOp: %s (%s)", n.Op, formatAstTree(n.Operand))
	case *ast.ParenExpr:
		return fmt.Sprintf("(%s)", formatAstTree(n.Expr))
	default:
		return n.String()
	}
}

func TestLocalAssignment(t *testing.T) {
	factory := ast.NewFactory()
	p, program := newParser(t, factory, `local x = 5`)
	checkParserErrors(t, p)

	if len(program.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Body))
	}

	stmt, ok := program.Body[0].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[0])
	}
	if len(stmt.Names) != 1 || stmt.Names[0] != "x" {
		t.Errorf("expected variable name %q, got %v", "x", stmt.Names)
	}
	if len(stmt.Values) != 1 {
		t.Fatalf("expected 1 value, got %d", len(stmt.Values))
	}

	literal, ok := stmt.Values[0].(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", stmt.Values[0])
	}
	if literal.Value != int64(5) {
		t.Errorf("expected literal value 5, got %v", literal.Value)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		operator string
		value    any
	}{
		{"negative number", "-15", "-", int64(15)},
		{"not boolean", "not true", "not", true},
		{"length operator", "#array", "#", "array"},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory.Reset()
			l := lexer.New(tt.input)
			p := New(l, factory)

			expr := p.parseExpression(LOWEST)
			checkParserErrors(t, p)

			if expr == nil {
				t.Fatalf("parseExpression returned nil for %q", tt.input)
			}

			unaryExp, ok := expr.(*ast.UnaryOp)
			if !ok {
				t.Fatalf("expected *ast.UnaryOp, got %T", expr)
			}
			if unaryExp.Op != tt.operator {
				t.Errorf("expected operator %q, got %q", tt.operator, unaryExp.Op)
			}
		})
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"addition", "5 + 5", int64(5), "+", int64(5)},
		{"subtraction", "5 - 5", int64(5), "-", int64(5)},
		{"multiplication", "5 * 5", int64(5), "*", int64(5)},
		{"division", "5 / 5", int64(5), "/", int64(5)},
		{"greater than", "5 > 5", int64(5), ">", int64(5)},
		{"less than", "5 < 5", int64(5), "<", int64(5)},
		{"equality integers", "5 == 5", int64(5), "==", int64(5)},
		{"inequality", "5 ~= 5", int64(5), "~=", int64(5)},
		{"equality booleans", "true == true", true, "==", true},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory.Reset()
			l := lexer.New(tt.input)
			p := New(l, factory)

			expr := p.parseExpression(LOWEST)
			checkParserErrors(t, p)

			binaryExp, ok := expr.(*ast.BinaryOp)
			if !ok {
				t.Fatalf("expected *ast.BinaryOp, got %T", expr)
			}
			if binaryExp.Op != tt.operator {
				t.Errorf("expected operator %q, got %q", tt.operator, binaryExp.Op)
			}

			testLiteralObject(t, binaryExp.Left, tt.leftValue)
			testLiteralObject(t, binaryExp.Right, tt.rightValue)
		})
	}
}

func TestAssignments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		targets  []string
		operator string
		values   []int64
	}{
		{"simple assignment", "x = 5", []string{"x"}, "=", []int64{5}},
		{"compound addition", "y += 10", []string{"y"}, "+=", []int64{10}},
		{"multi-target", "a, b = 1, 2", []string{"a", "b"}, "=", []int64{1, 2}},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, program := newParser(t, factory, tt.input)
			checkParserErrors(t, p)

			if len(program.Body) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Body))
			}

			stmt, ok := program.Body[0].(*ast.Assignment)
			if !ok {
				t.Fatalf("expected *ast.Assignment, got %T", program.Body[0])
			}
			if stmt.Operator != tt.operator {
				t.Errorf("expected operator %q, got %q", tt.operator, stmt.Operator)
			}
			if len(stmt.Targets) != len(tt.targets) {
				t.Fatalf("expected %d target(s), got %d", len(tt.targets), len(stmt.Targets))
			}
			for i, name := range tt.targets {
				testIdentifier(t, stmt.Targets[i], name)
			}
			if len(stmt.Values) != len(tt.values) {
				t.Fatalf("expected %d value(s), got %d", len(tt.values), len(stmt.Values))
			}
			for i, val := range tt.values {
				testLiteralObject(t, stmt.Values[i], val)
			}
		})
	}
}

func TestFunctionAndMethodCalls(t *testing.T) {
	input := `
		print("hello")
		player:Damage(50)
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	if len(program.Body) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Body))
	}

	t.Run("function call", func(t *testing.T) {
		stmt, ok := program.Body[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got %T", program.Body[0])
		}
		call, ok := stmt.Expr.(*ast.FunctionCall)
		if !ok {
			t.Fatalf("expected *ast.FunctionCall, got %T", stmt.Expr)
		}
		testIdentifier(t, call.Function, "print")
		if len(call.Args) != 1 {
			t.Errorf("expected 1 argument, got %d", len(call.Args))
		}
	})

	t.Run("method call", func(t *testing.T) {
		stmt, ok := program.Body[1].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got %T", program.Body[1])
		}
		call, ok := stmt.Expr.(*ast.MethodCall)
		if !ok {
			t.Fatalf("expected *ast.MethodCall, got %T", stmt.Expr)
		}
		testIdentifier(t, call.Object, "player")
		if call.Method != "Damage" {
			t.Errorf("expected method %q, got %q", "Damage", call.Method)
		}
	})
}

func TestControlFlowStatements(t *testing.T) {
	input := `
		if true then end
		while x < 10 do end
		repeat break until true
		for i = 1, 10 do end
		for k, v in pairs do end
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	if len(program.Body) != 5 {
		t.Fatalf("expected 5 statements, got %d", len(program.Body))
	}

	t.Run("if statement", func(t *testing.T) {
		if _, ok := program.Body[0].(*ast.IfStatement); !ok {
			t.Errorf("expected *ast.IfStatement, got %T", program.Body[0])
		}
	})

	t.Run("while loop", func(t *testing.T) {
		if _, ok := program.Body[1].(*ast.WhileLoop); !ok {
			t.Errorf("expected *ast.WhileLoop, got %T", program.Body[1])
		}
	})

	t.Run("repeat loop", func(t *testing.T) {
		if _, ok := program.Body[2].(*ast.RepeatLoop); !ok {
			t.Errorf("expected *ast.RepeatLoop, got %T", program.Body[2])
		}
	})

	t.Run("numeric for loop", func(t *testing.T) {
		forLoop, ok := program.Body[3].(*ast.ForLoop)
		if !ok {
			t.Fatalf("expected *ast.ForLoop, got %T", program.Body[3])
		}
		if forLoop.Variable != "i" {
			t.Errorf("expected loop variable %q, got %q", "i", forLoop.Variable)
		}
	})

	t.Run("generic for loop", func(t *testing.T) {
		forIn, ok := program.Body[4].(*ast.ForInLoop)
		if !ok {
			t.Fatalf("expected *ast.ForInLoop, got %T", program.Body[4])
		}
		if len(forIn.Variables) != 2 || forIn.Variables[0] != "k" {
			t.Errorf("expected for-in variables [k, v], got %v", forIn.Variables)
		}
	})
}

func TestFunctionDefinitions(t *testing.T) {
	input := `
		function add(a: number, b: number): number
			return a + b
		end
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	fn, ok := program.Body[0].(*ast.FunctionDef)
	if !ok {
		t.Fatalf("expected *ast.FunctionDef, got %T", program.Body[0])
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(fn.Parameters))
	}
	if fn.Parameters[0].Name != "a" || formatType(fn.Parameters[0].Type) != "number" {
		t.Errorf("expected param %q, got %q: %s", "a: number", fn.Parameters[0].Name, formatType(fn.Parameters[0].Type))
	}
	if fn.ReturnType == nil || formatType(fn.ReturnType) != "number" {
		t.Errorf("expected return type %q, got %q", "number", formatType(fn.ReturnType))
	}
}

func TestTablesAndAccess(t *testing.T) {
	input := `
		local myTable = {
			1,
			key = "value",
			[5] = true
		}
		local a = myTable.key
		local b = myTable[5]
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	t.Run("table literal fields", func(t *testing.T) {
		localAssign, ok := program.Body[0].(*ast.LocalAssignment)
		if !ok {
			t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[0])
		}
		tableLit, ok := localAssign.Values[0].(*ast.TableLiteral)
		if !ok {
			t.Fatalf("expected *ast.TableLiteral, got %T", localAssign.Values[0])
		}
		if len(tableLit.Fields) != 3 {
			t.Fatalf("expected 3 table fields, got %d", len(tableLit.Fields))
		}
		if tableLit.Fields[0].Key != nil {
			t.Errorf("expected array field key to be nil")
		}
		if tableLit.Fields[1].Key == nil {
			t.Errorf("expected dictionary field key to not be nil")
		}
	})

	t.Run("field access", func(t *testing.T) {
		stmt, ok := program.Body[1].(*ast.LocalAssignment)
		if !ok {
			t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[1])
		}
		fieldAcc, ok := stmt.Values[0].(*ast.FieldAccess)
		if !ok {
			t.Fatalf("expected *ast.FieldAccess, got %T", stmt.Values[0])
		}
		if fieldAcc.Field != "key" {
			t.Errorf("expected field %q, got %q", "key", fieldAcc.Field)
		}
	})

	t.Run("index access", func(t *testing.T) {
		stmt, ok := program.Body[2].(*ast.LocalAssignment)
		if !ok {
			t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[2])
		}
		idxAcc, ok := stmt.Values[0].(*ast.IndexAccess)
		if !ok {
			t.Fatalf("expected *ast.IndexAccess, got %T", stmt.Values[0])
		}
		testLiteralObject(t, idxAcc.Index, int64(5))
	})
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"additive vs multiplicative",
			"a + b * c",
			"BinaryOp: + (Left: a, Right: BinaryOp: * (Left: b, Right: c))",
		},
		{
			"additive vs equality",
			"a + b == c",
			"BinaryOp: == (Left: BinaryOp: + (Left: a, Right: b), Right: c)",
		},
		{
			"left-associative addition",
			"a + b + c",
			"BinaryOp: + (Left: BinaryOp: + (Left: a, Right: b), Right: c)",
		},
		{
			"left-associative mul and div",
			"a * b / c",
			"BinaryOp: / (Left: BinaryOp: * (Left: a, Right: b), Right: c)",
		},
		{
			"unary minus vs multiplication",
			"-a * b",
			"BinaryOp: * (Left: UnaryOp: - (a), Right: b)",
		},
		{
			"not vs equality",
			"not a == b",
			"BinaryOp: == (Left: UnaryOp: not (a), Right: b)",
		},
		{
			"right-associative exponentiation",
			"a ^ b ^ c",
			"BinaryOp: ^ (Left: a, Right: BinaryOp: ^ (Left: b, Right: c))",
		},
		{
			"right-associative concatenation",
			"a .. b .. c",
			"BinaryOp: .. (Left: a, Right: BinaryOp: .. (Left: b, Right: c))",
		},
		{
			"parentheses override precedence",
			"a + (b * c) + d",
			"BinaryOp: + (Left: BinaryOp: + (Left: a, Right: (BinaryOp: * (Left: b, Right: c))), Right: d)",
		},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory.Reset()
			l := lexer.New(tt.input)
			p := New(l, factory)

			expr := p.parseExpression(LOWEST)
			checkParserErrors(t, p)

			if actual := formatAstTree(expr); actual != tt.expected {
				t.Errorf("precedence mismatch for %q\n  expected: %s\n  got:      %s", tt.input, tt.expected, actual)
			}
		})
	}
}

func TestLuauSpecificExpressions(t *testing.T) {
	input := `
		local a = x :: string
		local b = if true then 1 else 2
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	t.Run("type cast", func(t *testing.T) {
		stmt, ok := program.Body[0].(*ast.LocalAssignment)
		if !ok {
			t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[0])
		}
		typeCast, ok := stmt.Values[0].(*ast.TypeCast)
		if !ok {
			t.Fatalf("expected *ast.TypeCast, got %T", stmt.Values[0])
		}
		if formatType(typeCast.Type) != "string" {
			t.Errorf("expected cast type %q, got %q", "string", formatType(typeCast.Type))
		}
	})

	t.Run("if expression", func(t *testing.T) {
		stmt, ok := program.Body[1].(*ast.LocalAssignment)
		if !ok {
			t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[1])
		}
		ifExpr, ok := stmt.Values[0].(*ast.IfExpr)
		if !ok {
			t.Fatalf("expected *ast.IfExpr, got %T", stmt.Values[0])
		}
		testLiteralObject(t, ifExpr.Condition, true)
	})
}

func TestAnonymousFunctionsAndVarargs(t *testing.T) {
	input := `
		local cb = function(...)
			print(...)
		end
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	assign, ok := program.Body[0].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[0])
	}

	fn, ok := assign.Values[0].(*ast.FunctionExpr)
	if !ok {
		t.Fatalf("expected *ast.FunctionExpr, got %T", assign.Values[0])
	}
	if len(fn.Parameters) != 1 || fn.Parameters[0].Name != "..." {
		t.Errorf("expected vararg parameter %q, got %v", "...", fn.Parameters)
	}

	exprStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement in function body, got %T", fn.Body.Statements[0])
	}
	call, ok := exprStmt.Expr.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("expected *ast.FunctionCall, got %T", exprStmt.Expr)
	}
	if _, ok := call.Args[0].(*ast.VarArgs); !ok {
		t.Errorf("expected *ast.VarArgs as argument, got %T", call.Args[0])
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{"local with non-identifier name", "local 5 = x", "expected Identifier"},
		{"function without name", "function() end", "expected IDENT for function name"},
		{"unclosed if block", "if true then", "expected END to close if statement"},
		{"unclosed parenthesis", "x = (5 + 5", "expected closing bracket"},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory.Reset()
			l := lexer.New(tt.input)
			p := New(l, factory)
			p.ParseProgram()

			errors := p.Errors()
			if len(errors) == 0 {
				t.Fatalf("expected a parse error containing %q, got none", tt.expectedErr)
			}

			for _, err := range errors {
				if strings.Contains(err.Error(), tt.expectedErr) {
					return
				}
			}
			t.Errorf("expected error containing %q, got: %v", tt.expectedErr, errors)
		})
	}
}

func TestStandardLuaFeatures(t *testing.T) {
	input := `
		do
			local x = 1
		end

		if a then
			print(1)
		elseif b then
			print(2)
		else
			print(3)
		end

		local function helper() end

		print "hello"
		print {1, 2, 3}
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	if len(program.Body) != 5 {
		t.Fatalf("expected 5 statements, got %d", len(program.Body))
	}

	t.Run("do block", func(t *testing.T) {
		if _, ok := program.Body[0].(*ast.DoBlock); !ok {
			t.Errorf("expected *ast.DoBlock, got %T", program.Body[0])
		}
	})

	t.Run("if/elseif/else", func(t *testing.T) {
		ifStmt, ok := program.Body[1].(*ast.IfStatement)
		if !ok {
			t.Fatalf("expected *ast.IfStatement, got %T", program.Body[1])
		}
		if len(ifStmt.ElseIfs) != 1 {
			t.Errorf("expected 1 elseif block, got %d", len(ifStmt.ElseIfs))
		}
		if ifStmt.Else == nil {
			t.Errorf("expected else block, got nil")
		}
	})

	t.Run("local function", func(t *testing.T) {
		localFn, ok := program.Body[2].(*ast.LocalFunction)
		if !ok {
			t.Fatalf("expected *ast.LocalFunction, got %T", program.Body[2])
		}
		if localFn.Name != "helper" {
			t.Errorf("expected local function name %q, got %q", "helper", localFn.Name)
		}
	})

	t.Run("call with string arg", func(t *testing.T) {
		stmt, ok := program.Body[3].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got %T", program.Body[3])
		}
		call, ok := stmt.Expr.(*ast.FunctionCall)
		if !ok || len(call.Args) != 1 {
			t.Fatalf("expected *ast.FunctionCall with 1 arg, got %T", stmt.Expr)
		}
		testLiteralObject(t, call.Args[0], "hello")
	})

	t.Run("call with table arg", func(t *testing.T) {
		stmt, ok := program.Body[4].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got %T", program.Body[4])
		}
		call, ok := stmt.Expr.(*ast.FunctionCall)
		if !ok {
			t.Fatalf("expected *ast.FunctionCall, got %T", stmt.Expr)
		}
		if _, ok := call.Args[0].(*ast.TableLiteral); !ok {
			t.Errorf("expected *ast.TableLiteral as argument, got %T", call.Args[0])
		}
	})
}

func TestLuauTypeAliasesAndContinue(t *testing.T) {
	input := `
		type Point = { x: number, y: number }
		export type ID = string | number

		for i = 1, 10 do
			continue
		end
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	if len(program.Body) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Body))
	}

	t.Run("non-exported type alias", func(t *testing.T) {
		typeAlias, ok := program.Body[0].(*ast.TypeAlias)
		if !ok {
			t.Fatalf("expected *ast.TypeAlias, got %T", program.Body[0])
		}
		if typeAlias.Name != "Point" {
			t.Errorf("expected type alias name %q, got %q", "Point", typeAlias.Name)
		}
		if typeAlias.IsExport {
			t.Errorf("expected non-exported type alias")
		}
	})

	t.Run("exported type alias with union", func(t *testing.T) {
		typeAlias, ok := program.Body[1].(*ast.TypeAlias)
		if !ok {
			t.Fatalf("expected *ast.TypeAlias, got %T", program.Body[1])
		}
		if typeAlias.Name != "ID" {
			t.Errorf("expected type alias name %q, got %q", "ID", typeAlias.Name)
		}
		if !typeAlias.IsExport {
			t.Errorf("expected exported type alias")
		}
		if typeAlias.Type == nil || formatType(typeAlias.Type) != "string | number" {
			t.Errorf("expected type %q, got %q", "string | number", formatType(typeAlias.Type))
		}
	})

	t.Run("continue inside for loop", func(t *testing.T) {
		forLoop, ok := program.Body[2].(*ast.ForLoop)
		if !ok {
			t.Fatalf("expected *ast.ForLoop, got %T", program.Body[2])
		}
		if len(forLoop.Body.Statements) != 1 {
			t.Fatalf("expected 1 statement in loop body, got %d", len(forLoop.Body.Statements))
		}
		if _, ok := forLoop.Body.Statements[0].(*ast.ContinueStatement); !ok {
			t.Errorf("expected *ast.ContinueStatement, got %T", forLoop.Body.Statements[0])
		}
	})
}

func TestLuauGenerics(t *testing.T) {
	input := `
		type Map<K, V> = { [K]: V }

		function reverse<T>(arr: {T}): {T}
			return arr
		end

		local function id<Value>(v: Value): Value
			return v
		end
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	if len(program.Body) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Body))
	}

	t.Run("generic type alias", func(t *testing.T) {
		typeAlias, ok := program.Body[0].(*ast.TypeAlias)
		if !ok {
			t.Fatalf("expected *ast.TypeAlias, got %T", program.Body[0])
		}
		if len(typeAlias.Generics) != 2 || typeAlias.Generics[0] != "K" || typeAlias.Generics[1] != "V" {
			t.Errorf("expected generics [K, V], got %v", typeAlias.Generics)
		}
	})

	t.Run("generic function def", func(t *testing.T) {
		funcDef, ok := program.Body[1].(*ast.FunctionDef)
		if !ok {
			t.Fatalf("expected *ast.FunctionDef, got %T", program.Body[1])
		}
		if len(funcDef.Generics) != 1 || funcDef.Generics[0] != "T" {
			t.Errorf("expected generics [T], got %v", funcDef.Generics)
		}
	})

	t.Run("generic local function", func(t *testing.T) {
		localFunc, ok := program.Body[2].(*ast.LocalFunction)
		if !ok {
			t.Fatalf("expected *ast.LocalFunction, got %T", program.Body[2])
		}
		if len(localFunc.Generics) != 1 || localFunc.Generics[0] != "Value" {
			t.Errorf("expected generics [Value], got %v", localFunc.Generics)
		}
	})
}

func TestInterpolatedStrings(t *testing.T) {
	input := "local greeting = `Hello {\"World\"}, I am {10 + 10} years old!`"

	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	assign, ok := program.Body[0].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[0])
	}
	interp, ok := assign.Values[0].(*ast.InterpolatedString)
	if !ok {
		t.Fatalf("expected *ast.InterpolatedString, got %T", assign.Values[0])
	}

	if len(interp.Segments) != 3 {
		t.Fatalf("expected 3 string segments, got %d", len(interp.Segments))
	}
	if interp.Segments[0] != "Hello " || interp.Segments[1] != ", I am " || interp.Segments[2] != " years old!" {
		t.Errorf("unexpected segments: %v", interp.Segments)
	}
	if len(interp.Expressions) != 2 {
		t.Errorf("expected 2 interpolated expressions, got %d", len(interp.Expressions))
	}
}

func TestFunctionAttributes(t *testing.T) {
	input := `
		@native
		function mathFast()
		end

		@checked @native
		local function secureMath()
		end
	`
	factory := ast.NewFactory()
	p, program := newParser(t, factory, input)
	checkParserErrors(t, p)

	if len(program.Body) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Body))
	}

	t.Run("single attribute on function def", func(t *testing.T) {
		fn, ok := program.Body[0].(*ast.FunctionDef)
		if !ok {
			t.Fatalf("expected *ast.FunctionDef, got %T", program.Body[0])
		}
		if len(fn.Attributes) != 1 {
			t.Fatalf("expected 1 attribute, got %d", len(fn.Attributes))
		}
		if fn.Attributes[0].Name != "native" {
			t.Errorf("expected attribute %q, got %q", "native", fn.Attributes[0].Name)
		}
	})

	t.Run("multiple attributes on local function", func(t *testing.T) {
		localFn, ok := program.Body[1].(*ast.LocalFunction)
		if !ok {
			t.Fatalf("expected *ast.LocalFunction, got %T", program.Body[1])
		}
		if len(localFn.Attributes) != 2 {
			t.Fatalf("expected 2 attributes, got %d", len(localFn.Attributes))
		}
		if localFn.Attributes[0].Name != "checked" {
			t.Errorf("expected first attribute %q, got %q", "checked", localFn.Attributes[0].Name)
		}
		if localFn.Attributes[1].Name != "native" {
			t.Errorf("expected second attribute %q, got %q", "native", localFn.Attributes[1].Name)
		}
	})
}

func TestAttributeErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{"attribute on local variable", "@native local x = 5", "attributes are only allowed on function declarations"},
		{"attribute on type alias", "@native type Point = {x: number, y: number}", "attributes are only allowed on function declarations"},
		{"bare at-sign", "@", "expected identifier after '@'"},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory.Reset()
			l := lexer.New(tt.input)
			p := New(l, factory)
			p.ParseProgram()

			errors := p.Errors()
			if len(errors) == 0 {
				t.Fatalf("expected a parse error containing %q, got none", tt.expectedErr)
			}

			for _, err := range errors {
				if strings.Contains(err.Error(), tt.expectedErr) {
					return
				}
			}
			t.Errorf("expected error containing %q, got: %v", tt.expectedErr, errors)
		})
	}
}

func TestSoftKeywordsAsIdentifiers(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"type as local variable", `local type = "hello"`},
		{"type as assignment target", `type = "hello"`},
		{"type as field access object", `type.Name = "Foo"`},
		{"type as function argument", `print(type)`},
		{"export as local variable", `local export = true`},
		{"export as assignment target", `export = true`},
		{"export as function argument", `print(export)`},
		{"type in binary expression", `local x = type + 1`},
		{"type as table field value", `local t = { type = "foo" }`},
		{"type as method call object", `type:GetName()`},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, program := newParser(t, factory, tt.input)
			checkParserErrors(t, p)

			if len(program.Body) == 0 {
				t.Fatal("expected at least one statement")
			}
		})
	}
}

func TestAssignmentTargetsAreNeverNil(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple assignment", `x = 5`},
		{"two targets", `x, y = 1, 2`},
		{"three targets", `x, y, z = 1, 2, 3`},
		{"field assignment", `a.b = 5`},
		{"index assignment", `a[1] = 5`},
		{"multiple field assignments", `a.b, c.d = 1, 2`},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, program := newParser(t, factory, tt.input)
			checkParserErrors(t, p)

			for _, stmt := range program.Body {
				assign, ok := stmt.(*ast.Assignment)
				if !ok {
					continue
				}
				for i, target := range assign.Targets {
					if target == nil {
						t.Errorf("Assignment.Targets[%d] is nil — would panic in visitor", i)
					}
				}
				for i, value := range assign.Values {
					if value == nil {
						t.Errorf("Assignment.Values[%d] is nil — would panic in visitor", i)
					}
				}
			}
		})
	}
}

func TestFunctionCallArgsAreNeverNil(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single arg", `print(1)`},
		{"multiple args", `print(1, 2, 3)`},
		{"expression arg", `print(a + b)`},
		{"nested call arg", `print(fn())`},
		{"mixed type args", `f(1, "two", true, nil)`},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, program := newParser(t, factory, tt.input)
			checkParserErrors(t, p)

			for _, stmt := range program.Body {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					continue
				}
				call, ok := exprStmt.Expr.(*ast.FunctionCall)
				if !ok {
					continue
				}
				for i, arg := range call.Args {
					if arg == nil {
						t.Errorf("FunctionCall.Args[%d] is nil — would panic in visitor", i)
					}
				}
			}
		})
	}
}

func TestReturnStatementValuesAreNeverNil(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty return", `return`},
		{"single value", `return 1`},
		{"multiple values", `return 1, 2`},
		{"expression value", `return a + b`},
		{"mixed type values", `return true, false, nil`},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, program := newParser(t, factory, tt.input)

			for _, stmt := range program.Body {
				ret, ok := stmt.(*ast.ReturnStatement)
				if !ok {
					continue
				}
				for i, val := range ret.Values {
					if val == nil {
						t.Errorf("ReturnStatement.Values[%d] is nil — would panic in visitor", i)
					}
				}
			}
		})
	}
}

func TestParserDoesNotPanicOnMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"bare assignment operator", `= 5`},
		{"double assignment operator", `x == 5`},
		{"unclosed table", `local t = {1, 2`},
		{"unclosed function call", `print(1, 2`},
		{"empty function call args with comma", `print(,)`},
		{"missing then in if", `if true end`},
		{"missing condition in while", `while do end`},
		{"missing in in for-in", `for k, v pairs do end`},
		{"missing do in for", `for i = 1, 10 end`},
		{"local without name", `local = 5`},
		{"function without name", `function() end`},
		{"chained assignment operators", `x = = 5`},
		{"unmatched bracket in expression", `local x = (1 + 2`},
		{"type cast without type", `local x = val ::`},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("parser panicked on input %q: %v", tt.input, r)
				}
			}()

			factory.Reset()
			l := lexer.New(tt.input)
			p := New(l, factory)
			p.ParseProgram()

			if len(p.Errors()) == 0 {
				t.Logf("no errors recorded for malformed input %q — may be worth checking", tt.input)
			}
		})
	}
}

func TestTreeIsWalkableAfterParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"local assignment", `local x = 5`},
		{"assignment", `x = 5`},
		{"function call expression", `print(x + 12)`},
		{"if statement", `if true then x = 1 end`},
		{"while loop", `while x < 10 do x += 1 end`},
		{"numeric for loop", `for i = 1, 10 do print(i) end`},
		{"generic for loop", `for k, v in pairs(t) do print(k, v) end`},
		{"repeat loop", `repeat x += 1 until x >= 10`},
		{"type as identifier", `local type = "hello"`},
		{"type assignment", `type = "world"`},
		{"typed function def", `function foo(a: number, b: number): number return a + b end`},
		{"generic local function", `local function bar<T>(v: T): T return v end`},
		{"table literal", `local t = { 1, key = "val", [5] = true }`},
		{"type cast", `local x = a :: string`},
		{"if expression", `local b = if true then 1 else 2`},
		{"interpolated string", "local s = `hello {name}!`"},
		{"do block", `do local x = 1 end`},
		{"exported type alias", `export type ID = string | number`},
		{"type alias", `type Point = { x: number, y: number }`},
		{"attributed function", `@native function fast() end`},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("tree walk panicked for input %q: %v", tt.input, r)
				}
			}()

			p, program := newParser(t, factory, tt.input)
			checkParserErrors(t, p)
			walkProgram(t, program)
		})
	}
}

func walkProgram(t *testing.T, program *ast.Program) {
	t.Helper()
	for i, stmt := range program.Body {
		if stmt == nil {
			t.Errorf("Program.Body[%d] is nil", i)
			return
		}
		walkStmt(t, stmt)
	}
}

func walkStmt(t *testing.T, stmt ast.Stmt) {
	t.Helper()
	switch s := stmt.(type) {
	case *ast.LocalAssignment:
		for i, v := range s.Values {
			if v == nil {
				t.Errorf("LocalAssignment.Values[%d] is nil", i)
			}
		}
	case *ast.Assignment:
		for i, target := range s.Targets {
			if target == nil {
				t.Errorf("Assignment.Targets[%d] is nil", i)
			}
		}
		for i, v := range s.Values {
			if v == nil {
				t.Errorf("Assignment.Values[%d] is nil", i)
			}
		}
	case *ast.ExpressionStatement:
		if s.Expr == nil {
			t.Error("ExpressionStatement.Expr is nil")
		}
	case *ast.ReturnStatement:
		for i, v := range s.Values {
			if v == nil {
				t.Errorf("ReturnStatement.Values[%d] is nil", i)
			}
		}
	case *ast.IfStatement:
		if s.Condition == nil {
			t.Error("IfStatement.Condition is nil")
		}
		if s.Then == nil {
			t.Error("IfStatement.Then is nil")
			return
		}
		for _, stmt := range s.Then.Statements {
			walkStmt(t, stmt)
		}
	case *ast.WhileLoop:
		if s.Condition == nil {
			t.Error("WhileLoop.Condition is nil")
		}
	case *ast.ForLoop:
		if s.Start == nil {
			t.Error("ForLoop.Start is nil")
		}
		if s.End == nil {
			t.Error("ForLoop.End is nil")
		}
	case *ast.ForInLoop:
		for i, iter := range s.Iterables {
			if iter == nil {
				t.Errorf("ForInLoop.Iterables[%d] is nil", i)
			}
		}
	case *ast.FunctionDef:
		if s.Body == nil {
			t.Error("FunctionDef.Body is nil")
		}
	case *ast.LocalFunction:
		if s.Body == nil {
			t.Error("LocalFunction.Body is nil")
		}
	}
}

func TestNestedTableLiteralAsLastField(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		fieldCount int
	}{
		{
			name:       "nested table as last field",
			input:      `local t = { 1, {2, 3} }`,
			fieldCount: 2,
		},
		{
			name:       "nested table as non-last field",
			input:      `local t = { {1, 2}, 3 }`,
			fieldCount: 2,
		},
		{
			name:       "all nested tables",
			input:      `local t = { {1}, {2}, {3} }`,
			fieldCount: 3,
		},
		{
			name:       "deeply nested tables",
			input:      `local t = { 1, { 2, {3, 4} } }`,
			fieldCount: 2,
		},
		{
			name:       "nested table then statement",
			input:      "local t = { 1, {2, 3} }\nprint(t)",
			fieldCount: 2,
		},
	}

	factory := ast.NewFactory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, program := newParser(t, factory, tt.input)
			checkParserErrors(t, p)

			if len(program.Body) == 0 {
				t.Fatal("expected at least one statement, got none")
			}

			assign, ok := program.Body[0].(*ast.LocalAssignment)
			if !ok {
				t.Fatalf("expected *ast.LocalAssignment, got %T", program.Body[0])
			}
			if len(assign.Values) == 0 {
				t.Fatal("expected at least one value in assignment")
			}

			table, ok := assign.Values[0].(*ast.TableLiteral)
			if !ok {
				t.Fatalf("expected *ast.TableLiteral as assigned value, got %T", assign.Values[0])
			}
			if len(table.Fields) != tt.fieldCount {
				t.Errorf("expected %d fields, got %d", tt.fieldCount, len(table.Fields))
			}

			for i, field := range table.Fields {
				if field == nil {
					t.Errorf("Fields[%d] is nil", i)
				}
			}
		})
	}
}

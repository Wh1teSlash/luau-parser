package parser

import (
	"fmt"
	"strings"
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

func testIdentifier(t *testing.T, exp ast.Expr, value string) {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected ast.Identifier, got %T", exp)
		return
	}
	if ident.Name != value {
		t.Errorf("Expected identifier name %s, got %s", value, ident.Name)
	}
}

func TestAssignments(t *testing.T) {
	tests := []struct {
		input    string
		targets  []string
		operator string
		values   []int64
	}{
		{"x = 5", []string{"x"}, "=", []int64{5}},
		{"y += 10", []string{"y"}, "+=", []int64{10}},
		{"a, b = 1, 2", []string{"a", "b"}, "=", []int64{1, 2}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Body) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(program.Body))
		}

		stmt, ok := program.Body[0].(*ast.Assignment)
		if !ok {
			t.Fatalf("Statement not equal to Assignment. Got=%T", program.Body[0])
		}

		if stmt.Operator != tt.operator {
			t.Errorf("Expected operator %q, got %q", tt.operator, stmt.Operator)
		}

		if len(stmt.Targets) != len(tt.targets) {
			t.Fatalf("Expected %d targets, got %d", len(tt.targets), len(stmt.Targets))
		}
		for i, name := range tt.targets {
			testIdentifier(t, stmt.Targets[i], name)
		}

		if len(stmt.Values) != len(tt.values) {
			t.Fatalf("Expected %d values, got %d", len(tt.values), len(stmt.Values))
		}
		for i, val := range tt.values {
			testLiteralObject(t, stmt.Values[i], val)
		}
	}
}

func TestFunctionAndMethodCalls(t *testing.T) {
	input := `
		print("hello")
		player:Damage(50)
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Body) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(program.Body))
	}

	stmt1, ok := program.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got=%T", program.Body[0])
	}
	call1, ok := stmt1.Expr.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("Expected FunctionCall, got=%T", stmt1.Expr)
	}
	testIdentifier(t, call1.Function, "print")
	if len(call1.Args) != 1 {
		t.Fatalf("Expected 1 arg, got %d", len(call1.Args))
	}

	stmt2, ok := program.Body[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got=%T", program.Body[1])
	}
	call2, ok := stmt2.Expr.(*ast.MethodCall)
	if !ok {
		t.Fatalf("Expected MethodCall, got=%T", stmt2.Expr)
	}
	testIdentifier(t, call2.Object, "player")
	if call2.Method != "Damage" {
		t.Errorf("Expected method 'Damage', got %s", call2.Method)
	}
}

func TestControlFlowStatements(t *testing.T) {
	input := `
		if true then end
		while x < 10 do end
		repeat break until true
		for i = 1, 10 do end
		for k, v in pairs do end
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Body) != 5 {
		t.Fatalf("Expected 5 statements, got %d", len(program.Body))
	}

	_, ok := program.Body[0].(*ast.IfStatement)
	if !ok {
		t.Errorf("Expected IfStatement, got=%T", program.Body[0])
	}

	_, ok = program.Body[1].(*ast.WhileLoop)
	if !ok {
		t.Errorf("Expected WhileLoop, got=%T", program.Body[1])
	}

	_, ok = program.Body[2].(*ast.RepeatLoop)
	if !ok {
		t.Errorf("Expected RepeatLoop, got=%T", program.Body[2])
	}

	forLoop, ok := program.Body[3].(*ast.ForLoop)
	if !ok {
		t.Errorf("Expected ForLoop, got=%T", program.Body[3])
	}
	if forLoop.Variable != "i" {
		t.Errorf("Expected ForLoop variable 'i', got %s", forLoop.Variable)
	}

	forIn, ok := program.Body[4].(*ast.ForInLoop)
	if !ok {
		t.Errorf("Expected ForInLoop, got=%T", program.Body[4])
	}
	if len(forIn.Variables) != 2 || forIn.Variables[0] != "k" {
		t.Errorf("Expected ForInLoop vars [k, v], got %v", forIn.Variables)
	}
}

func TestFunctionDefinitions(t *testing.T) {
	input := `
		function add(a: number, b: number): number
			return a + b
		end
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Body) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Body))
	}

	fn, ok := program.Body[0].(*ast.FunctionDef)
	if !ok {
		t.Fatalf("Expected FunctionDef, got=%T", program.Body[0])
	}

	if fn.Name != "add" {
		t.Errorf("Expected function name 'add', got %s", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(fn.Parameters))
	}
	if fn.Parameters[0].Name != "a" || fn.Parameters[0].Type.Type != "number" {
		t.Errorf("Expected param 'a: number'")
	}

	if fn.ReturnType == nil || fn.ReturnType.Type != "number" {
		t.Errorf("Expected return type 'number'")
	}

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in body, got %d", len(fn.Body.Statements))
	}

	_, ok = fn.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Errorf("Expected ReturnStatement, got=%T", fn.Body.Statements[0])
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
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	localAssign, ok := program.Body[0].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("Expected LocalAssignment, got=%T", program.Body[0])
	}

	tableLit, ok := localAssign.Values[0].(*ast.TableLiteral)
	if !ok {
		t.Fatalf("Expected TableLiteral, got=%T", localAssign.Values[0])
	}

	if len(tableLit.Fields) != 3 {
		t.Fatalf("Expected 3 table fields, got %d", len(tableLit.Fields))
	}
	if tableLit.Fields[0].Key != nil {
		t.Errorf("Expected array field key to be nil")
	}
	if tableLit.Fields[1].Key == nil {
		t.Errorf("Expected dictionary field key to not be nil")
	}

	stmt2, ok := program.Body[1].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("Expected LocalAssignment, got=%T", program.Body[1])
	}

	fieldAcc, ok := stmt2.Values[0].(*ast.FieldAccess)
	if !ok {
		t.Fatalf("Expected FieldAccess, got=%T", stmt2.Values[0])
	}
	if fieldAcc.Field != "key" {
		t.Errorf("Expected field 'key', got %s", fieldAcc.Field)
	}

	stmt3, ok := program.Body[2].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("Expected LocalAssignment, got=%T", program.Body[2])
	}

	idxAcc, ok := stmt3.Values[0].(*ast.IndexAccess)
	if !ok {
		t.Fatalf("Expected IndexAccess, got=%T", stmt3.Values[0])
	}
	testLiteralObject(t, idxAcc.Index, int64(5))
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a + b * c",
			"BinaryOp: + (Left: a, Right: BinaryOp: * (Left: b, Right: c))",
		},
		{
			"a + b == c",
			"BinaryOp: == (Left: BinaryOp: + (Left: a, Right: b), Right: c)",
		},
		{
			"a + b + c",
			"BinaryOp: + (Left: BinaryOp: + (Left: a, Right: b), Right: c)",
		},
		{
			"a * b / c",
			"BinaryOp: / (Left: BinaryOp: * (Left: a, Right: b), Right: c)",
		},
		{
			"-a * b",
			"BinaryOp: * (Left: UnaryOp: - (a), Right: b)",
		},
		{
			"not a == b",
			"BinaryOp: == (Left: UnaryOp: not (a), Right: b)",
		},
		{
			"a ^ b ^ c",
			"BinaryOp: ^ (Left: a, Right: BinaryOp: ^ (Left: b, Right: c))",
		},
		{
			"a .. b .. c",
			"BinaryOp: .. (Left: a, Right: BinaryOp: .. (Left: b, Right: c))",
		},
		{
			"a + (b * c) + d",
			"BinaryOp: + (Left: BinaryOp: + (Left: a, Right: (BinaryOp: * (Left: b, Right: c))), Right: d)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		expr := p.parseExpression(LOWEST)
		checkParserErrors(t, p)

		actual := formatAstTree(expr)
		if actual != tt.expected {
			t.Errorf("Precedence mismatch for %q\nExpected: %s\nGot:      %s", tt.input, tt.expected, actual)
		}
	}
}

func TestLuauSpecificExpressions(t *testing.T) {
	input := `
		local a = x :: string
		local b = if true then 1 else 2
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt1, ok := program.Body[0].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("Expected LocalAssignment, got=%T", program.Body[0])
	}

	typeCast, ok := stmt1.Values[0].(*ast.TypeCast)
	if !ok {
		t.Fatalf("Expected TypeCast, got=%T", stmt1.Values[0])
	}
	if typeCast.Type.Type != "string" {
		t.Errorf("Expected cast to string, got %s", typeCast.Type.Type)
	}

	stmt2, ok := program.Body[1].(*ast.LocalAssignment)
	if !ok {
		t.Fatalf("Expected LocalAssignment, got=%T", program.Body[1])
	}

	ifExpr, ok := stmt2.Values[0].(*ast.IfExpr)
	if !ok {
		t.Fatalf("Expected IfExpr, got=%T", stmt2.Values[0])
	}
	testLiteralObject(t, ifExpr.Condition, true)
	testLiteralObject(t, ifExpr.Then, int64(1))
	testLiteralObject(t, ifExpr.Else, int64(2))
}

func TestAnonymousFunctionsAndVarargs(t *testing.T) {
	input := `
		local cb = function(...)
			print(...)
		end
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assign := program.Body[0].(*ast.LocalAssignment)
	fn, ok := assign.Values[0].(*ast.FunctionExpr)
	if !ok {
		t.Fatalf("Expected FunctionExpr, got=%T", assign.Values[0])
	}

	if len(fn.Parameters) != 1 || fn.Parameters[0].Name != "..." {
		t.Errorf("Expected vararg parameter '...', got %v", fn.Parameters)
	}

	exprStmt := fn.Body.Statements[0].(*ast.ExpressionStatement)
	call := exprStmt.Expr.(*ast.FunctionCall)

	_, ok = call.Args[0].(*ast.VarArgs)
	if !ok {
		t.Errorf("Expected VarArgs as argument, got=%T", call.Args[0])
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedErrFrag string
	}{
		{"local 5 = x", "expected Identifier"},
		{"function() end", "expected IDENT for function name"},
		{"if true then", "expected END to close if statement"},
		{"x = (5 + 5", "expected closing bracket"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		p.ParseProgram()

		errors := p.Errors()
		if len(errors) == 0 {
			t.Errorf("Expected parsing error for input %q, but got none", tt.input)
			continue
		}

		found := false
		for _, err := range errors {
			if strings.Contains(err.Error(), tt.expectedErrFrag) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected error containing %q, got: %v", tt.expectedErrFrag, errors)
		}
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
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Body) != 5 {
		t.Fatalf("Expected 5 statements, got %d", len(program.Body))
	}

	_, ok := program.Body[0].(*ast.DoBlock)
	if !ok {
		t.Errorf("Expected DoBlock, got=%T", program.Body[0])
	}

	ifStmt, ok := program.Body[1].(*ast.IfStatement)
	if !ok {
		t.Fatalf("Expected IfStatement, got=%T", program.Body[1])
	}
	if len(ifStmt.ElseIfs) != 1 {
		t.Errorf("Expected 1 ElseIf block, got %d", len(ifStmt.ElseIfs))
	}
	if ifStmt.Else == nil {
		t.Errorf("Expected Else block to not be nil")
	}

	localFn, ok := program.Body[2].(*ast.LocalFunction)
	if !ok {
		t.Errorf("Expected LocalFunction, got=%T", program.Body[2])
	}
	if localFn.Name != "helper" {
		t.Errorf("Expected local function 'helper', got %s", localFn.Name)
	}

	call1Stmt, ok := program.Body[3].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got=%T", program.Body[3])
	}
	call1, ok := call1Stmt.Expr.(*ast.FunctionCall)
	if !ok || len(call1.Args) != 1 {
		t.Fatalf("Expected FunctionCall with 1 arg, got=%T", call1Stmt.Expr)
	}
	testLiteralObject(t, call1.Args[0], "hello")

	call2Stmt, ok := program.Body[4].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got=%T", program.Body[4])
	}
	call2, ok := call2Stmt.Expr.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("Expected FunctionCall, got=%T", call2Stmt.Expr)
	}
	_, isTable := call2.Args[0].(*ast.TableLiteral)
	if !isTable {
		t.Errorf("Expected argument to be TableLiteral, got=%T", call2.Args[0])
	}
}

func TestLuauTypeAliasesAndContinue(t *testing.T) {
	input := `
		type Point = { x: number, y: number }
		export type ID = string | number

		for i = 1, 10 do
			continue
		end
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Body) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Body))
	}

	typeAlias, ok := program.Body[0].(*ast.TypeAlias)
	if !ok {
		t.Fatalf("Expected TypeAlias, got=%T", program.Body[0])
	}
	if typeAlias.Name != "Point" {
		t.Errorf("Expected type alias 'Point', got %s", typeAlias.Name)
	}
	if typeAlias.IsExport {
		t.Errorf("Expected non-exported type 'Point'")
	}

	exportedType, ok := program.Body[1].(*ast.TypeAlias)
	if !ok {
		t.Fatalf("Expected TypeAlias, got=%T", program.Body[1])
	}
	if exportedType.Name != "ID" {
		t.Errorf("Expected type alias 'ID', got %s", exportedType.Name)
	}
	if !exportedType.IsExport {
		t.Errorf("Expected exported type")
	}

	if exportedType.Type == nil || exportedType.Type.Type != "string | number" {
		t.Errorf("Expected Type annotation 'string | number'")
	}

	forLoop, ok := program.Body[2].(*ast.ForLoop)
	if !ok {
		t.Fatalf("Expected ForLoop, got=%T", program.Body[2])
	}
	if len(forLoop.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in ForLoop body, got %d", len(forLoop.Body.Statements))
	}
	_, hasContinue := forLoop.Body.Statements[0].(*ast.ContinueStatement)
	if !hasContinue {
		t.Errorf("Expected ContinueStatement inside loop body, got=%T", forLoop.Body.Statements[0])
	}
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
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Body) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Body))
	}

	typeAlias := program.Body[0].(*ast.TypeAlias)
	if len(typeAlias.Generics) != 2 || typeAlias.Generics[0] != "K" || typeAlias.Generics[1] != "V" {
		t.Errorf("Expected generics [K, V], got %v", typeAlias.Generics)
	}

	funcDef := program.Body[1].(*ast.FunctionDef)
	if len(funcDef.Generics) != 1 || funcDef.Generics[0] != "T" {
		t.Errorf("Expected generic [T], got %v", funcDef.Generics)
	}

	localFunc := program.Body[2].(*ast.LocalFunction)
	if len(localFunc.Generics) != 1 || localFunc.Generics[0] != "Value" {
		t.Errorf("Expected generic [Value], got %v", localFunc.Generics)
	}
}

func TestInterpolatedStrings(t *testing.T) {
	input := "local greeting = `Hello {\"World\"}, I am {10 + 10} years old!`"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assign := program.Body[0].(*ast.LocalAssignment)
	interp, ok := assign.Values[0].(*ast.InterpolatedString)
	if !ok {
		t.Fatalf("Expected InterpolatedString, got %T", assign.Values[0])
	}

	if len(interp.Segments) != 3 {
		t.Fatalf("Expected 3 string segments, got %d", len(interp.Segments))
	}
	if interp.Segments[0] != "Hello " || interp.Segments[1] != ", I am " || interp.Segments[2] != " years old!" {
		t.Errorf("Segments mismatch: %v", interp.Segments)
	}

	if len(interp.Expressions) != 2 {
		t.Fatalf("Expected 2 expressions, got %d", len(interp.Expressions))
	}
}

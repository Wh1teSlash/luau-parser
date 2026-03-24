package lexer

import "testing"

func checkToken(t *testing.T, index int, tok Token, expectedType TokenType, expectedLiteral string) {
	t.Helper()
	if tok.Type != expectedType {
		t.Errorf("token[%d] type wrong: expected %q, got %q (literal %q)",
			index, expectedType, tok.Type, tok.Literal)
	}
	if tok.Literal != expectedLiteral {
		t.Errorf("token[%d] literal wrong: expected %q, got %q",
			index, expectedLiteral, tok.Literal)
	}
}

func TestNextToken(t *testing.T) {
	input := `
		local x = 5
		local hex = 0xFF
		local bin = 0b1010
		local float = 3.14e-2
		local str1 = "hello"
		local str2 = [[multi
line]]
		local str3 = [=[ nested [[ brackets ]] ]=]
		-- single comment
		--[[ multi
		comment ]]
		x += 1
		y //= 2
		z ..= "!"
		type Point = { x: number }
		export type ID = string | number
		if a ~= b then continue end
		const MAX = 100
		const name: string = "luau"
	`

	tests := []struct {
		name            string
		expectedType    TokenType
		expectedLiteral string
	}{
		{"local keyword", LOCAL, "local"},
		{"variable name x", IDENT, "x"},
		{"assign operator", ASSIGN, "="},
		{"integer literal 5", INT, "5"},
		{"local keyword", LOCAL, "local"},
		{"variable name hex", IDENT, "hex"},
		{"assign operator", ASSIGN, "="},
		{"hex integer literal", INT, "0xFF"},
		{"local keyword", LOCAL, "local"},
		{"variable name bin", IDENT, "bin"},
		{"assign operator", ASSIGN, "="},
		{"binary integer literal", INT, "0b1010"},
		{"local keyword", LOCAL, "local"},
		{"variable name float", IDENT, "float"},
		{"assign operator", ASSIGN, "="},
		{"float literal with exponent", FLOAT, "3.14e-2"},
		{"local keyword", LOCAL, "local"},
		{"variable name str1", IDENT, "str1"},
		{"assign operator", ASSIGN, "="},
		{"double-quoted string", STRING, "hello"},
		{"local keyword", LOCAL, "local"},
		{"variable name str2", IDENT, "str2"},
		{"assign operator", ASSIGN, "="},
		{"multi-line string", STRING, "multi\nline"},
		{"local keyword", LOCAL, "local"},
		{"variable name str3", IDENT, "str3"},
		{"assign operator", ASSIGN, "="},
		{"level-1 long string with inner brackets", STRING, " nested [[ brackets ]] "},
		{"single-line comment", COMMENT, " single comment"},
		{"multi-line comment", COMMENT, " multi\n\t\tcomment "},
		{"variable name x", IDENT, "x"},
		{"plus-assign operator", PLUS_ASSIGN, "+="},
		{"integer literal 1", INT, "1"},
		{"variable name y", IDENT, "y"},
		{"floor-div-assign operator", FLOOR_DIV_ASSIGN, "//="},
		{"integer literal 2", INT, "2"},
		{"variable name z", IDENT, "z"},
		{"concat-assign operator", CONCAT_ASSIGN, "..="},
		{"string literal exclamation", STRING, "!"},
		{"type keyword", TYPE, "type"},
		{"type name Point", IDENT, "Point"},
		{"assign operator", ASSIGN, "="},
		{"left brace", LBRACE, "{"},
		{"field name x", IDENT, "x"},
		{"colon", COLON, ":"},
		{"type name number", IDENT, "number"},
		{"right brace", RBRACE, "}"},
		{"export keyword", EXPORT, "export"},
		{"type keyword", TYPE, "type"},
		{"type name ID", IDENT, "ID"},
		{"assign operator", ASSIGN, "="},
		{"type name string", IDENT, "string"},
		{"pipe operator", PIPE, "|"},
		{"type name number", IDENT, "number"},
		{"if keyword", IF, "if"},
		{"variable name a", IDENT, "a"},
		{"not-equal operator", NOT_EQ, "~="},
		{"variable name b", IDENT, "b"},
		{"then keyword", THEN, "then"},
		{"continue keyword", CONTINUE, "continue"},
		{"end keyword", END, "end"},
		{"const keyword", CONST, "const"},
		{"const variable MAX", IDENT, "MAX"},
		{"assign operator", ASSIGN, "="},
		{"integer literal 100", INT, "100"},
		{"const keyword", CONST, "const"},
		{"const variable name", IDENT, "name"},
		{"colon", COLON, ":"},
		{"type annotation string", IDENT, "string"},
		{"assign operator", ASSIGN, "="},
		{"string literal luau", STRING, "luau"},

		{"end of file", EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := l.NextToken()
			checkToken(t, i, tok, tt.expectedType, tt.expectedLiteral)
		})
	}
}

func TestInterpolatedStrings(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []struct {
			expectedType    TokenType
			expectedLiteral string
		}
	}{
		{
			name:  "simple interpolated string",
			input: "`hello {name}!`",
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{INTERP_BEGIN, "hello "},
				{IDENT, "name"},
				{INTERP_END, "!"},
			},
		},
		{
			name:  "interpolated string with expression",
			input: "`value: {x + 1} done`",
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{INTERP_BEGIN, "value: "},
				{IDENT, "x"},
				{PLUS, "+"},
				{INT, "1"},
				{INTERP_END, " done"},
			},
		},
		{
			name:  "interpolated string with multiple expressions",
			input: "`{a}, {b}`",
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{INTERP_BEGIN, ""},
				{IDENT, "a"},
				{INTERP_MID, ", "},
				{IDENT, "b"},
				{INTERP_END, ""},
			},
		},
		{
			name:  "plain interpolated string with no expressions",
			input: "`just text`",
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{STRING, "just text"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, expected := range tt.tokens {
				tok := l.NextToken()
				checkToken(t, i, tok, expected.expectedType, expected.expectedLiteral)
			}
		})
	}
}

func TestOperators(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{"addition", "+", PLUS, "+"},
		{"subtraction", "-", MINUS, "-"},
		{"multiplication", "*", ASTERISK, "*"},
		{"division", "/", SLASH, "/"},
		{"floor division", "//", FLOOR_DIV, "//"},
		{"modulo", "%", MODULO, "%"},
		{"exponentiation", "^", CARET, "^"},
		{"concatenation", "..", CONCAT, ".."},
		{"length", "#", HASH, "#"},
		{"equality", "==", EQ, "=="},
		{"inequality", "~=", NOT_EQ, "~="},
		{"less than", "<", LT, "<"},
		{"less than or equal", "<=", LTE, "<="},
		{"greater than", ">", GT, ">"},
		{"greater than or equal", ">=", GTE, ">="},
		{"assignment", "=", ASSIGN, "="},
		{"plus assign", "+=", PLUS_ASSIGN, "+="},
		{"minus assign", "-=", MINUS_ASSIGN, "-="},
		{"multiply assign", "*=", ASTERISK_ASSIGN, "*="},
		{"divide assign", "/=", SLASH_ASSIGN, "/="},
		{"floor div assign", "//=", FLOOR_DIV_ASSIGN, "//="},
		{"modulo assign", "%=", MODULO_ASSIGN, "%="},
		{"exponent assign", "^=", CARET_ASSIGN, "^="},
		{"concat assign", "..=", CONCAT_ASSIGN, "..="},
		{"colon", ":", COLON, ":"},
		{"double colon", "::", DOUBLE_COLON, "::"},
		{"dot", ".", DOT, "."},
		{"ellipsis", "...", ELLIPSIS, "..."},
		{"pipe", "|", PIPE, "|"},
		{"ampersand", "&", AMPERSAND, "&"},
		{"question mark", "?", QUESTION, "?"},
		{"at sign", "@", AT, "@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()
			checkToken(t, 0, tok, tt.expectedType, tt.expectedLiteral)
		})
	}
}

func TestKeywords(t *testing.T) {
	tests := []struct {
		keyword      string
		expectedType TokenType
	}{
		{"and", AND},
		{"break", BREAK},
		{"continue", CONTINUE},
		{"do", DO},
		{"else", ELSE},
		{"elseif", ELSEIF},
		{"end", END},
		{"false", FALSE},
		{"for", FOR},
		{"function", FUNCTION},
		{"if", IF},
		{"in", IN},
		{"local", LOCAL},
		{"nil", NIL},
		{"not", NOT},
		{"or", OR},
		{"repeat", REPEAT},
		{"return", RETURN},
		{"then", THEN},
		{"true", TRUE},
		{"until", UNTIL},
		{"while", WHILE},
		{"export", EXPORT},
		{"type", TYPE},
		{"const", CONST},
	}

	for _, tt := range tests {
		t.Run(tt.keyword, func(t *testing.T) {
			l := New(tt.keyword)
			tok := l.NextToken()
			checkToken(t, 0, tok, tt.expectedType, tt.keyword)
		})
	}
}

func TestConstTokenSequences(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []struct {
			expectedType    TokenType
			expectedLiteral string
		}
	}{
		{
			name:  "simple const declaration",
			input: `const x = 5`,
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{CONST, "const"},
				{IDENT, "x"},
				{ASSIGN, "="},
				{INT, "5"},
				{EOF, ""},
			},
		},
		{
			name:  "typed const declaration",
			input: `const x: number = 5`,
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{CONST, "const"},
				{IDENT, "x"},
				{COLON, ":"},
				{IDENT, "number"},
				{ASSIGN, "="},
				{INT, "5"},
				{EOF, ""},
			},
		},
		{
			name:  "multi-name const declaration",
			input: `const a, b = 1, 2`,
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{CONST, "const"},
				{IDENT, "a"},
				{COMMA, ","},
				{IDENT, "b"},
				{ASSIGN, "="},
				{INT, "1"},
				{COMMA, ","},
				{INT, "2"},
				{EOF, ""},
			},
		},
		{
			name:  "const function declaration",
			input: `const function f() end`,
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{CONST, "const"},
				{FUNCTION, "function"},
				{IDENT, "f"},
				{LPAREN, "("},
				{RPAREN, ")"},
				{END, "end"},
				{EOF, ""},
			},
		},
		{
			name:  "const not confused with identifier prefix",
			input: `constant = 5`,
			tokens: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{IDENT, "constant"},
				{ASSIGN, "="},
				{INT, "5"},
				{EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, expected := range tt.tokens {
				tok := l.NextToken()
				checkToken(t, i, tok, expected.expectedType, expected.expectedLiteral)
			}
		})
	}
}

func TestNumberLiterals(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{"plain integer", "42", INT, "42"},
		{"zero", "0", INT, "0"},
		{"hex lowercase", "0xff", INT, "0xff"},
		{"hex uppercase", "0xFF", INT, "0xFF"},
		{"hex mixed", "0xDeAdBeEf", INT, "0xDeAdBeEf"},
		{"binary", "0b1010", INT, "0b1010"},
		{"binary zeros", "0b0000", INT, "0b0000"},
		{"float simple", "3.14", FLOAT, "3.14"},
		{"float exponent lowercase", "1e10", FLOAT, "1e10"},
		{"float exponent uppercase", "1E10", FLOAT, "1E10"},
		{"float exponent negative", "3.14e-2", FLOAT, "3.14e-2"},
		{"float exponent positive", "2.5e+3", FLOAT, "2.5e+3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()
			checkToken(t, 0, tok, tt.expectedType, tt.expectedLiteral)
		})
	}
}

func TestStringLiterals(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{"double-quoted", `"hello"`, "hello"},
		{"single-quoted", `'world'`, "world"},
		{"empty double-quoted", `""`, ""},
		{"empty single-quoted", `''`, ""},
		{"escaped quote", `"say \"hi\""`, `say \"hi\"`},
		{"level-0 long string", "[[hello]]", "hello"},
		{"level-0 long string multiline", "[[line1\nline2]]", "line1\nline2"},
		{"level-1 long string", "[=[hello]=]", "hello"},
		{"level-1 long string with inner brackets", "[=[ foo [[bar]] baz ]=]", " foo [[bar]] baz "},
		{"level-2 long string", "[==[hello]==]", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()
			checkToken(t, 0, tok, STRING, tt.expectedLiteral)
		})
	}
}

func TestComments(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{"single-line comment", "-- hello world", " hello world"},
		{"single-line comment empty", "--", ""},
		{"multi-line comment", "--[[ foo\nbar ]]", " foo\nbar "},
		{"multi-line comment level-1", "--[=[ foo [[nested]] bar ]=]", " foo [[nested]] bar "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()
			checkToken(t, 0, tok, COMMENT, tt.expectedLiteral)
		})
	}
}

func TestTokenPositions(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		line   int
		column int
	}{
		{"first token column", "local", 1, 1},
		{"token after newline", "\nlocal", 2, 1},
		{"token after spaces", "   x", 1, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()
			if tok.Pos.Line != tt.line {
				t.Errorf("expected line %d, got %d", tt.line, tok.Pos.Line)
			}
			if tok.Pos.Column != tt.column {
				t.Errorf("expected column %d, got %d", tt.column, tok.Pos.Column)
			}
		})
	}
}

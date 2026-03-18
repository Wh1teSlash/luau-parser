package lexer

import (
	"testing"
)

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
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LOCAL, "local"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT, "5"},

		{LOCAL, "local"},
		{IDENT, "hex"},
		{ASSIGN, "="},
		{INT, "0xFF"},

		{LOCAL, "local"},
		{IDENT, "bin"},
		{ASSIGN, "="},
		{INT, "0b1010"},

		{LOCAL, "local"},
		{IDENT, "float"},
		{ASSIGN, "="},
		{FLOAT, "3.14e-2"},

		{LOCAL, "local"},
		{IDENT, "str1"},
		{ASSIGN, "="},
		{STRING, "hello"},

		{LOCAL, "local"},
		{IDENT, "str2"},
		{ASSIGN, "="},
		{STRING, "multi\nline"},

		{LOCAL, "local"},
		{IDENT, "str3"},
		{ASSIGN, "="},
		{STRING, " nested [[ brackets ]] "},

		{COMMENT, " single comment"},

		{COMMENT, " multi\n\t\tcomment "},

		{IDENT, "x"},
		{PLUS_ASSIGN, "+="},
		{INT, "1"},

		{IDENT, "y"},
		{FLOOR_DIV_ASSIGN, "//="},
		{INT, "2"},

		{IDENT, "z"},
		{CONCAT_ASSIGN, "..="},
		{STRING, "!"},

		{TYPE, "type"},
		{IDENT, "Point"},
		{ASSIGN, "="},
		{LBRACE, "{"},
		{IDENT, "x"},
		{COLON, ":"},
		{IDENT, "number"},
		{RBRACE, "}"},

		{EXPORT, "export"},
		{TYPE, "type"},
		{IDENT, "ID"},
		{ASSIGN, "="},
		{IDENT, "string"},
		{PIPE, "|"},
		{IDENT, "number"},

		{IF, "if"},
		{IDENT, "a"},
		{NOT_EQ, "~="},
		{IDENT, "b"},
		{THEN, "then"},
		{CONTINUE, "continue"},
		{END, "end"},

		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

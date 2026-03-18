package lexer

import "github.com/Wh1teSlash/luau-parser/ast"

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	COMMENT TokenType = "COMMENT"

	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"

	PLUS      TokenType = "+"
	MINUS     TokenType = "-"
	ASTERISK  TokenType = "*"
	SLASH     TokenType = "/"
	FLOOR_DIV TokenType = "//"
	MODULO    TokenType = "%"
	CARET     TokenType = "^"
	CONCAT    TokenType = ".."
	HASH      TokenType = "#"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "~="
	LT     TokenType = "<"
	LTE    TokenType = "<="
	GT     TokenType = ">"
	GTE    TokenType = ">="

	ASSIGN           TokenType = "="
	PLUS_ASSIGN      TokenType = "+="
	MINUS_ASSIGN     TokenType = "-="
	ASTERISK_ASSIGN  TokenType = "*="
	SLASH_ASSIGN     TokenType = "/="
	FLOOR_DIV_ASSIGN TokenType = "//="
	MODULO_ASSIGN    TokenType = "%="
	CARET_ASSIGN     TokenType = "^="
	CONCAT_ASSIGN    TokenType = "..="

	COMMA        TokenType = ","
	SEMICOLON    TokenType = ";"
	COLON        TokenType = ":"
	DOUBLE_COLON TokenType = "::"
	DOT          TokenType = "."
	ELLIPSIS     TokenType = "..."

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	ARROW     TokenType = "->"
	QUESTION  TokenType = "?"
	PIPE      TokenType = "|"
	AMPERSAND TokenType = "&"

	AND      TokenType = "AND"
	BREAK    TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"
	DO       TokenType = "DO"
	ELSE     TokenType = "ELSE"
	ELSEIF   TokenType = "ELSEIF"
	END      TokenType = "END"
	FALSE    TokenType = "FALSE"
	FOR      TokenType = "FOR"
	FUNCTION TokenType = "FUNCTION"
	IF       TokenType = "IF"
	IN       TokenType = "IN"
	LOCAL    TokenType = "LOCAL"
	NIL      TokenType = "NIL"
	NOT      TokenType = "NOT"
	OR       TokenType = "OR"
	REPEAT   TokenType = "REPEAT"
	RETURN   TokenType = "RETURN"
	THEN     TokenType = "THEN"
	TRUE     TokenType = "TRUE"
	UNTIL    TokenType = "UNTIL"
	WHILE    TokenType = "WHILE"

	EXPORT TokenType = "EXPORT"
	TYPE   TokenType = "TYPE"

	INTERP_BEGIN TokenType = "INTERP_BEGIN"
	INTERP_MID   TokenType = "INTERP_MID"
	INTERP_END   TokenType = "INTERP_END"
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     ast.Position
}

var keywords = map[string]TokenType{
	"and":      AND,
	"break":    BREAK,
	"continue": CONTINUE,
	"do":       DO,
	"else":     ELSE,
	"elseif":   ELSEIF,
	"end":      END,
	"false":    FALSE,
	"for":      FOR,
	"function": FUNCTION,
	"if":       IF,
	"in":       IN,
	"local":    LOCAL,
	"nil":      NIL,
	"not":      NOT,
	"or":       OR,
	"repeat":   REPEAT,
	"return":   RETURN,
	"then":     THEN,
	"true":     TRUE,
	"until":    UNTIL,
	"while":    WHILE,
	"export":   EXPORT,
	"type":     TYPE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

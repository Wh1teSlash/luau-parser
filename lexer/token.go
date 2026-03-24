package lexer

import "github.com/Wh1teSlash/luau-parser/ast"

type TokenType uint8

const (
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	IDENT
	INT
	FLOAT
	STRING

	PLUS
	MINUS
	ASTERISK
	SLASH
	FLOOR_DIV
	MODULO
	CARET
	CONCAT
	HASH

	EQ
	NOT_EQ
	LT
	LTE
	GT
	GTE

	ASSIGN
	PLUS_ASSIGN
	MINUS_ASSIGN
	ASTERISK_ASSIGN
	SLASH_ASSIGN
	FLOOR_DIV_ASSIGN
	MODULO_ASSIGN
	CARET_ASSIGN
	CONCAT_ASSIGN

	COMMA
	SEMICOLON
	COLON
	DOUBLE_COLON
	DOT
	ELLIPSIS

	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	ARROW
	QUESTION
	PIPE
	AMPERSAND
	AT

	AND
	BREAK
	CONTINUE
	DO
	ELSE
	ELSEIF
	END
	FALSE
	FOR
	FUNCTION
	IF
	IN
	LOCAL
	NIL
	NOT
	OR
	REPEAT
	RETURN
	THEN
	TRUE
	UNTIL
	WHILE

	EXPORT
	TYPE
	CONST

	INTERP_BEGIN
	INTERP_MID
	INTERP_END

	tokenTypeCount
)

var tokenTypeNames = [tokenTypeCount]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	PLUS:      "+",
	MINUS:     "-",
	ASTERISK:  "*",
	SLASH:     "/",
	FLOOR_DIV: "//",
	MODULO:    "%",
	CARET:     "^",
	CONCAT:    "..",
	HASH:      "#",

	EQ:     "==",
	NOT_EQ: "~=",
	LT:     "<",
	LTE:    "<=",
	GT:     ">",
	GTE:    ">=",

	ASSIGN:           "=",
	PLUS_ASSIGN:      "+=",
	MINUS_ASSIGN:     "-=",
	ASTERISK_ASSIGN:  "*=",
	SLASH_ASSIGN:     "/=",
	FLOOR_DIV_ASSIGN: "//=",
	MODULO_ASSIGN:    "%=",
	CARET_ASSIGN:     "^=",
	CONCAT_ASSIGN:    "..=",

	COMMA:        ",",
	SEMICOLON:    ";",
	COLON:        ":",
	DOUBLE_COLON: "::",
	DOT:          ".",
	ELLIPSIS:     "...",

	LPAREN:   "(",
	RPAREN:   ")",
	LBRACE:   "{",
	RBRACE:   "}",
	LBRACKET: "[",
	RBRACKET: "]",

	ARROW:     "->",
	QUESTION:  "?",
	PIPE:      "|",
	AMPERSAND: "&",
	AT:        "@",

	AND:      "and",
	BREAK:    "break",
	CONTINUE: "continue",
	DO:       "do",
	ELSE:     "else",
	ELSEIF:   "elseif",
	END:      "end",
	FALSE:    "false",
	FOR:      "for",
	FUNCTION: "function",
	IF:       "if",
	IN:       "in",
	LOCAL:    "local",
	NIL:      "nil",
	NOT:      "not",
	OR:       "or",
	REPEAT:   "repeat",
	RETURN:   "return",
	THEN:     "then",
	TRUE:     "true",
	UNTIL:    "until",
	WHILE:    "while",

	EXPORT: "export",
	TYPE:   "type",
	CONST:  "const",

	INTERP_BEGIN: "INTERP_BEGIN",
	INTERP_MID:   "INTERP_MID",
	INTERP_END:   "INTERP_END",
}

func (t TokenType) String() string {
	if t < tokenTypeCount {
		return tokenTypeNames[t]
	}
	return "UNKNOWN"
}

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
	"const":    CONST,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

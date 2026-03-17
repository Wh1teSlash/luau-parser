package parser

import "github.com/Wh1teSlash/luau-parser/lexer"

const (
	_ int = iota
	LOWEST
	OR          // or
	AND         // and
	EQUALS      // ==, ~=
	LESSGREATER // >, <, >=, <=
	CONCAT      // ..
	SUM         // +, -
	PRODUCT     // *, /, //, %
	CARET       // ^
	PREFIX      // -X, not X, #X
	CALL        // myFunction(X)
	INDEX       // array[index], obj.field
)

var precedences = map[lexer.TokenType]int{
	lexer.OR:        OR,
	lexer.AND:       AND,
	lexer.EQ:        EQUALS,
	lexer.NOT_EQ:    EQUALS,
	lexer.LT:        LESSGREATER,
	lexer.GT:        LESSGREATER,
	lexer.LTE:       LESSGREATER,
	lexer.GTE:       LESSGREATER,
	lexer.CONCAT:    CONCAT,
	lexer.PLUS:      SUM,
	lexer.MINUS:     SUM,
	lexer.SLASH:     PRODUCT,
	lexer.ASTERISK:  PRODUCT,
	lexer.FLOOR_DIV: PRODUCT,
	lexer.MODULO:    PRODUCT,
	lexer.CARET:     CARET,
	lexer.LPAREN:    CALL,
	lexer.LBRACKET:  INDEX,
	lexer.DOT:       INDEX,
	lexer.COLON:     CALL,
}

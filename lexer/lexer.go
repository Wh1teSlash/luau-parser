package lexer

import (
	"github.com/Wh1teSlash/luau-parser/ast"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte

	line   int
	column int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peekNextChar() byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+1]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	pos := ast.Position{Line: l.line, Column: l.column}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(EQ, pos)
		} else {
			tok = Token{Type: ASSIGN, Literal: string(l.ch), Pos: pos}
		}
	case '~':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(NOT_EQ, pos)
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Pos: pos}
		}
	case '<':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(LTE, pos)
		} else {
			tok = Token{Type: LT, Literal: string(l.ch), Pos: pos}
		}
	case '>':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(GTE, pos)
		} else {
			tok = Token{Type: GT, Literal: string(l.ch), Pos: pos}
		}

	case '+':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(PLUS_ASSIGN, pos)
		} else {
			tok = Token{Type: PLUS, Literal: string(l.ch), Pos: pos}
		}
	case '-':
		if l.peekChar() == '-' {
			tok.Pos = pos
			tok.Type = COMMENT

			l.readChar()
			l.readChar()

			if l.ch == '[' && (l.peekChar() == '[' || l.peekChar() == '=') {
				tok.Literal = l.readMultiLineString()
			} else {
				position := l.position
				for l.ch != '\n' && l.ch != 0 {
					l.readChar()
				}
				tok.Literal = l.input[position:l.position]
			}
			return tok
		} else if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(MINUS_ASSIGN, pos)
		} else {
			tok = Token{Type: MINUS, Literal: string(l.ch), Pos: pos}
		}
	case '*':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(ASTERISK_ASSIGN, pos)
		} else {
			tok = Token{Type: ASTERISK, Literal: string(l.ch), Pos: pos}
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			if l.peekChar() == '=' {
				tok = l.makeTwoCharToken(FLOOR_DIV_ASSIGN, pos)
				tok.Literal = "//="
			} else {
				tok = Token{Type: FLOOR_DIV, Literal: "//", Pos: pos}
			}
		} else if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(SLASH_ASSIGN, pos)
		} else {
			tok = Token{Type: SLASH, Literal: string(l.ch), Pos: pos}
		}
	case '%':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(MODULO_ASSIGN, pos)
		} else {
			tok = Token{Type: MODULO, Literal: string(l.ch), Pos: pos}
		}
	case '^':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(CARET_ASSIGN, pos)
		} else {
			tok = Token{Type: CARET, Literal: string(l.ch), Pos: pos}
		}

	case ':':
		if l.peekChar() == ':' {
			tok = l.makeTwoCharToken(DOUBLE_COLON, pos)
		} else {
			tok = Token{Type: COLON, Literal: string(l.ch), Pos: pos}
		}
	case '.':
		if l.peekChar() == '.' {
			if l.peekNextChar() == '.' {
				l.readChar()
				tok = l.makeTwoCharToken(ELLIPSIS, pos)
				tok.Literal = "..."
			} else if l.peekNextChar() == '=' {
				l.readChar()
				tok = l.makeTwoCharToken(CONCAT_ASSIGN, pos)
				tok.Literal = "..="
			} else {
				tok = l.makeTwoCharToken(CONCAT, pos)
			}
		} else {
			tok = Token{Type: DOT, Literal: string(l.ch), Pos: pos}
		}
	case '#':
		tok = Token{Type: HASH, Literal: string(l.ch), Pos: pos}
	case '?':
		tok = Token{Type: QUESTION, Literal: string(l.ch), Pos: pos}
	case '|':
		tok = Token{Type: PIPE, Literal: string(l.ch), Pos: pos}
	case '&':
		tok = Token{Type: AMPERSAND, Literal: string(l.ch), Pos: pos}

	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch), Pos: pos}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch), Pos: pos}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch), Pos: pos}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch), Pos: pos}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch), Pos: pos}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch), Pos: pos}
	case '[':
		if l.peekChar() == '[' || l.peekChar() == '=' {
			tok.Type = STRING
			tok.Literal = l.readMultiLineString()
			tok.Pos = pos
			return tok
		} else {
			tok = Token{Type: LBRACKET, Literal: string(l.ch), Pos: pos}
		}
	case ']':
		tok = Token{Type: RBRACKET, Literal: string(l.ch), Pos: pos}

	case '"', '\'':
		tok.Type = STRING
		tok.Literal = l.readString(l.ch)
		tok.Pos = pos

	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Pos = pos

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Pos = pos
			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			tok.Pos = pos
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Pos: pos}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) makeTwoCharToken(tokenType TokenType, pos ast.Position) Token {
	ch := l.ch
	l.readChar()
	literal := string(ch) + string(l.ch)
	return Token{Type: tokenType, Literal: literal, Pos: pos}
}

func (l *Lexer) peekCharOffset(offset int) byte {
	if l.readPosition+offset >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+offset]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position
	tokType := INT

	if l.ch == '0' {
		if l.peekChar() == 'x' || l.peekChar() == 'X' {
			l.readChar()
			l.readChar()
			for isHexDigit(l.ch) || l.ch == '_' {
				l.readChar()
			}
			return l.input[position:l.position], tokType
		} else if l.peekChar() == 'b' || l.peekChar() == 'B' {
			l.readChar()
			l.readChar()
			for l.ch == '0' || l.ch == '1' || l.ch == '_' {
				l.readChar()
			}
			return l.input[position:l.position], tokType
		}
	}

	for isDigit(l.ch) || l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '_' {
		if l.ch == '.' {
			tokType = FLOAT
		}
		if l.ch == 'e' || l.ch == 'E' {
			tokType = FLOAT
			if l.peekChar() == '+' || l.peekChar() == '-' {
				l.readChar()
			}
		}
		l.readChar()
	}
	return l.input[position:l.position], tokType
}

func (l *Lexer) readString(quote byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\\' {
			l.readChar()
			continue
		}
		if l.ch == quote || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readMultiLineString() string {
	l.readChar()

	eqCount := 0
	for l.ch == '=' {
		eqCount++
		l.readChar()
	}

	if l.ch == '[' {
		l.readChar()
	}

	if l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}

	position := l.position

	for l.ch != 0 {
		if l.ch == ']' {
			closingEq := 0

			for l.peekCharOffset(closingEq) == '=' {
				closingEq++
			}

			if closingEq == eqCount && l.peekCharOffset(closingEq) == ']' {
				str := l.input[position:l.position]

				for i := 0; i < closingEq+2; i++ {
					l.readChar()
				}
				return str
			}
		}

		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

package parser

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

type prefixParseFn func() ast.Expr
type infixParseFn func(ast.Expr) ast.Expr

type Parser struct {
	l      *lexer.Lexer
	errors []error

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []error{},
		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),
	}
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NIL, p.parseNilLiteral)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.HASH, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerInfix(lexer.LPAREN, p.parseFunctionCall)
	p.registerInfix(lexer.LBRACKET, p.parseIndexAccess)
	p.registerInfix(lexer.DOT, p.parseFieldAccess)
	p.registerInfix(lexer.COLON, p.parseMethodCall)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseTableLiteral)
	p.registerPrefix(lexer.FUNCTION, p.parseFunctionExpr)
	p.registerInfix(lexer.DOUBLE_COLON, p.parseTypeCast)
	p.registerPrefix(lexer.IF, p.parseIfExpr)
	p.registerPrefix(lexer.ELLIPSIS, p.parseVarArgs)
	p.registerInfix(lexer.STRING, p.parseFunctionCallStringOrTable)
	p.registerInfix(lexer.LBRACE, p.parseFunctionCallStringOrTable)

	binaryOps := []lexer.TokenType{
		lexer.PLUS, lexer.MINUS, lexer.SLASH, lexer.ASTERISK,
		lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LTE, lexer.GTE,
		lexer.CONCAT, lexer.MODULO, lexer.FLOOR_DIV, lexer.CARET,
		lexer.AND, lexer.OR,
	}
	for _, op := range binaryOps {
		p.registerInfix(op, p.parseInfixExpression)
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []error { return p.errors }

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		err := fmt.Errorf("expected next token to be %s, got %s instead at line %d, col %d",
			t, p.peekToken.Type, p.peekToken.Pos.Line, p.peekToken.Pos.Column)
		p.errors = append(p.errors, err)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Body:     []ast.Stmt{},
	}

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseBlock() *ast.Block {
	block := &ast.Block{
		BaseNode:   ast.BaseNode{Position: p.curToken.Pos},
		Statements: []ast.Stmt{},
	}

	for p.curToken.Type != lexer.END &&
		p.curToken.Type != lexer.ELSE &&
		p.curToken.Type != lexer.ELSEIF &&
		p.curToken.Type != lexer.UNTIL &&
		p.curToken.Type != lexer.EOF {

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

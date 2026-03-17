package parser

import (
	"fmt"
	"strconv"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

func (p *Parser) parseExpression(precedence int) ast.Expr {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Errorf("no parsing function for token prefix %s", p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != lexer.EOF && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Name:     p.curToken.Literal,
	}
}

func (p *Parser) parseStringLiteral() ast.Expr {
	return &ast.Literal{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Type:     "string",
		Value:    p.curToken.Literal,
	}
}

func (p *Parser) parseBooleanLiteral() ast.Expr {
	return &ast.Literal{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Type:     "boolean",
		Value:    p.curToken.Type == lexer.TRUE,
	}
}

func (p *Parser) parseNilLiteral() ast.Expr {
	return &ast.Literal{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Type:     "nil",
		Value:    nil,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expr {
	lit := &ast.Literal{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Type:     "number",
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("unable to parse %q as int", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expr {
	expression := &ast.UnaryOp{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Op:       p.curToken.Literal,
	}

	p.nextToken()
	expression.Operand = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expr {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, fmt.Errorf("expected closing bracket, got %s", p.peekToken.Type))
		return nil
	}
	p.nextToken()

	return &ast.ParenExpr{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Expr:     exp,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expr) ast.Expr {
	expression := &ast.BinaryOp{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Left:     left,
		Op:       p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	return expression
}

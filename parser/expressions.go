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

func (p *Parser) parseFloatLiteral() ast.Expr {
	lit := &ast.Literal{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Type:     "number",
	}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("unable to parse %q as float", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseTableLiteral() ast.Expr {
	table := &ast.TableLiteral{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Fields:   []*ast.TableField{},
	}

	for p.peekToken.Type != lexer.RBRACE && p.peekToken.Type != lexer.EOF {
		p.nextToken()

		if p.curToken.Type == lexer.RBRACE {
			break
		}

		field := &ast.TableField{}

		if p.curToken.Type == lexer.LBRACKET {
			p.nextToken()
			field.Key = p.parseExpression(LOWEST)
			if !p.expectPeek(lexer.RBRACKET) || !p.expectPeek(lexer.ASSIGN) {
				return nil
			}
			p.nextToken()
			field.Value = p.parseExpression(LOWEST)

		} else if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.ASSIGN {
			field.Key = &ast.Literal{
				BaseNode: ast.BaseNode{Position: p.curToken.Pos},
				Type:     "string",
				Value:    p.curToken.Literal,
			}
			p.nextToken()
			p.nextToken()
			field.Value = p.parseExpression(LOWEST)

		} else {
			field.Value = p.parseExpression(LOWEST)
		}

		table.Fields = append(table.Fields, field)

		if p.peekToken.Type == lexer.COMMA || p.peekToken.Type == lexer.SEMICOLON {
			p.nextToken()
		} else if p.peekToken.Type != lexer.RBRACE {
			p.errors = append(p.errors, fmt.Errorf("expected ',' or ';' or '}' in table literal, got %s", p.peekToken.Type))
			return nil
		}
	}

	if p.curToken.Type != lexer.RBRACE {
		if !p.expectPeek(lexer.RBRACE) {
			return nil
		}
	}

	return table
}

func (p *Parser) parseFunctionExpr() ast.Expr {
	expr := &ast.FunctionExpr{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
	}

	expr.Generics = p.parseGenericParams()

	params, returnType := p.parseFunctionSignature()

	expr.Parameters = params
	expr.ReturnType = returnType
	expr.Body = p.parseBlock()

	return expr
}

func (p *Parser) parseIndexAccess(left ast.Expr) ast.Expr {
	expr := &ast.IndexAccess{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Table:    left,
	}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return expr
}

func (p *Parser) parseFieldAccess(left ast.Expr) ast.Expr {
	expr := &ast.FieldAccess{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Object:   left,
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	expr.Field = p.curToken.Literal

	return expr
}

func (p *Parser) parseMethodCall(left ast.Expr) ast.Expr {
	expr := &ast.MethodCall{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Object:   left,
		Args:     []ast.Expr{},
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	expr.Method = p.curToken.Literal

	switch p.peekToken.Type {
	case lexer.LPAREN:
		p.nextToken()
		expr.Args = p.parseCallArguments()
	case lexer.LBRACE, lexer.STRING:
		p.nextToken()
		expr.Args = append(expr.Args, p.parseExpression(LOWEST))
	default:
		p.errors = append(p.errors, fmt.Errorf("expected function arguments for method %s", expr.Method))
		return nil
	}

	return expr
}

func (p *Parser) parseCallArguments() []ast.Expr {
	args := []ast.Expr{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return args
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

	if p.curToken.Type == lexer.CARET || p.curToken.Type == lexer.CONCAT {
		precedence--
	}

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseTypeCast(left ast.Expr) ast.Expr {
	expr := &ast.TypeCast{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Value:    left,
	}

	p.nextToken()

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Errorf("expected type identifier after ::, got %s", p.curToken.Type))
		return nil
	}

	expr.Type = &ast.TypeAnnotation{Type: p.curToken.Literal}

	return expr
}

func (p *Parser) parseIfExpr() ast.Expr {
	expr := &ast.IfExpr{BaseNode: ast.BaseNode{Position: p.curToken.Pos}}
	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.THEN) {
		return nil
	}
	p.nextToken()

	expr.Then = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.ELSE) {
		return nil
	}
	p.nextToken()

	expr.Else = p.parseExpression(LOWEST)
	return expr
}

func (p *Parser) parseVarArgs() ast.Expr {
	return &ast.VarArgs{BaseNode: ast.BaseNode{Position: p.curToken.Pos}}
}

func (p *Parser) parseFunctionCallStringOrTable(left ast.Expr) ast.Expr {
	call := &ast.FunctionCall{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Function: left,
		Args:     []ast.Expr{},
	}

	switch p.curToken.Type {
	case lexer.STRING:
		call.Args = append(call.Args, p.parseStringLiteral())
	case lexer.LBRACE:
		call.Args = append(call.Args, p.parseTableLiteral())
	}

	return call
}

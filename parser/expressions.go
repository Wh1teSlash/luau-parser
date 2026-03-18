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

func (p *Parser) parseInterpolatedString() ast.Expr {
	pos := p.curToken.Pos
	segments := []string{p.curToken.Literal}
	expressions := []ast.Expr{}

	for {
		p.nextToken()

		expr := p.parseExpression(LOWEST)
		if expr != nil {
			expressions = append(expressions, expr)
		}

		if p.peekToken.Type == lexer.INTERP_MID {
			p.nextToken()
			segments = append(segments, p.curToken.Literal)
		} else if p.peekToken.Type == lexer.INTERP_END {
			p.nextToken()
			segments = append(segments, p.curToken.Literal)
			break
		} else {
			p.errors = append(p.errors, fmt.Errorf("expected INTERP_MID or INTERP_END, got %s", p.peekToken.Type))
			break
		}
	}

	return p.factory.InterpolatedString(pos, segments, expressions)
}

func (p *Parser) parseIdentifier() ast.Expr {
	return p.factory.Identifier(p.curToken.Pos, p.curToken.Literal)
}

func (p *Parser) parseStringLiteral() ast.Expr {
	return p.factory.Literal(p.curToken.Pos, "string", p.curToken.Literal)
}

func (p *Parser) parseBooleanLiteral() ast.Expr {
	return p.factory.Literal(p.curToken.Pos, "boolean", p.curToken.Type == lexer.TRUE)
}

func (p *Parser) parseNilLiteral() ast.Expr {
	return p.factory.Literal(p.curToken.Pos, "nil", nil)
}

func (p *Parser) parseIntegerLiteral() ast.Expr {
	pos := p.curToken.Pos
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("unable to parse %q as int", p.curToken.Literal))
		return nil
	}
	return p.factory.Literal(pos, "number", value)
}

func (p *Parser) parseFloatLiteral() ast.Expr {
	pos := p.curToken.Pos
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("unable to parse %q as float", p.curToken.Literal))
		return nil
	}
	return p.factory.Literal(pos, "number", value)
}

func (p *Parser) parseTableLiteral() ast.Expr {
	pos := p.curToken.Pos
	fields := []*ast.TableField{}

	for p.peekToken.Type != lexer.RBRACE && p.peekToken.Type != lexer.EOF {
		p.nextToken()

		if p.curToken.Type == lexer.RBRACE {
			break
		}

		var key, value ast.Expr

		if p.curToken.Type == lexer.LBRACKET {
			p.nextToken()
			key = p.parseExpression(LOWEST)
			if !p.expectPeek(lexer.RBRACKET) || !p.expectPeek(lexer.ASSIGN) {
				return nil
			}
			p.nextToken()
			value = p.parseExpression(LOWEST)

		} else if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.ASSIGN {
			key = p.factory.Literal(p.curToken.Pos, "string", p.curToken.Literal)
			p.nextToken()
			p.nextToken()
			value = p.parseExpression(LOWEST)

		} else {
			value = p.parseExpression(LOWEST)
		}

		fields = append(fields, p.factory.TableField(key, value))

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

	return p.factory.TableLiteral(pos, fields)
}

func (p *Parser) parseFunctionExpr() ast.Expr {
	pos := p.curToken.Pos
	generics := p.parseGenericParams()
	params, returnType := p.parseFunctionSignature()

	p.nextToken()
	body := p.parseBlock()

	return p.factory.FunctionExpr(pos, params, body,
		ast.WithExprGenerics(generics...),
		ast.WithExprReturnType(returnType),
	)
}

func (p *Parser) parseIndexAccess(left ast.Expr) ast.Expr {
	pos := p.curToken.Pos
	p.nextToken()
	index := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}
	return p.factory.IndexAccess(pos, left, index)
}

func (p *Parser) parseFieldAccess(left ast.Expr) ast.Expr {
	pos := p.curToken.Pos
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	return p.factory.FieldAccess(pos, left, p.curToken.Literal)
}

func (p *Parser) parseMethodCall(left ast.Expr) ast.Expr {
	pos := p.curToken.Pos
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	method := p.curToken.Literal
	args := []ast.Expr{}

	switch p.peekToken.Type {
	case lexer.LPAREN:
		p.nextToken()
		args = p.parseCallArguments()
	case lexer.LBRACE, lexer.STRING:
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	default:
		p.errors = append(p.errors, fmt.Errorf("expected function arguments for method %s", method))
		return nil
	}

	return p.factory.MethodCall(pos, left, method, args)
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
	pos := p.curToken.Pos
	op := p.curToken.Literal
	p.nextToken()
	operand := p.parseExpression(PREFIX)

	return p.factory.UnaryOp(pos, op, operand)
}

func (p *Parser) parseGroupedExpression() ast.Expr {
	pos := p.curToken.Pos
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, fmt.Errorf("expected closing bracket, got %s", p.peekToken.Type))
		return nil
	}
	p.nextToken()

	return p.factory.ParenExpr(pos, exp)
}

func (p *Parser) parseInfixExpression(left ast.Expr) ast.Expr {
	pos := p.curToken.Pos
	op := p.curToken.Literal
	precedence := p.curPrecedence()

	if p.curToken.Type == lexer.CARET || p.curToken.Type == lexer.CONCAT {
		precedence--
	}

	p.nextToken()
	right := p.parseExpression(precedence)

	return p.factory.BinaryOp(pos, left, op, right)
}

func (p *Parser) parseTypeCast(left ast.Expr) ast.Expr {
	pos := p.curToken.Pos
	p.nextToken()
	typeNode := p.parseType(TYPE_LOWEST)
	return p.factory.TypeCast(pos, left, typeNode)
}

func (p *Parser) parseIfExpr() ast.Expr {
	pos := p.curToken.Pos
	elseIfs := []*ast.ElseIfExprClause{}

	p.nextToken()
	condition := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.THEN) {
		return nil
	}
	p.nextToken()
	then := p.parseExpression(LOWEST)

	for p.peekToken.Type == lexer.ELSEIF {
		p.nextToken()
		p.nextToken()
		clauseCond := p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.THEN) {
			return nil
		}
		p.nextToken()
		clauseThen := p.parseExpression(LOWEST)

		elseIfs = append(elseIfs, p.factory.ElseIfExprClause(clauseCond, clauseThen))
	}

	if !p.expectPeek(lexer.ELSE) {
		return nil
	}
	p.nextToken()
	elseExpr := p.parseExpression(LOWEST)

	return p.factory.IfExpr(
		pos,
		condition,
		then,
		ast.WithElseIfExprs(elseIfs...),
		ast.WithElseExpr(elseExpr),
	)
}

func (p *Parser) parseVarArgs() ast.Expr {
	return p.factory.VarArgs(p.curToken.Pos)
}

func (p *Parser) parseFunctionCallStringOrTable(left ast.Expr) ast.Expr {
	pos := p.curToken.Pos
	args := []ast.Expr{}

	switch p.curToken.Type {
	case lexer.STRING:
		args = append(args, p.parseStringLiteral())
	case lexer.LBRACE:
		args = append(args, p.parseTableLiteral())
	}

	return p.factory.FunctionCall(pos, left, args)
}

func (p *Parser) parseFunctionCall(function ast.Expr) ast.Expr {
	pos := p.curToken.Pos
	args := []ast.Expr{}

	if p.peekToken.Type != lexer.RPAREN {
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
		for p.peekToken.Type == lexer.COMMA {
			p.nextToken()
			p.nextToken()
			args = append(args, p.parseExpression(LOWEST))
		}
	}

	p.expectPeek(lexer.RPAREN)
	return p.factory.FunctionCall(pos, function, args)
}

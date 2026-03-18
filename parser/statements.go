package parser

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case lexer.SEMICOLON:
		return p.factory.EmptyStatement(p.curToken.Pos)
	case lexer.LOCAL:
		return p.parseLocalStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.REPEAT:
		return p.parseRepeatLoop()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.FUNCTION:
		return p.parseFunctionStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.COMMENT:
		return p.parseComment()
	case lexer.BREAK:
		return p.factory.BreakStatement(p.curToken.Pos)
	case lexer.CONTINUE:
		return p.factory.ContinueStatement(p.curToken.Pos)
	case lexer.DO:
		return p.parseDoBlock()
	case lexer.TYPE:
		return p.parseTypeAlias(false)
	case lexer.EXPORT:
		return p.parseExportStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseGenericParams() []string {
	var generics []string

	if p.peekToken.Type != lexer.LT {
		return generics
	}
	p.nextToken()

	if p.peekToken.Type == lexer.GT {
		p.nextToken()
		return generics
	}

	p.nextToken()
	for {
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Errorf("expected IDENT in generics, got %s", p.curToken.Type))
			return nil
		}
		generics = append(generics, p.curToken.Literal)

		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken()
		p.nextToken()
	}

	if !p.expectPeek(lexer.GT) {
		return nil
	}
	return generics
}

func (p *Parser) parseTypeAlias(isExport bool) ast.Stmt {
	pos := p.curToken.Pos
	p.nextToken()

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Errorf("expected Identifier, got %s", p.curToken.Type))
		return nil
	}
	name := p.curToken.Literal

	generics := p.parseGenericParams()

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}
	p.nextToken()

	typeNode := p.parseType(TYPE_LOWEST)

	return p.factory.TypeAlias(pos, name, generics, typeNode, isExport)
}

func (p *Parser) parseExportStatement() ast.Stmt {
	if p.peekToken.Type == lexer.TYPE {
		p.nextToken()
		return p.parseTypeAlias(true)
	}
	p.errors = append(p.errors, fmt.Errorf("expected TYPE after EXPORT"))
	return nil
}

func (p *Parser) parseFunctionStatement() ast.Stmt {
	pos := p.curToken.Pos
	p.nextToken()

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Errorf("expected IDENT for function name, got %s", p.curToken.Type))
		return nil
	}

	name := p.curToken.Literal
	isMethod := false

	for p.peekToken.Type == lexer.DOT {
		p.nextToken()
		name += "."
		p.nextToken()
		name += p.curToken.Literal
	}

	if p.peekToken.Type == lexer.COLON {
		p.nextToken()
		name += ":"
		p.nextToken()
		name += p.curToken.Literal
		isMethod = true
	}

	generics := p.parseGenericParams()
	params, returnType := p.parseFunctionSignature()
	p.nextToken()
	body := p.parseBlock()

	if p.curToken.Type != lexer.END {
		p.errors = append(p.errors, fmt.Errorf("expected END to close function statement, got %s", p.curToken.Type))
	}

	if isMethod {
		return p.factory.MetamethodDef(pos, name, params, body)
	}

	return p.factory.FunctionDef(pos, name, generics, params, body, returnType)
}

func (p *Parser) parseLocalStatement() ast.Stmt {
	pos := p.curToken.Pos

	if p.peekToken.Type == lexer.FUNCTION {
		p.nextToken()
		p.nextToken()

		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Errorf("expected Identifier, got %s", p.curToken.Type))
			return nil
		}
		name := p.curToken.Literal
		generics := p.parseGenericParams()
		params, returnType := p.parseFunctionSignature()
		p.nextToken()
		body := p.parseBlock()

		if p.curToken.Type != lexer.END {
			p.errors = append(p.errors, fmt.Errorf("expected END, got %s", p.curToken.Type))
		}

		return p.factory.LocalFunction(pos, name, generics, params, body, returnType)
	}

	names := []string{}
	values := []ast.Expr{}
	types := []ast.TypeNode{}

	for {
		p.nextToken()
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Errorf("expected Identifier, got %s", p.curToken.Type))
			return nil
		}
		names = append(names, p.curToken.Literal)

		if p.peekToken.Type == lexer.COLON {
			p.nextToken()
			p.nextToken()
			types = append(types, p.parseType(TYPE_LOWEST))
		} else {
			types = append(types, nil)
		}

		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken()
	}

	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken()
		p.nextToken()

		for {
			val := p.parseExpression(LOWEST)
			if val != nil {
				values = append(values, val)
			}
			if p.peekToken.Type != lexer.COMMA {
				break
			}
			p.nextToken()
			p.nextToken()
		}
	}

	return p.factory.LocalAssignment(pos, names, types, values)
}

func (p *Parser) parseDoBlock() ast.Stmt {
	pos := p.curToken.Pos
	p.nextToken()
	body := p.parseBlock()

	if p.curToken.Type != lexer.END {
		p.errors = append(p.errors, fmt.Errorf("expected END to close do block, got %s", p.curToken.Type))
	}
	return p.factory.DoBlock(pos, body)
}

func (p *Parser) parseComment() ast.Stmt {
	return p.factory.Comment(p.curToken.Pos, p.curToken.Literal)
}

func (p *Parser) parseForStatement() ast.Stmt {
	pos := p.curToken.Pos
	p.nextToken()

	var names []string
	names = append(names, p.curToken.Literal)

	if p.peekToken.Type == lexer.COMMA || p.peekToken.Type == lexer.IN {
		for p.peekToken.Type == lexer.COMMA {
			p.nextToken()
			p.nextToken()
			names = append(names, p.curToken.Literal)
		}
		p.expectPeek(lexer.IN)
		p.nextToken()

		iterables := []ast.Expr{p.parseExpression(LOWEST)}

		p.expectPeek(lexer.DO)
		p.nextToken()

		return p.factory.ForInLoop(pos, names, iterables, p.parseBlock())
	} else {
		p.expectPeek(lexer.ASSIGN)
		p.nextToken()
		start := p.parseExpression(LOWEST)
		p.expectPeek(lexer.COMMA)
		p.nextToken()
		end := p.parseExpression(LOWEST)

		var step ast.Expr
		if p.peekToken.Type == lexer.COMMA {
			p.nextToken()
			p.nextToken()
			step = p.parseExpression(LOWEST)
		}

		p.expectPeek(lexer.DO)
		p.nextToken()

		return p.factory.ForLoop(pos, names[0], start, end, step, p.parseBlock())
	}
}

func (p *Parser) parseWhileStatement() ast.Stmt {
	pos := p.curToken.Pos
	p.nextToken()
	condition := p.parseExpression(LOWEST)

	p.expectPeek(lexer.DO)
	p.nextToken()
	body := p.parseBlock()
	return p.factory.WhileLoop(pos, condition, body)
}

func (p *Parser) parseRepeatLoop() ast.Stmt {
	pos := p.curToken.Pos
	p.nextToken()
	body := p.parseBlock()

	if p.curToken.Type != lexer.UNTIL {
		p.errors = append(p.errors, fmt.Errorf("expected current token to be UNTIL, got %s", p.curToken.Type))
		return nil
	}
	p.nextToken()

	condition := p.parseExpression(LOWEST)
	return p.factory.RepeatLoop(pos, body, condition)
}

func (p *Parser) parseFunctionSignature() ([]*ast.Parameter, ast.TypeNode) {
	params := []*ast.Parameter{}
	p.expectPeek(lexer.LPAREN)

	if p.peekToken.Type != lexer.RPAREN {
		p.nextToken()
		for {
			name := p.curToken.Literal
			var typeNode ast.TypeNode

			if p.peekToken.Type == lexer.COLON {
				p.nextToken()
				p.nextToken()
				typeNode = p.parseType(TYPE_LOWEST)
			}

			params = append(params, p.factory.Parameter(name, typeNode))

			if p.peekToken.Type != lexer.COMMA {
				break
			}
			p.nextToken()
			p.nextToken()
		}
	}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
	}

	var returnType ast.TypeNode
	if p.peekToken.Type == lexer.COLON || p.peekToken.Type == lexer.ARROW {
		p.nextToken()
		p.nextToken()
		returnType = p.parseType(TYPE_LOWEST)
	}

	return params, returnType
}

func (p *Parser) parseIfStatement() ast.Stmt {
	pos := p.curToken.Pos
	elseIfs := []*ast.ElseIfClause{}

	p.nextToken()
	condition := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.THEN) {
		return nil
	}
	p.nextToken()
	thenBlock := p.parseBlock()

	for p.curToken.Type == lexer.ELSEIF {
		p.nextToken()
		clauseCond := p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.THEN) {
			return nil
		}
		p.nextToken()
		clauseBody := p.parseBlock()

		elseIfs = append(elseIfs, p.factory.ElseIfClause(clauseCond, clauseBody))
	}

	var elseBlock *ast.Block
	if p.curToken.Type == lexer.ELSE {
		p.nextToken()
		elseBlock = p.parseBlock()
	}

	if p.curToken.Type != lexer.END {
		p.errors = append(p.errors, fmt.Errorf("expected END to close if statement, got %s", p.curToken.Type))
	}

	return p.factory.IfStatement(pos, condition, thenBlock, elseIfs, elseBlock)
}

func (p *Parser) parseReturnStatement() ast.Stmt {
	pos := p.curToken.Pos
	values := []ast.Expr{}

	if p.peekToken.Type == lexer.EOF || p.peekToken.Type == lexer.END || p.peekToken.Type == lexer.ELSE || p.peekToken.Type == lexer.ELSEIF || p.peekToken.Type == lexer.UNTIL {
		return p.factory.ReturnStatement(pos, values)
	}

	p.nextToken()
	for {
		values = append(values, p.parseExpression(LOWEST))
		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken()
		p.nextToken()
	}
	return p.factory.ReturnStatement(pos, values)
}

func (p *Parser) parseExpressionStatement() ast.Stmt {
	targets := []ast.Expr{p.parseExpression(LOWEST)}

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		targets = append(targets, p.parseExpression(LOWEST))
	}

	if p.isAssignmentOperator(p.peekToken.Type) {
		op := p.peekToken.Literal
		p.nextToken()
		p.nextToken()

		values := []ast.Expr{p.parseExpression(LOWEST)}

		for p.peekToken.Type == lexer.COMMA {
			p.nextToken()
			p.nextToken()
			values = append(values, p.parseExpression(LOWEST))
		}

		return p.factory.Assignment(targets[0].Pos(), targets, op, values)
	}

	if len(targets) == 1 {
		switch targets[0].(type) {
		case *ast.FunctionCall, *ast.MethodCall:
			return p.factory.ExpressionStatement(targets[0].Pos(), targets[0])
		}
	}

	err := fmt.Errorf("syntax error: expected assignment or function call at line %d, col %d",
		p.curToken.Pos.Line, p.curToken.Pos.Column)
	p.errors = append(p.errors, err)

	return nil
}

func (p *Parser) isAssignmentOperator(t lexer.TokenType) bool {
	switch t {
	case lexer.ASSIGN, lexer.PLUS_ASSIGN, lexer.MINUS_ASSIGN,
		lexer.ASTERISK_ASSIGN, lexer.SLASH_ASSIGN, lexer.FLOOR_DIV_ASSIGN,
		lexer.MODULO_ASSIGN, lexer.CARET_ASSIGN, lexer.CONCAT_ASSIGN:
		return true
	default:
		return false
	}
}

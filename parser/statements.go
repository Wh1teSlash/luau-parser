package parser

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case lexer.SEMICOLON:
		stmt := &ast.EmptyStatement{BaseNode: ast.BaseNode{Position: p.curToken.Pos}}
		return stmt
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
		return &ast.BreakStatement{BaseNode: ast.BaseNode{Position: p.curToken.Pos}}
	case lexer.CONTINUE:
		return &ast.ContinueStatement{BaseNode: ast.BaseNode{Position: p.curToken.Pos}}
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

	return &ast.TypeAlias{
		BaseNode: ast.BaseNode{Position: pos},
		Name:     name,
		Generics: generics,
		Type:     typeNode,
		IsExport: isExport,
	}
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
		return &ast.MetamethodDef{
			BaseNode:   ast.BaseNode{Position: pos},
			Name:       name,
			Parameters: params,
			Body:       body,
		}
	}

	return &ast.FunctionDef{
		BaseNode:   ast.BaseNode{Position: pos},
		Name:       name,
		Generics:   generics,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
	}
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

		return &ast.LocalFunction{
			BaseNode:   ast.BaseNode{Position: pos},
			Name:       name,
			Generics:   generics,
			Parameters: params,
			ReturnType: returnType,
			Body:       body,
		}
	}

	stmt := &ast.LocalAssignment{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Names:    []string{},
		Values:   []ast.Expr{},
		Types:    []ast.TypeNode{},
	}

	for {
		p.nextToken()
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Errorf("expected Identifier, got %s", p.curToken.Type))
			return nil
		}
		stmt.Names = append(stmt.Names, p.curToken.Literal)

		if p.peekToken.Type == lexer.COLON {
			p.nextToken()
			p.nextToken()
			stmt.Types = append(stmt.Types, p.parseType(TYPE_LOWEST))
		} else {
			stmt.Types = append(stmt.Types, nil)
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
				stmt.Values = append(stmt.Values, val)
			}
			if p.peekToken.Type != lexer.COMMA {
				break
			}
			p.nextToken()
			p.nextToken()
		}
	}

	return stmt
}

func (p *Parser) parseDoBlock() ast.Stmt {
	stmt := &ast.DoBlock{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
	}
	p.nextToken()
	stmt.Body = p.parseBlock()

	if p.curToken.Type != lexer.END {
		p.errors = append(p.errors, fmt.Errorf("expected END to close do block, got %s", p.curToken.Type))
	}
	return stmt
}

func (p *Parser) parseComment() ast.Stmt {
	stmt := &ast.Comment{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Text:     p.curToken.Literal,
	}

	return stmt
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
		return &ast.ForInLoop{BaseNode: ast.BaseNode{Position: pos}, Variables: names, Iterables: iterables, Body: p.parseBlock()}
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
		return &ast.ForLoop{BaseNode: ast.BaseNode{Position: pos}, Variable: names[0], Start: start, End: end, Step: step, Body: p.parseBlock()}
	}
}

func (p *Parser) parseWhileStatement() ast.Stmt {
	stmt := &ast.WhileLoop{BaseNode: ast.BaseNode{Position: p.curToken.Pos}}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	p.expectPeek(lexer.DO)
	p.nextToken()
	stmt.Body = p.parseBlock()
	return stmt
}

func (p *Parser) parseRepeatLoop() ast.Stmt {
	stmt := &ast.RepeatLoop{BaseNode: ast.BaseNode{Position: p.curToken.Pos}}
	p.nextToken()
	stmt.Body = p.parseBlock()

	if p.curToken.Type != lexer.UNTIL {
		p.errors = append(p.errors, fmt.Errorf("expected current token to be UNTIL, got %s", p.curToken.Type))
		return nil
	}
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseFunctionCall(function ast.Expr) ast.Expr {
	call := &ast.FunctionCall{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Function: function,
		Args:     []ast.Expr{},
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.nextToken()
		call.Args = append(call.Args, p.parseExpression(LOWEST))
		for p.peekToken.Type == lexer.COMMA {
			p.nextToken()
			p.nextToken()
			call.Args = append(call.Args, p.parseExpression(LOWEST))
		}
	}

	p.expectPeek(lexer.RPAREN)
	return call
}

func (p *Parser) parseFunctionSignature() ([]*ast.Parameter, ast.TypeNode) {
	params := []*ast.Parameter{}
	p.expectPeek(lexer.LPAREN)

	if p.peekToken.Type != lexer.RPAREN {
		p.nextToken()
		for {
			param := &ast.Parameter{Name: p.curToken.Literal}
			if p.peekToken.Type == lexer.COLON {
				p.nextToken()
				p.nextToken()
				param.Type = p.parseType(TYPE_LOWEST)
			}
			params = append(params, param)

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
	stmt := &ast.IfStatement{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		ElseIfs:  []*ast.ElseIfClause{},
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.THEN) {
		return nil
	}
	p.nextToken()
	stmt.Then = p.parseBlock()

	for p.curToken.Type == lexer.ELSEIF {
		clause := &ast.ElseIfClause{}
		p.nextToken()
		clause.Condition = p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.THEN) {
			return nil
		}
		p.nextToken()
		clause.Body = p.parseBlock()
		stmt.ElseIfs = append(stmt.ElseIfs, clause)
	}

	if p.curToken.Type == lexer.ELSE {
		p.nextToken()
		stmt.Else = p.parseBlock()
	}

	if p.curToken.Type != lexer.END {
		p.errors = append(p.errors, fmt.Errorf("expected END to close if statement, got %s", p.curToken.Type))
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Stmt {
	stmt := &ast.ReturnStatement{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Values:   []ast.Expr{},
	}

	if p.peekToken.Type == lexer.EOF || p.peekToken.Type == lexer.END || p.peekToken.Type == lexer.ELSE || p.peekToken.Type == lexer.ELSEIF || p.peekToken.Type == lexer.UNTIL {
		return stmt
	}

	p.nextToken()
	for {
		stmt.Values = append(stmt.Values, p.parseExpression(LOWEST))
		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken()
		p.nextToken()
	}
	return stmt
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

		return &ast.Assignment{
			BaseNode: ast.BaseNode{Position: targets[0].Pos()},
			Targets:  targets,
			Operator: op,
			Values:   values,
		}
	}

	if len(targets) == 1 {
		switch targets[0].(type) {
		case *ast.FunctionCall, *ast.MethodCall:
			return &ast.ExpressionStatement{
				BaseNode: ast.BaseNode{Position: targets[0].Pos()},
				Expr:     targets[0],
			}
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

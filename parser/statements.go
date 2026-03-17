package parser

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case lexer.LOCAL:
		return p.parseLocalStatement()
	case lexer.RETURN:
		// return p.parseReturnStatement()
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLocalStatement() ast.Stmt {
	stmt := &ast.LocalAssignment{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Names:    []string{},
		Values:   []ast.Expr{},
		Types:    []*ast.TypeAnnotation{},
	}

	for {
		p.nextToken()
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Errorf("ожидался идентификатор, получено %s", p.curToken.Type))
			return nil
		}
		stmt.Names = append(stmt.Names, p.curToken.Literal)

		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken()
	}

	if p.peekToken.Type != lexer.ASSIGN {
		return stmt
	}

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

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Stmt {
	expr := p.parseExpression(LOWEST)

	if expr != nil {
		_ = expr
	}

	return nil
}

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

	// 1. Читаем список имен переменных
	for {
		p.nextToken()
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Errorf("ожидался идентификатор, получено %s", p.curToken.Type))
			return nil
		}
		stmt.Names = append(stmt.Names, p.curToken.Literal)

		// TODO: здесь можно добавить логику парсинга аннотаций типов (например, `local x: number`)

		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken() // пропускаем запятую
	}

	if p.peekToken.Type != lexer.ASSIGN {
		return stmt // Просто объявление (local x)
	}

	p.nextToken() // Съедаем последний идентификатор
	p.nextToken() // Съедаем `=`

	// 2. Читаем список значений
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

	// TODO: Ожидаем точку с запятой (если она есть)

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Stmt {
	expr := p.parseExpression(LOWEST)

	// Временная заглушка, чтобы код компилировался.
	// Настоящий парсер здесь проверит, не является ли это Assignment, FunctionCall или MethodCall.
	if expr != nil {
		_ = expr
	}

	return nil
}

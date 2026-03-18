package parser

import (
	"fmt"

	"github.com/Wh1teSlash/luau-parser/ast"
	"github.com/Wh1teSlash/luau-parser/lexer"
)

func (p *Parser) peekTypePrecedence() int {
	if p, ok := typePrecedences[p.peekToken.Type]; ok {
		return p
	}
	return TYPE_LOWEST
}

func (p *Parser) parseType(precedence int) ast.TypeNode {
	prefix := p.prefixTypeFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Errorf("no type parsing function for %s", p.curToken.Type))
		return nil
	}
	leftType := prefix()

	for p.peekToken.Type != lexer.EOF && precedence < p.peekTypePrecedence() {
		infix := p.infixTypeFns[p.peekToken.Type]
		if infix == nil {
			return leftType
		}

		p.nextToken()
		leftType = infix(leftType)
	}

	return leftType
}

func (p *Parser) parsePrimitiveType() ast.TypeNode {
	return p.factory.PrimitiveType(p.curToken.Pos, p.curToken.Literal)
}

func (p *Parser) parseUnionType(left ast.TypeNode) ast.TypeNode {
	pos := p.curToken.Pos
	precedence := TYPE_UNION
	p.nextToken()
	right := p.parseType(precedence)

	return p.factory.UnionType(pos, left, right)
}

func (p *Parser) parseOptionalType(left ast.TypeNode) ast.TypeNode {
	return p.factory.OptionalType(p.curToken.Pos, left)
}

func (p *Parser) parseTableType() ast.TypeNode {
	pos := p.curToken.Pos
	fields := []*ast.TableTypeField{}

	if p.peekToken.Type == lexer.RBRACE {
		p.nextToken()
		return p.factory.TableType(pos, fields)
	}

	for p.peekToken.Type != lexer.RBRACE && p.peekToken.Type != lexer.EOF {
		p.nextToken()

		var key, value ast.TypeNode
		var keyName string
		isAccess := false

		if p.curToken.Type == lexer.LBRACKET {
			p.nextToken()
			key = p.parseType(TYPE_LOWEST)
			p.expectPeek(lexer.RBRACKET)
			p.expectPeek(lexer.COLON)
			p.nextToken()
			value = p.parseType(TYPE_LOWEST)
			isAccess = true
		} else if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.COLON {
			keyName = p.curToken.Literal
			p.nextToken()
			p.nextToken()
			value = p.parseType(TYPE_LOWEST)
		} else {
			value = p.parseType(TYPE_LOWEST)
		}

		fields = append(fields, p.factory.TableTypeField(key, keyName, value, isAccess))

		if p.peekToken.Type == lexer.COMMA || p.peekToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
	}
	p.expectPeek(lexer.RBRACE)
	return p.factory.TableType(pos, fields)
}

func (p *Parser) parseGenericType(left ast.TypeNode) ast.TypeNode {
	pos := p.curToken.Pos
	types := []ast.TypeNode{}

	p.nextToken()
	types = append(types, p.parseType(TYPE_LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		types = append(types, p.parseType(TYPE_LOWEST))
	}

	p.expectPeek(lexer.GT)
	return p.factory.GenericType(pos, left, types)
}

func (p *Parser) parseParenType() ast.TypeNode {
	p.nextToken()
	t := p.parseType(TYPE_LOWEST)
	p.expectPeek(lexer.RPAREN)
	return t
}

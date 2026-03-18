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
	return &ast.PrimitiveType{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Name:     p.curToken.Literal,
	}
}

func (p *Parser) parseUnionType(left ast.TypeNode) ast.TypeNode {
	node := &ast.UnionType{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Left:     left,
	}

	precedence := TYPE_UNION
	p.nextToken()
	node.Right = p.parseType(precedence)

	return node
}

func (p *Parser) parseOptionalType(left ast.TypeNode) ast.TypeNode {
	return &ast.OptionalType{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		BaseType: left,
	}
}

func (p *Parser) parseTableType() ast.TypeNode {
	node := &ast.TableType{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		Fields:   []*ast.TableTypeField{},
	}

	if p.peekToken.Type == lexer.RBRACE {
		p.nextToken()
		return node
	}

	for p.peekToken.Type != lexer.RBRACE && p.peekToken.Type != lexer.EOF {
		p.nextToken()
		field := &ast.TableTypeField{}

		if p.curToken.Type == lexer.LBRACKET {
			p.nextToken()
			field.Key = p.parseType(TYPE_LOWEST)
			p.expectPeek(lexer.RBRACKET)
			p.expectPeek(lexer.COLON)
			p.nextToken()
			field.Value = p.parseType(TYPE_LOWEST)
			field.IsAccess = true
		} else if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.COLON {
			field.KeyName = p.curToken.Literal
			p.nextToken()
			p.nextToken()
			field.Value = p.parseType(TYPE_LOWEST)
		} else {
			field.Value = p.parseType(TYPE_LOWEST)
		}

		node.Fields = append(node.Fields, field)
		if p.peekToken.Type == lexer.COMMA || p.peekToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
	}
	p.expectPeek(lexer.RBRACE)
	return node
}

func (p *Parser) parseGenericType(left ast.TypeNode) ast.TypeNode {
	node := &ast.GenericType{
		BaseNode: ast.BaseNode{Position: p.curToken.Pos},
		BaseType: left,
		Types:    []ast.TypeNode{},
	}

	p.nextToken()
	node.Types = append(node.Types, p.parseType(TYPE_LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		node.Types = append(node.Types, p.parseType(TYPE_LOWEST))
	}

	p.expectPeek(lexer.GT)
	return node
}

func (p *Parser) parseParenType() ast.TypeNode {
	p.nextToken()
	t := p.parseType(TYPE_LOWEST)
	p.expectPeek(lexer.RPAREN)
	return t
}

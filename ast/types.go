package ast

import "fmt"

type TypeNode interface {
	Node
	typeNode()
}

type PrimitiveType struct {
	BaseNode
	Name string
}

func (p *PrimitiveType) String() string       { return p.Name }
func (p *PrimitiveType) Accept(v Visitor) any { return v.VisitPrimitiveType(p) }
func (p *PrimitiveType) typeNode()            {}

type UnionType struct {
	BaseNode
	Left  TypeNode
	Right TypeNode
}

func (u *UnionType) String() string       { return fmt.Sprintf("%s | %s", u.Left.String(), u.Right.String()) }
func (u *UnionType) Accept(v Visitor) any { return v.VisitUnionType(u) }
func (u *UnionType) typeNode()            {}

type OptionalType struct {
	BaseNode
	BaseType TypeNode
}

func (o *OptionalType) String() string       { return o.BaseType.String() + "?" }
func (o *OptionalType) Accept(v Visitor) any { return v.VisitOptionalType(o) }
func (o *OptionalType) typeNode()            {}

type TableType struct {
	BaseNode
	Fields []*TableTypeField
}

func (t *TableType) String() string       { return "table" }
func (t *TableType) Accept(v Visitor) any { return v.VisitTableType(t) }
func (t *TableType) typeNode()            {}

type TableTypeField struct {
	Key      TypeNode
	KeyName  string
	Value    TypeNode
	IsAccess bool
}

type GenericType struct {
	BaseNode
	BaseType TypeNode
	Types    []TypeNode
}

func (g *GenericType) String() string       { return "generic" }
func (g *GenericType) Accept(v Visitor) any { return v.VisitGenericType(g) }
func (g *GenericType) typeNode()            {}

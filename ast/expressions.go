package ast

import "fmt"

type Identifier struct {
	BaseNode
	Name string
}

func (i *Identifier) String() string       { return fmt.Sprintf("Identifier{%s}", i.Name) }
func (i *Identifier) Accept(v Visitor) any { return v.VisitIdentifier(i) }
func (i *Identifier) expressionNode()      {}

type Literal struct {
	BaseNode
	Type  string // "number", "string", "boolean", "nil"
	Value any
}

func (l *Literal) String() string {
	return fmt.Sprintf("Literal{type: %s, value: %v}", l.Type, l.Value)
}
func (l *Literal) Accept(v Visitor) any { return v.VisitLiteral(l) }
func (l *Literal) expressionNode()      {}

type BinaryOp struct {
	BaseNode
	Left  Expr
	Op    string // "+", "-", "*", "/", "==", ">", etc.
	Right Expr
}

func (b *BinaryOp) String() string       { return fmt.Sprintf("BinaryOp{op: %s}", b.Op) }
func (b *BinaryOp) Accept(v Visitor) any { return v.VisitBinaryOp(b) }
func (b *BinaryOp) expressionNode()      {}

type UnaryOp struct {
	BaseNode
	Op      string // "-", "not", "#"
	Operand Expr
}

func (u *UnaryOp) String() string       { return fmt.Sprintf("UnaryOp{op: %s}", u.Op) }
func (u *UnaryOp) Accept(v Visitor) any { return v.VisitUnaryOp(u) }
func (u *UnaryOp) expressionNode()      {}

type FunctionCall struct {
	BaseNode
	Function Expr // may be Identifier or other Expression
	Args     []Expr
}

func (f *FunctionCall) String() string       { return fmt.Sprintf("FunctionCall{args: %d}", len(f.Args)) }
func (f *FunctionCall) Accept(v Visitor) any { return v.VisitFunctionCall(f) }
func (f *FunctionCall) expressionNode()      {}

type MethodCall struct {
	BaseNode
	Object Expr
	Method string
	Args   []Expr
}

func (m *MethodCall) String() string {
	return fmt.Sprintf("MethodCall{method: %s, args: %d}", m.Method, len(m.Args))
}
func (m *MethodCall) Accept(v Visitor) any { return v.VisitMethodCall(m) }
func (m *MethodCall) expressionNode()      {}

type IndexAccess struct {
	BaseNode
	Table Expr
	Index Expr
}

func (i *IndexAccess) String() string       { return "IndexAccess" }
func (i *IndexAccess) Accept(v Visitor) any { return v.VisitIndexAccess(i) }
func (i *IndexAccess) expressionNode()      {}

type FieldAccess struct {
	BaseNode
	Object Expr
	Field  string
}

func (f *FieldAccess) String() string       { return fmt.Sprintf("FieldAccess{field: %s}", f.Field) }
func (f *FieldAccess) Accept(v Visitor) any { return v.VisitFieldAccess(f) }
func (f *FieldAccess) expressionNode()      {}

type TableLiteral struct {
	BaseNode
	Fields []*TableField
}

func (t *TableLiteral) String() string       { return fmt.Sprintf("TableLiteral{fields: %d}", len(t.Fields)) }
func (t *TableLiteral) Accept(v Visitor) any { return v.VisitTableLiteral(t) }
func (t *TableLiteral) expressionNode()      {}

type TableField struct {
	Key   Expr // may nil, if its an array
	Value Expr
}

type FunctionExpr struct {
	BaseNode
	Generics   []string
	Parameters []*Parameter
	Body       *Block
	ReturnType TypeNode
}

func (f *FunctionExpr) String() string {
	return fmt.Sprintf("FunctionExpr{params: %d}", len(f.Parameters))
}
func (f *FunctionExpr) Accept(v Visitor) any { return v.VisitFunctionExpr(f) }
func (f *FunctionExpr) expressionNode()      {}

type Parameter struct {
	Name string
	Type TypeNode
}

type TypeAnnotation struct {
	Type string // "number", "string", "boolean", "nil", "any", "table", and etc.
}

type TypeCast struct {
	BaseNode
	Value Expr
	Type  TypeNode
}

func (t *TypeCast) String() string       { return "TypeCast" }
func (t *TypeCast) Accept(v Visitor) any { return v.VisitTypeCast(t) }
func (t *TypeCast) expressionNode()      {}

type IfExpr struct {
	BaseNode
	Condition Expr
	Then      Expr
	ElseIfs   []*ElseIfExprClause
	Else      Expr // may be nil
}

func (i *IfExpr) String() string       { return "IfExpr" }
func (i *IfExpr) Accept(v Visitor) any { return v.VisitIfExpr(i) }
func (i *IfExpr) expressionNode()      {}

type IfExprOption func(*IfExpr)

func WithElseIfExprs(clauses ...*ElseIfExprClause) IfExprOption {
	return func(i *IfExpr) {
		i.ElseIfs = append(i.ElseIfs, clauses...)
	}
}

func WithElseExpr(elseExpr Expr) IfExprOption {
	return func(i *IfExpr) {
		i.Else = elseExpr
	}
}

type ElseIfExprClause struct {
	Condition Expr
	Then      Expr
}

type VarArgs struct {
	BaseNode
}

func (v *VarArgs) String() string         { return "VarArgs" }
func (v *VarArgs) Accept(vis Visitor) any { return vis.VisitVarArgs(v) }
func (v *VarArgs) expressionNode()        {}

type ParenExpr struct {
	BaseNode
	Expr Expr
}

func (p *ParenExpr) String() string       { return "ParenExpr" }
func (p *ParenExpr) Accept(v Visitor) any { return v.VisitParenExpr(p) }
func (p *ParenExpr) expressionNode()      {}

type InterpolatedString struct {
	BaseNode
	Segments    []string
	Expressions []Expr
}

func (i *InterpolatedString) String() string       { return "InterpolatedString" }
func (i *InterpolatedString) Accept(v Visitor) any { return v.VisitInterpolatedString(i) }
func (i *InterpolatedString) expressionNode()      {}

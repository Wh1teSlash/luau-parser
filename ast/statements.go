package ast

import "fmt"

type Assignment struct {
	BaseNode
	Targets  []Expr
	Operator string
	Values   []Expr
}

func (a *Assignment) String() string       { return fmt.Sprintf("Assignment{%s}", a.Operator) }
func (a *Assignment) Accept(v Visitor) any { return v.VisitAssignment(a) }
func (a *Assignment) statementNode()       {}

type LocalAssignment struct {
	BaseNode
	Names  []string
	Types  []TypeNode
	Values []Expr
}

func (l *LocalAssignment) String() string {
	return fmt.Sprintf("LocalAssignment{names: %d}", len(l.Names))
}
func (l *LocalAssignment) Accept(v Visitor) any { return v.VisitLocalAssignment(l) }
func (l *LocalAssignment) statementNode()       {}

type IfStatement struct {
	BaseNode
	Condition Expr
	Then      *Block
	ElseIfs   []*ElseIfClause
	Else      *Block // may be nil
}

func (i *IfStatement) String() string       { return "IfStatement" }
func (i *IfStatement) Accept(v Visitor) any { return v.VisitIfStatement(i) }
func (i *IfStatement) statementNode()       {}

type ElseIfClause struct {
	Condition Expr
	Body      *Block
}

type WhileLoop struct {
	BaseNode
	Condition Expr
	Body      *Block
}

func (w *WhileLoop) String() string       { return "WhileLoop" }
func (w *WhileLoop) Accept(v Visitor) any { return v.VisitWhileLoop(w) }
func (w *WhileLoop) statementNode()       {}

type RepeatLoop struct {
	BaseNode
	Body      *Block
	Condition Expr
}

func (r *RepeatLoop) String() string       { return "RepeatLoop" }
func (r *RepeatLoop) Accept(v Visitor) any { return v.VisitRepeatLoop(r) }
func (r *RepeatLoop) statementNode()       {}

type ForLoop struct {
	BaseNode
	Variable string
	Start    Expr
	End      Expr
	Step     Expr // may be nil, by default: 1
	Body     *Block
}

func (f *ForLoop) String() string       { return fmt.Sprintf("ForLoop{var: %s}", f.Variable) }
func (f *ForLoop) Accept(v Visitor) any { return v.VisitForLoop(f) }
func (f *ForLoop) statementNode()       {}

type ForInLoop struct {
	BaseNode
	Variables []string
	Iterables []Expr
	Body      *Block
}

func (f *ForInLoop) String() string       { return fmt.Sprintf("ForInLoop{vars: %d}", len(f.Variables)) }
func (f *ForInLoop) Accept(v Visitor) any { return v.VisitForInLoop(f) }
func (f *ForInLoop) statementNode()       {}

type DoBlock struct {
	BaseNode
	Body *Block
}

func (d *DoBlock) String() string       { return "DoBlock" }
func (d *DoBlock) Accept(v Visitor) any { return v.VisitDoBlock(d) }
func (d *DoBlock) statementNode()       {}

type FunctionDef struct {
	BaseNode
	Name       string
	Generics   []string
	Parameters []*Parameter
	Body       *Block
	ReturnType TypeNode
}

func (f *FunctionDef) String() string {
	return fmt.Sprintf("FunctionDef{name: %s, params: %d}", f.Name, len(f.Parameters))
}
func (f *FunctionDef) Accept(v Visitor) any { return v.VisitFunctionDef(f) }
func (f *FunctionDef) statementNode()       {}

type LocalFunction struct {
	BaseNode
	Name       string
	Generics   []string
	Parameters []*Parameter
	Body       *Block
	ReturnType TypeNode
}

func (l *LocalFunction) String() string {
	return fmt.Sprintf("LocalFunction{name: %s, params: %d}", l.Name, len(l.Parameters))
}
func (l *LocalFunction) Accept(v Visitor) any { return v.VisitLocalFunction(l) }
func (l *LocalFunction) statementNode()       {}

type ReturnStatement struct {
	BaseNode
	Values []Expr
}

func (r *ReturnStatement) String() string       { return "ReturnStatement" }
func (r *ReturnStatement) Accept(v Visitor) any { return v.VisitReturnStatement(r) }
func (r *ReturnStatement) statementNode()       {}

type BreakStatement struct {
	BaseNode
}

func (b *BreakStatement) String() string       { return "BreakStatement" }
func (b *BreakStatement) Accept(v Visitor) any { return v.VisitBreakStatement(b) }
func (b *BreakStatement) statementNode()       {}

type ContinueStatement struct {
	BaseNode
}

func (c *ContinueStatement) String() string       { return "ContinueStatement" }
func (c *ContinueStatement) Accept(v Visitor) any { return v.VisitContinueStatement(c) }
func (c *ContinueStatement) statementNode()       {}

type TypeAlias struct {
	BaseNode
	Name     string
	Generics []string
	Type     TypeNode
	IsExport bool
}

func (t *TypeAlias) String() string       { return fmt.Sprintf("TypeAlias{name: %s}", t.Name) }
func (t *TypeAlias) Accept(v Visitor) any { return v.VisitTypeAlias(t) }
func (t *TypeAlias) statementNode()       {}

type MetamethodDef struct {
	BaseNode
	Name       string
	Parameters []*Parameter
	Body       *Block
}

func (m *MetamethodDef) String() string       { return fmt.Sprintf("MetamethodDef{name: %s}", m.Name) }
func (m *MetamethodDef) Accept(v Visitor) any { return v.VisitMetamethodDef(m) }
func (m *MetamethodDef) statementNode()       {}

type EmptyStatement struct {
	BaseNode
}

func (e *EmptyStatement) String() string       { return "EmptyStatement" }
func (e *EmptyStatement) Accept(v Visitor) any { return v.VisitEmptyStatement(e) }
func (e *EmptyStatement) statementNode()       {}

type ExpressionStatement struct {
	BaseNode
	Expr Expr
}

func (e *ExpressionStatement) String() string       { return e.Expr.String() }
func (e *ExpressionStatement) Accept(v Visitor) any { return v.VisitExpressionStatement(e) }
func (e *ExpressionStatement) statementNode()       {}

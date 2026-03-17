package ast

type Position struct {
	Line   int
	Column int
}

type Node interface {
	String() string
	Pos() Position
	Accept(visitor Visitor) any
}

type Expr interface {
	Node
	expressionNode()
}

type Stmt interface {
	Node
	statementNode()
}

type BaseNode struct {
	Position Position
}

func (b *BaseNode) Pos() Position {
	return b.Position
}

type Visitor interface {
	// Expressions
	VisitIdentifier(node *Identifier) any
	VisitLiteral(node *Literal) any
	VisitBinaryOp(node *BinaryOp) any
	VisitUnaryOp(node *UnaryOp) any
	VisitFunctionCall(node *FunctionCall) any
	VisitMethodCall(node *MethodCall) any
	VisitIndexAccess(node *IndexAccess) any
	VisitFieldAccess(node *FieldAccess) any
	VisitTableLiteral(node *TableLiteral) any
	VisitFunctionExpr(node *FunctionExpr) any
	VisitTypeCast(node *TypeCast) any
	VisitIfExpr(node *IfExpr) any
	VisitVarArgs(node *VarArgs) any
	VisitParenExpr(node *ParenExpr) any

	// Statements
	VisitAssignment(node *Assignment) any
	VisitLocalAssignment(node *LocalAssignment) any
	VisitIfStatement(node *IfStatement) any
	VisitWhileLoop(node *WhileLoop) any
	VisitRepeatLoop(node *RepeatLoop) any
	VisitForLoop(node *ForLoop) any
	VisitForInLoop(node *ForInLoop) any
	VisitDoBlock(node *DoBlock) any
	VisitFunctionDef(node *FunctionDef) any
	VisitLocalFunction(node *LocalFunction) any
	VisitReturnStatement(node *ReturnStatement) any
	VisitBreakStatement(node *BreakStatement) any
	VisitContinueStatement(node *ContinueStatement) any
	VisitTypeAlias(node *TypeAlias) any
	VisitMetamethodDef(node *MetamethodDef) any
	VisitEmptyStatement(node *EmptyStatement) any
	VisitExpressionStatement(node *ExpressionStatement) any

	// Composites
	VisitBlock(node *Block) any
	VisitProgram(node *Program) any
	VisitComment(node *Comment) any
	VisitModule(node *Module) any
}

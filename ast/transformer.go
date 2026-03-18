package ast

type Transformer interface {
	TransformIdentifier(node *Identifier) Expr
	TransformLiteral(node *Literal) Expr
	TransformBinaryOp(node *BinaryOp) Expr
	TransformUnaryOp(node *UnaryOp) Expr
	TransformFunctionCall(node *FunctionCall) Expr
	TransformMethodCall(node *MethodCall) Expr
	TransformIndexAccess(node *IndexAccess) Expr
	TransformFieldAccess(node *FieldAccess) Expr
	TransformTableLiteral(node *TableLiteral) Expr
	TransformFunctionExpr(node *FunctionExpr) Expr
	TransformTypeCast(node *TypeCast) Expr
	TransformIfExpr(node *IfExpr) Expr
	TransformVarArgs(node *VarArgs) Expr
	TransformParenExpr(node *ParenExpr) Expr
	TransformInterpolatedString(node *InterpolatedString) Expr

	TransformAssignment(node *Assignment) Stmt
	TransformLocalAssignment(node *LocalAssignment) Stmt
	TransformIfStatement(node *IfStatement) Stmt
	TransformWhileLoop(node *WhileLoop) Stmt
	TransformRepeatLoop(node *RepeatLoop) Stmt
	TransformForLoop(node *ForLoop) Stmt
	TransformForInLoop(node *ForInLoop) Stmt
	TransformDoBlock(node *DoBlock) Stmt
	TransformFunctionDef(node *FunctionDef) Stmt
	TransformLocalFunction(node *LocalFunction) Stmt
	TransformReturnStatement(node *ReturnStatement) Stmt
	TransformBreakStatement(node *BreakStatement) Stmt
	TransformContinueStatement(node *ContinueStatement) Stmt
	TransformTypeAlias(node *TypeAlias) Stmt
	TransformMetamethodDef(node *MetamethodDef) Stmt
	TransformEmptyStatement(node *EmptyStatement) Stmt
	TransformExpressionStatement(node *ExpressionStatement) Stmt
	TransformComment(node *Comment) Stmt

	TransformBlock(node *Block) *Block
	TransformProgram(node *Program) *Program

	TransformExpr(node Expr) Expr
	TransformStmt(node Stmt) Stmt
	TransformTypeNode(node TypeNode) TypeNode
}

type BaseTransformer struct{}

func (t *BaseTransformer) TransformExpr(node Expr) Expr {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *Identifier:
		return t.TransformIdentifier(n)
	case *Literal:
		return t.TransformLiteral(n)
	case *BinaryOp:
		return t.TransformBinaryOp(n)
	case *UnaryOp:
		return t.TransformUnaryOp(n)
	case *FunctionCall:
		return t.TransformFunctionCall(n)
	case *MethodCall:
		return t.TransformMethodCall(n)
	case *IndexAccess:
		return t.TransformIndexAccess(n)
	case *FieldAccess:
		return t.TransformFieldAccess(n)
	case *TableLiteral:
		return t.TransformTableLiteral(n)
	case *FunctionExpr:
		return t.TransformFunctionExpr(n)
	case *TypeCast:
		return t.TransformTypeCast(n)
	case *IfExpr:
		return t.TransformIfExpr(n)
	case *VarArgs:
		return t.TransformVarArgs(n)
	case *ParenExpr:
		return t.TransformParenExpr(n)
	case *InterpolatedString:
		return t.TransformInterpolatedString(n)
	default:
		return node
	}
}

func (t *BaseTransformer) TransformStmt(node Stmt) Stmt {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *Assignment:
		return t.TransformAssignment(n)
	case *LocalAssignment:
		return t.TransformLocalAssignment(n)
	case *IfStatement:
		return t.TransformIfStatement(n)
	case *WhileLoop:
		return t.TransformWhileLoop(n)
	case *RepeatLoop:
		return t.TransformRepeatLoop(n)
	case *ForLoop:
		return t.TransformForLoop(n)
	case *ForInLoop:
		return t.TransformForInLoop(n)
	case *DoBlock:
		return t.TransformDoBlock(n)
	case *FunctionDef:
		return t.TransformFunctionDef(n)
	case *LocalFunction:
		return t.TransformLocalFunction(n)
	case *ReturnStatement:
		return t.TransformReturnStatement(n)
	case *BreakStatement:
		return t.TransformBreakStatement(n)
	case *ContinueStatement:
		return t.TransformContinueStatement(n)
	case *TypeAlias:
		return t.TransformTypeAlias(n)
	case *MetamethodDef:
		return t.TransformMetamethodDef(n)
	case *EmptyStatement:
		return t.TransformEmptyStatement(n)
	case *ExpressionStatement:
		return t.TransformExpressionStatement(n)
	case *Comment:
		return t.TransformComment(n)
	default:
		return node
	}
}

func (t *BaseTransformer) TransformTypeNode(node TypeNode) TypeNode {
	return node
}

func (t *BaseTransformer) TransformIdentifier(node *Identifier) Expr { return node }
func (t *BaseTransformer) TransformLiteral(node *Literal) Expr       { return node }
func (t *BaseTransformer) TransformVarArgs(node *VarArgs) Expr       { return node }

func (t *BaseTransformer) TransformBinaryOp(node *BinaryOp) Expr {
	node.Left = t.TransformExpr(node.Left)
	node.Right = t.TransformExpr(node.Right)
	return node
}

func (t *BaseTransformer) TransformUnaryOp(node *UnaryOp) Expr {
	node.Operand = t.TransformExpr(node.Operand)
	return node
}

func (t *BaseTransformer) TransformFunctionCall(node *FunctionCall) Expr {
	node.Function = t.TransformExpr(node.Function)
	for i, arg := range node.Args {
		node.Args[i] = t.TransformExpr(arg)
	}
	return node
}

func (t *BaseTransformer) TransformMethodCall(node *MethodCall) Expr {
	node.Object = t.TransformExpr(node.Object)
	for i, arg := range node.Args {
		node.Args[i] = t.TransformExpr(arg)
	}
	return node
}

func (t *BaseTransformer) TransformIndexAccess(node *IndexAccess) Expr {
	node.Table = t.TransformExpr(node.Table)
	node.Index = t.TransformExpr(node.Index)
	return node
}

func (t *BaseTransformer) TransformFieldAccess(node *FieldAccess) Expr {
	node.Object = t.TransformExpr(node.Object)
	return node
}

func (t *BaseTransformer) TransformTableLiteral(node *TableLiteral) Expr {
	for _, field := range node.Fields {
		if field.Key != nil {
			field.Key = t.TransformExpr(field.Key)
		}
		field.Value = t.TransformExpr(field.Value)
	}
	return node
}

func (t *BaseTransformer) TransformFunctionExpr(node *FunctionExpr) Expr {
	node.Body = t.TransformBlock(node.Body)
	node.ReturnType = t.TransformTypeNode(node.ReturnType)
	return node
}

func (t *BaseTransformer) TransformTypeCast(node *TypeCast) Expr {
	node.Value = t.TransformExpr(node.Value)
	node.Type = t.TransformTypeNode(node.Type)
	return node
}

func (t *BaseTransformer) TransformIfExpr(node *IfExpr) Expr {
	node.Condition = t.TransformExpr(node.Condition)
	node.Then = t.TransformExpr(node.Then)
	for _, clause := range node.ElseIfs {
		clause.Condition = t.TransformExpr(clause.Condition)
		clause.Then = t.TransformExpr(clause.Then)
	}
	if node.Else != nil {
		node.Else = t.TransformExpr(node.Else)
	}
	return node
}

func (t *BaseTransformer) TransformParenExpr(node *ParenExpr) Expr {
	node.Expr = t.TransformExpr(node.Expr)
	return node
}

func (t *BaseTransformer) TransformInterpolatedString(node *InterpolatedString) Expr {
	for i, expr := range node.Expressions {
		node.Expressions[i] = t.TransformExpr(expr)
	}
	return node
}

func (t *BaseTransformer) TransformBreakStatement(node *BreakStatement) Stmt       { return node }
func (t *BaseTransformer) TransformContinueStatement(node *ContinueStatement) Stmt { return node }
func (t *BaseTransformer) TransformEmptyStatement(node *EmptyStatement) Stmt       { return node }
func (t *BaseTransformer) TransformComment(node *Comment) Stmt                     { return node }

func (t *BaseTransformer) TransformAssignment(node *Assignment) Stmt {
	for i, target := range node.Targets {
		node.Targets[i] = t.TransformExpr(target)
	}
	for i, value := range node.Values {
		node.Values[i] = t.TransformExpr(value)
	}
	return node
}

func (t *BaseTransformer) TransformLocalAssignment(node *LocalAssignment) Stmt {
	for i, value := range node.Values {
		node.Values[i] = t.TransformExpr(value)
	}
	for i, typ := range node.Types {
		node.Types[i] = t.TransformTypeNode(typ)
	}
	return node
}

func (t *BaseTransformer) TransformIfStatement(node *IfStatement) Stmt {
	node.Condition = t.TransformExpr(node.Condition)
	node.Then = t.TransformBlock(node.Then)
	for _, clause := range node.ElseIfs {
		clause.Condition = t.TransformExpr(clause.Condition)
		clause.Body = t.TransformBlock(clause.Body)
	}
	if node.Else != nil {
		node.Else = t.TransformBlock(node.Else)
	}
	return node
}

func (t *BaseTransformer) TransformWhileLoop(node *WhileLoop) Stmt {
	node.Condition = t.TransformExpr(node.Condition)
	node.Body = t.TransformBlock(node.Body)
	return node
}

func (t *BaseTransformer) TransformRepeatLoop(node *RepeatLoop) Stmt {
	node.Body = t.TransformBlock(node.Body)
	node.Condition = t.TransformExpr(node.Condition)
	return node
}

func (t *BaseTransformer) TransformForLoop(node *ForLoop) Stmt {
	node.Start = t.TransformExpr(node.Start)
	node.End = t.TransformExpr(node.End)
	if node.Step != nil {
		node.Step = t.TransformExpr(node.Step)
	}
	node.Body = t.TransformBlock(node.Body)
	return node
}

func (t *BaseTransformer) TransformForInLoop(node *ForInLoop) Stmt {
	for i, iter := range node.Iterables {
		node.Iterables[i] = t.TransformExpr(iter)
	}
	node.Body = t.TransformBlock(node.Body)
	return node
}

func (t *BaseTransformer) TransformDoBlock(node *DoBlock) Stmt {
	node.Body = t.TransformBlock(node.Body)
	return node
}

func (t *BaseTransformer) TransformFunctionDef(node *FunctionDef) Stmt {
	node.Body = t.TransformBlock(node.Body)
	node.ReturnType = t.TransformTypeNode(node.ReturnType)
	return node
}

func (t *BaseTransformer) TransformLocalFunction(node *LocalFunction) Stmt {
	node.Body = t.TransformBlock(node.Body)
	node.ReturnType = t.TransformTypeNode(node.ReturnType)
	return node
}

func (t *BaseTransformer) TransformReturnStatement(node *ReturnStatement) Stmt {
	for i, value := range node.Values {
		node.Values[i] = t.TransformExpr(value)
	}
	return node
}

func (t *BaseTransformer) TransformTypeAlias(node *TypeAlias) Stmt {
	node.Type = t.TransformTypeNode(node.Type)
	return node
}

func (t *BaseTransformer) TransformMetamethodDef(node *MetamethodDef) Stmt {
	node.Body = t.TransformBlock(node.Body)
	return node
}

func (t *BaseTransformer) TransformExpressionStatement(node *ExpressionStatement) Stmt {
	node.Expr = t.TransformExpr(node.Expr)
	return node
}

func (t *BaseTransformer) TransformBlock(node *Block) *Block {
	if node == nil {
		return nil
	}
	for i, stmt := range node.Statements {
		node.Statements[i] = t.TransformStmt(stmt)
	}
	return node
}

func (t *BaseTransformer) TransformProgram(node *Program) *Program {
	for i, stmt := range node.Body {
		node.Body[i] = t.TransformStmt(stmt)
	}
	return node
}

package ast

type NodeFactory struct {
	identifiers         *TypedArena[Identifier]
	literals            *TypedArena[Literal]
	binaryOps           *TypedArena[BinaryOp]
	unaryOps            *TypedArena[UnaryOp]
	functionCalls       *TypedArena[FunctionCall]
	methodCalls         *TypedArena[MethodCall]
	indexAccesses       *TypedArena[IndexAccess]
	fieldAccesses       *TypedArena[FieldAccess]
	tableLiterals       *TypedArena[TableLiteral]
	functionExprs       *TypedArena[FunctionExpr]
	typeCasts           *TypedArena[TypeCast]
	ifExprs             *TypedArena[IfExpr]
	varArgs             *TypedArena[VarArgs]
	parenExprs          *TypedArena[ParenExpr]
	interpolatedStrings *TypedArena[InterpolatedString]

	assignments          *TypedArena[Assignment]
	localAssignments     *TypedArena[LocalAssignment]
	ifStatements         *TypedArena[IfStatement]
	whileLoops           *TypedArena[WhileLoop]
	repeatLoops          *TypedArena[RepeatLoop]
	forLoops             *TypedArena[ForLoop]
	forInLoops           *TypedArena[ForInLoop]
	doBlocks             *TypedArena[DoBlock]
	functionDefs         *TypedArena[FunctionDef]
	localFunctions       *TypedArena[LocalFunction]
	returnStatements     *TypedArena[ReturnStatement]
	breakStatements      *TypedArena[BreakStatement]
	continueStatements   *TypedArena[ContinueStatement]
	typeAliases          *TypedArena[TypeAlias]
	metamethodDefs       *TypedArena[MetamethodDef]
	emptyStatements      *TypedArena[EmptyStatement]
	expressionStatements *TypedArena[ExpressionStatement]

	blocks     *TypedArena[Block]
	programs   *TypedArena[Program]
	comments   *TypedArena[Comment]
	modules    *TypedArena[Module]
	attributes *TypedArena[Attribute]

	primitiveTypes *TypedArena[PrimitiveType]
	unionTypes     *TypedArena[UnionType]
	optionalTypes  *TypedArena[OptionalType]
	tableTypes     *TypedArena[TableType]
	genericTypes   *TypedArena[GenericType]

	parameters        *TypedArena[Parameter]
	tableFields       *TypedArena[TableField]
	elseIfClauses     *TypedArena[ElseIfClause]
	elseIfExprClauses *TypedArena[ElseIfExprClause]
	tableTypeFields   *TypedArena[TableTypeField]
}

func NewFactory() *NodeFactory {
	return &NodeFactory{
		identifiers:         NewTypedArena[Identifier](),
		literals:            NewTypedArena[Literal](),
		binaryOps:           NewTypedArena[BinaryOp](),
		unaryOps:            NewTypedArena[UnaryOp](),
		functionCalls:       NewTypedArena[FunctionCall](),
		methodCalls:         NewTypedArena[MethodCall](),
		indexAccesses:       NewTypedArena[IndexAccess](),
		fieldAccesses:       NewTypedArena[FieldAccess](),
		tableLiterals:       NewTypedArena[TableLiteral](),
		functionExprs:       NewTypedArena[FunctionExpr](),
		typeCasts:           NewTypedArena[TypeCast](),
		ifExprs:             NewTypedArena[IfExpr](),
		varArgs:             NewTypedArena[VarArgs](),
		parenExprs:          NewTypedArena[ParenExpr](),
		interpolatedStrings: NewTypedArena[InterpolatedString](),

		assignments:          NewTypedArena[Assignment](),
		localAssignments:     NewTypedArena[LocalAssignment](),
		ifStatements:         NewTypedArena[IfStatement](),
		whileLoops:           NewTypedArena[WhileLoop](),
		repeatLoops:          NewTypedArena[RepeatLoop](),
		forLoops:             NewTypedArena[ForLoop](),
		forInLoops:           NewTypedArena[ForInLoop](),
		doBlocks:             NewTypedArena[DoBlock](),
		functionDefs:         NewTypedArena[FunctionDef](),
		localFunctions:       NewTypedArena[LocalFunction](),
		returnStatements:     NewTypedArena[ReturnStatement](),
		breakStatements:      NewTypedArena[BreakStatement](),
		continueStatements:   NewTypedArena[ContinueStatement](),
		typeAliases:          NewTypedArena[TypeAlias](),
		metamethodDefs:       NewTypedArena[MetamethodDef](),
		emptyStatements:      NewTypedArena[EmptyStatement](),
		expressionStatements: NewTypedArena[ExpressionStatement](),

		blocks:     NewTypedArena[Block](),
		programs:   NewTypedArena[Program](),
		comments:   NewTypedArena[Comment](),
		modules:    NewTypedArena[Module](),
		attributes: NewTypedArena[Attribute](),

		primitiveTypes: NewTypedArena[PrimitiveType](),
		unionTypes:     NewTypedArena[UnionType](),
		optionalTypes:  NewTypedArena[OptionalType](),
		tableTypes:     NewTypedArena[TableType](),
		genericTypes:   NewTypedArena[GenericType](),

		parameters:        NewTypedArena[Parameter](),
		tableFields:       NewTypedArena[TableField](),
		elseIfClauses:     NewTypedArena[ElseIfClause](),
		elseIfExprClauses: NewTypedArena[ElseIfExprClause](),
		tableTypeFields:   NewTypedArena[TableTypeField](),
	}
}

func (f *NodeFactory) Reset() {
	f.identifiers.Reset()
	f.literals.Reset()
	f.binaryOps.Reset()
	f.unaryOps.Reset()
	f.functionCalls.Reset()
	f.methodCalls.Reset()
	f.indexAccesses.Reset()
	f.fieldAccesses.Reset()
	f.tableLiterals.Reset()
	f.functionExprs.Reset()
	f.typeCasts.Reset()
	f.ifExprs.Reset()
	f.varArgs.Reset()
	f.parenExprs.Reset()
	f.interpolatedStrings.Reset()

	f.assignments.Reset()
	f.localAssignments.Reset()
	f.ifStatements.Reset()
	f.whileLoops.Reset()
	f.repeatLoops.Reset()
	f.forLoops.Reset()
	f.forInLoops.Reset()
	f.doBlocks.Reset()
	f.functionDefs.Reset()
	f.localFunctions.Reset()
	f.returnStatements.Reset()
	f.breakStatements.Reset()
	f.continueStatements.Reset()
	f.typeAliases.Reset()
	f.metamethodDefs.Reset()
	f.emptyStatements.Reset()
	f.expressionStatements.Reset()

	f.blocks.Reset()
	f.programs.Reset()
	f.comments.Reset()
	f.modules.Reset()
	f.attributes.Reset()

	f.primitiveTypes.Reset()
	f.unionTypes.Reset()
	f.optionalTypes.Reset()
	f.tableTypes.Reset()
	f.genericTypes.Reset()

	f.parameters.Reset()
	f.tableFields.Reset()
	f.elseIfClauses.Reset()
	f.elseIfExprClauses.Reset()
	f.tableTypeFields.Reset()
}

func (f *NodeFactory) Attribute(pos Position, name string) *Attribute {
	node := f.attributes.Alloc()
	node.Position = pos
	node.Name = name
	return node
}

func (f *NodeFactory) Identifier(pos Position, name string) *Identifier {
	node := f.identifiers.Alloc()
	node.Position = pos
	node.Name = name
	return node
}

func (f *NodeFactory) Literal(pos Position, litType string, value any) *Literal {
	node := f.literals.Alloc()
	node.Position = pos
	node.Type = litType
	node.Value = value
	return node
}

func (f *NodeFactory) BinaryOp(pos Position, left Expr, op string, right Expr) *BinaryOp {
	node := f.binaryOps.Alloc()
	node.Position = pos
	node.Left = left
	node.Op = op
	node.Right = right
	return node
}

func (f *NodeFactory) UnaryOp(pos Position, op string, operand Expr) *UnaryOp {
	node := f.unaryOps.Alloc()
	node.Position = pos
	node.Op = op
	node.Operand = operand
	return node
}

func (f *NodeFactory) FunctionCall(pos Position, function Expr, args []Expr) *FunctionCall {
	if args == nil {
		args = []Expr{}
	}
	node := f.functionCalls.Alloc()
	node.Position = pos
	node.Function = function
	node.Args = args
	return node
}

func (f *NodeFactory) MethodCall(pos Position, object Expr, method string, args []Expr) *MethodCall {
	if args == nil {
		args = []Expr{}
	}
	node := f.methodCalls.Alloc()
	node.Position = pos
	node.Object = object
	node.Method = method
	node.Args = args
	return node
}

func (f *NodeFactory) IndexAccess(pos Position, table Expr, index Expr) *IndexAccess {
	node := f.indexAccesses.Alloc()
	node.Position = pos
	node.Table = table
	node.Index = index
	return node
}

func (f *NodeFactory) FieldAccess(pos Position, object Expr, field string) *FieldAccess {
	node := f.fieldAccesses.Alloc()
	node.Position = pos
	node.Object = object
	node.Field = field
	return node
}

func (f *NodeFactory) TableLiteral(pos Position, fields []*TableField) *TableLiteral {
	if fields == nil {
		fields = []*TableField{}
	}
	node := f.tableLiterals.Alloc()
	node.Position = pos
	node.Fields = fields
	return node
}

func (f *NodeFactory) FunctionExpr(pos Position, generics []string, params []*Parameter, body *Block, returnType TypeNode) *FunctionExpr {
	if generics == nil {
		generics = []string{}
	}
	if params == nil {
		params = []*Parameter{}
	}
	node := f.functionExprs.Alloc()
	node.Position = pos
	node.Generics = generics
	node.Parameters = params
	node.Body = body
	node.ReturnType = returnType
	return node
}

func (f *NodeFactory) TypeCast(pos Position, value Expr, typeNode TypeNode) *TypeCast {
	node := f.typeCasts.Alloc()
	node.Position = pos
	node.Value = value
	node.Type = typeNode
	return node
}

func (f *NodeFactory) IfExpr(pos Position, condition Expr, then Expr, opts ...IfExprOption) *IfExpr {
	node := f.ifExprs.Alloc()
	node.Position = pos
	node.Condition = condition
	node.Then = then

	node.ElseIfs = []*ElseIfExprClause{}
	node.Else = nil

	for _, opt := range opts {
		opt(node)
	}

	return node
}

func (f *NodeFactory) VarArgs(pos Position) *VarArgs {
	node := f.varArgs.Alloc()
	node.Position = pos
	return node
}

func (f *NodeFactory) ParenExpr(pos Position, expr Expr) *ParenExpr {
	node := f.parenExprs.Alloc()
	node.Position = pos
	node.Expr = expr
	return node
}

func (f *NodeFactory) InterpolatedString(pos Position, segments []string, expressions []Expr) *InterpolatedString {
	if segments == nil {
		segments = []string{}
	}
	if expressions == nil {
		expressions = []Expr{}
	}
	node := f.interpolatedStrings.Alloc()
	node.Position = pos
	node.Segments = segments
	node.Expressions = expressions
	return node
}

func (f *NodeFactory) Block(pos Position, statements []Stmt) *Block {
	if statements == nil {
		statements = []Stmt{}
	}
	node := f.blocks.Alloc()
	node.Position = pos
	node.Statements = statements
	return node
}

func (f *NodeFactory) Assignment(pos Position, targets []Expr, operator string, values []Expr) *Assignment {
	if targets == nil {
		targets = []Expr{}
	}
	if values == nil {
		values = []Expr{}
	}
	node := f.assignments.Alloc()
	node.Position = pos
	node.Targets = targets
	node.Operator = operator
	node.Values = values
	return node
}

func (f *NodeFactory) LocalAssignment(pos Position, names []string, types []TypeNode, values []Expr) *LocalAssignment {
	node := f.localAssignments.Alloc()
	node.Position = pos
	node.Names = names
	node.Types = types
	node.Values = values
	return node
}

func (f *NodeFactory) IfStatement(pos Position, condition Expr, then *Block, opts ...IfStmtOption) *IfStatement {
	node := f.ifStatements.Alloc()

	node.Position = pos
	node.Condition = condition
	node.Then = then

	node.ElseIfs = []*ElseIfClause{}
	node.Else = nil

	for _, opt := range opts {
		opt(node)
	}

	return node
}

func (f *NodeFactory) WhileLoop(pos Position, condition Expr, body *Block) *WhileLoop {
	node := f.whileLoops.Alloc()
	node.Position = pos
	node.Condition = condition
	node.Body = body
	return node
}

func (f *NodeFactory) RepeatLoop(pos Position, body *Block, condition Expr) *RepeatLoop {
	node := f.repeatLoops.Alloc()
	node.Position = pos
	node.Body = body
	node.Condition = condition
	return node
}

func (f *NodeFactory) ForLoop(pos Position, variable string, start Expr, end Expr, body *Block, opts ...ForLoopOption) *ForLoop {
	node := f.forLoops.Alloc()
	node.Position = pos
	node.Variable = variable
	node.Start = start
	node.End = end
	node.Body = body
	node.Step = nil

	for _, opt := range opts {
		opt(node)
	}
	return node
}

func (f *NodeFactory) ForInLoop(pos Position, variables []string, iterables []Expr, body *Block) *ForInLoop {
	if variables == nil {
		variables = []string{}
	}
	if iterables == nil {
		iterables = []Expr{}
	}
	node := f.forInLoops.Alloc()
	node.Position = pos
	node.Variables = variables
	node.Iterables = iterables
	node.Body = body
	return node
}

func (f *NodeFactory) DoBlock(pos Position, body *Block) *DoBlock {
	node := f.doBlocks.Alloc()
	node.Position = pos
	node.Body = body
	return node
}

func (f *NodeFactory) FunctionDef(pos Position, name string, body *Block, opts ...FunctionDefOption) *FunctionDef {
	node := f.functionDefs.Alloc()

	node.Position = pos
	node.Name = name
	node.Body = body

	node.Generics = []string{}
	node.Parameters = []*Parameter{}
	node.ReturnType = nil
	node.Attributes = nil

	for _, opt := range opts {
		opt(node)
	}

	return node
}

func (f *NodeFactory) LocalFunction(pos Position, name string, generics []string, params []*Parameter, body *Block, returnType TypeNode) *LocalFunction {
	if generics == nil {
		generics = []string{}
	}
	if params == nil {
		params = []*Parameter{}
	}
	node := f.localFunctions.Alloc()
	node.Position = pos
	node.Name = name
	node.Generics = generics
	node.Parameters = params
	node.Body = body
	node.ReturnType = returnType
	return node
}

func (f *NodeFactory) ReturnStatement(pos Position, values []Expr) *ReturnStatement {
	if values == nil {
		values = []Expr{}
	}
	node := f.returnStatements.Alloc()
	node.Position = pos
	node.Values = values
	return node
}

func (f *NodeFactory) BreakStatement(pos Position) *BreakStatement {
	node := f.breakStatements.Alloc()
	node.Position = pos
	return node
}

func (f *NodeFactory) ContinueStatement(pos Position) *ContinueStatement {
	node := f.continueStatements.Alloc()
	node.Position = pos
	return node
}

func (f *NodeFactory) TypeAlias(pos Position, name string, typeNode TypeNode, opts ...TypeAliasOption) *TypeAlias {
	node := f.typeAliases.Alloc()

	node.Position = pos
	node.Name = name
	node.Type = typeNode

	node.Generics = []string{}
	node.IsExport = false

	for _, opt := range opts {
		opt(node)
	}

	return node
}

func (f *NodeFactory) MetamethodDef(pos Position, name string, params []*Parameter, body *Block) *MetamethodDef {
	if params == nil {
		params = []*Parameter{}
	}
	node := f.metamethodDefs.Alloc()
	node.Position = pos
	node.Name = name
	node.Parameters = params
	node.Body = body
	return node
}

func (f *NodeFactory) EmptyStatement(pos Position) *EmptyStatement {
	node := f.emptyStatements.Alloc()
	node.Position = pos
	return node
}

func (f *NodeFactory) ExpressionStatement(pos Position, expr Expr) *ExpressionStatement {
	node := f.expressionStatements.Alloc()
	node.Position = pos
	node.Expr = expr
	return node
}

func (f *NodeFactory) PrimitiveType(pos Position, name string) *PrimitiveType {
	node := f.primitiveTypes.Alloc()
	node.Position = pos
	node.Name = name
	return node
}

func (f *NodeFactory) UnionType(pos Position, left TypeNode, right TypeNode) *UnionType {
	node := f.unionTypes.Alloc()
	node.Position = pos
	node.Left = left
	node.Right = right
	return node
}

func (f *NodeFactory) OptionalType(pos Position, baseType TypeNode) *OptionalType {
	node := f.optionalTypes.Alloc()
	node.Position = pos
	node.BaseType = baseType
	return node
}

func (f *NodeFactory) TableType(pos Position, fields []*TableTypeField) *TableType {
	if fields == nil {
		fields = []*TableTypeField{}
	}
	node := f.tableTypes.Alloc()
	node.Position = pos
	node.Fields = fields
	return node
}

func (f *NodeFactory) GenericType(pos Position, baseType TypeNode, types []TypeNode) *GenericType {
	if types == nil {
		types = []TypeNode{}
	}
	node := f.genericTypes.Alloc()
	node.Position = pos
	node.BaseType = baseType
	node.Types = types
	return node
}

func (f *NodeFactory) Program(pos Position, body []Stmt) *Program {
	if body == nil {
		body = []Stmt{}
	}
	node := f.programs.Alloc()
	node.Position = pos
	node.Body = body
	return node
}

func (f *NodeFactory) Comment(pos Position, text string) *Comment {
	node := f.comments.Alloc()
	node.Position = pos
	node.Text = text
	return node
}

func (f *NodeFactory) Module(pos Position, name string, body *Block) *Module {
	node := f.modules.Alloc()
	node.Position = pos
	node.Name = name
	node.Body = body
	return node
}

func (f *NodeFactory) Parameter(name string, typeNode TypeNode) *Parameter {
	node := f.parameters.Alloc()
	node.Name = name
	node.Type = typeNode
	return node
}

func (f *NodeFactory) TableField(key Expr, value Expr) *TableField {
	node := f.tableFields.Alloc()
	node.Key = key
	node.Value = value
	return node
}

func (f *NodeFactory) ElseIfClause(condition Expr, body *Block) *ElseIfClause {
	node := f.elseIfClauses.Alloc()
	node.Condition = condition
	node.Body = body
	return node
}

func (f *NodeFactory) ElseIfExprClause(condition Expr, then Expr) *ElseIfExprClause {
	node := f.elseIfExprClauses.Alloc()
	node.Condition = condition
	node.Then = then
	return node
}

func (f *NodeFactory) TableTypeField(key TypeNode, keyName string, value TypeNode, isAccess bool) *TableTypeField {
	node := f.tableTypeFields.Alloc()
	node.Key = key
	node.KeyName = keyName
	node.Value = value
	node.IsAccess = isAccess
	return node
}

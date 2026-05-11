package ast

type EditFunc func(node Node, parents []Node) Node

func Rewrite(root Node, fn EditFunc) Node {
	return rewriteNode(root, nil, fn)
}

func rewriteNode(node Node, parents []Node, fn EditFunc) Node {
	if node == nil {
		return nil
	}

	next := append(parents, node)
	node = rewriteChildren(node, next, fn)

	if replacement := fn(node, parents); replacement != nil {
		return replacement
	}
	return node
}

func RewriteProgram(root *Program, fn EditFunc) *Program {
	result := rewriteNode(root, nil, fn)
	if p, ok := result.(*Program); ok {
		return p
	}
	return root
}

func rewriteExpr(e Expr, parents []Node, fn EditFunc) Expr {
	if e == nil {
		return nil
	}
	result := rewriteNode(e, parents, fn)
	if result == nil {
		return e
	}
	return result.(Expr)
}

func rewriteStmt(s Stmt, parents []Node, fn EditFunc) Stmt {
	if s == nil {
		return nil
	}
	result := rewriteNode(s, parents, fn)
	if result == nil {
		return s
	}
	return result.(Stmt)
}

func rewriteType(t TypeNode, parents []Node, fn EditFunc) TypeNode {
	if t == nil {
		return nil
	}
	result := rewriteNode(t, parents, fn)
	if result == nil {
		return t
	}
	return result.(TypeNode)
}

func rewriteBlock(b *Block, parents []Node, fn EditFunc) *Block {
	if b == nil {
		return nil
	}
	result := rewriteNode(b, parents, fn)
	if result == nil {
		return b
	}
	return result.(*Block)
}

func rewriteChildren(node Node, parents []Node, fn EditFunc) Node {
	switch n := node.(type) {

	case *Program:
		stmts, changed := rewriteStmtSlice(n.Body, parents, fn)
		if !changed {
			return node
		}
		cp := *n
		cp.Body = stmts
		return &cp

	case *Block:
		stmts, changed := rewriteStmtSlice(n.Statements, parents, fn)
		if !changed {
			return node
		}
		cp := *n
		cp.Statements = stmts
		return &cp

	case *Module:
		newBody := rewriteBlock(n.Body, parents, fn)
		if newBody == n.Body {
			return node
		}
		cp := *n
		cp.Body = newBody
		return &cp

	case *Identifier, *Literal, *VarArgs,
		*BreakStatement, *ContinueStatement, *EmptyStatement,
		*Comment, *Attribute, *PrimitiveType:
		return node

	case *BinaryOp:
		l, r := rewriteExpr(n.Left, parents, fn), rewriteExpr(n.Right, parents, fn)
		if l == n.Left && r == n.Right {
			return node
		}
		cp := *n
		cp.Left, cp.Right = l, r
		return &cp

	case *UnaryOp:
		op := rewriteExpr(n.Operand, parents, fn)
		if op == n.Operand {
			return node
		}
		cp := *n
		cp.Operand = op
		return &cp

	case *FunctionCall:
		fn2 := rewriteExpr(n.Function, parents, fn)
		args, changed := rewriteExprSlice(n.Args, parents, fn)
		if fn2 == n.Function && !changed {
			return node
		}
		cp := *n
		cp.Function = fn2
		cp.Args = args
		return &cp

	case *MethodCall:
		obj := rewriteExpr(n.Object, parents, fn)
		args, changed := rewriteExprSlice(n.Args, parents, fn)
		if obj == n.Object && !changed {
			return node
		}
		cp := *n
		cp.Object = obj
		cp.Args = args
		return &cp

	case *IndexAccess:
		tbl, idx := rewriteExpr(n.Table, parents, fn), rewriteExpr(n.Index, parents, fn)
		if tbl == n.Table && idx == n.Index {
			return node
		}
		cp := *n
		cp.Table, cp.Index = tbl, idx
		return &cp

	case *FieldAccess:
		obj := rewriteExpr(n.Object, parents, fn)
		if obj == n.Object {
			return node
		}
		cp := *n
		cp.Object = obj
		return &cp

	case *TableLiteral:
		fields, changed := rewriteTableFields(n.Fields, parents, fn)
		if !changed {
			return node
		}
		cp := *n
		cp.Fields = fields
		return &cp

	case *FunctionExpr:
		params, pChanged := rewriteParams(n.Parameters, parents, fn)
		body := rewriteBlock(n.Body, parents, fn)
		ret := rewriteType(n.ReturnType, parents, fn)
		if !pChanged && body == n.Body && ret == n.ReturnType {
			return node
		}
		cp := *n
		cp.Parameters = params
		cp.Body = body
		cp.ReturnType = ret
		return &cp

	case *TypeCast:
		val := rewriteExpr(n.Value, parents, fn)
		typ := rewriteType(n.Type, parents, fn)
		if val == n.Value && typ == n.Type {
			return node
		}
		cp := *n
		cp.Value = val
		cp.Type = typ
		return &cp

	case *IfExpr:
		cond := rewriteExpr(n.Condition, parents, fn)
		then := rewriteExpr(n.Then, parents, fn)
		elseIfs, eiChanged := rewriteElseIfExprs(n.ElseIfs, parents, fn)
		els := rewriteExpr(n.Else, parents, fn)
		if cond == n.Condition && then == n.Then && !eiChanged && els == n.Else {
			return node
		}
		cp := *n
		cp.Condition, cp.Then, cp.ElseIfs, cp.Else = cond, then, elseIfs, els
		return &cp

	case *ParenExpr:
		e := rewriteExpr(n.Expr, parents, fn)
		if e == n.Expr {
			return node
		}
		cp := *n
		cp.Expr = e
		return &cp

	case *InterpolatedString:
		exprs, changed := rewriteExprSlice(n.Expressions, parents, fn)
		if !changed {
			return node
		}
		cp := *n
		cp.Expressions = exprs
		return &cp

	case *Assignment:
		targets, tChanged := rewriteExprSlice(n.Targets, parents, fn)
		values, vChanged := rewriteExprSlice(n.Values, parents, fn)
		if !tChanged && !vChanged {
			return node
		}
		cp := *n
		cp.Targets, cp.Values = targets, values
		return &cp

	case *LocalAssignment:
		values, vChanged := rewriteExprSlice(n.Values, parents, fn)
		types, tChanged := rewriteTypeSlice(n.Types, parents, fn)
		if !vChanged && !tChanged {
			return node
		}
		cp := *n
		cp.Values, cp.Types = values, types
		return &cp

	case *IfStatement:
		cond := rewriteExpr(n.Condition, parents, fn)
		then := rewriteBlock(n.Then, parents, fn)
		elseIfs, eiChanged := rewriteElseIfStmts(n.ElseIfs, parents, fn)
		els := rewriteBlock(n.Else, parents, fn)
		if cond == n.Condition && then == n.Then && !eiChanged && els == n.Else {
			return node
		}
		cp := *n
		cp.Condition, cp.Then, cp.ElseIfs, cp.Else = cond, then, elseIfs, els
		return &cp

	case *WhileLoop:
		cond := rewriteExpr(n.Condition, parents, fn)
		body := rewriteBlock(n.Body, parents, fn)
		if cond == n.Condition && body == n.Body {
			return node
		}
		cp := *n
		cp.Condition, cp.Body = cond, body
		return &cp

	case *RepeatLoop:
		body := rewriteBlock(n.Body, parents, fn)
		cond := rewriteExpr(n.Condition, parents, fn)
		if body == n.Body && cond == n.Condition {
			return node
		}
		cp := *n
		cp.Body, cp.Condition = body, cond
		return &cp

	case *ForLoop:
		start := rewriteExpr(n.Start, parents, fn)
		end := rewriteExpr(n.End, parents, fn)
		step := rewriteExpr(n.Step, parents, fn)
		body := rewriteBlock(n.Body, parents, fn)
		if start == n.Start && end == n.End && step == n.Step && body == n.Body {
			return node
		}
		cp := *n
		cp.Start, cp.End, cp.Step, cp.Body = start, end, step, body
		return &cp

	case *ForInLoop:
		iters, changed := rewriteExprSlice(n.Iterables, parents, fn)
		body := rewriteBlock(n.Body, parents, fn)
		if !changed && body == n.Body {
			return node
		}
		cp := *n
		cp.Iterables, cp.Body = iters, body
		return &cp

	case *DoBlock:
		body := rewriteBlock(n.Body, parents, fn)
		if body == n.Body {
			return node
		}
		cp := *n
		cp.Body = body
		return &cp

	case *FunctionDef:
		params, pChanged := rewriteParams(n.Parameters, parents, fn)
		body := rewriteBlock(n.Body, parents, fn)
		ret := rewriteType(n.ReturnType, parents, fn)
		if !pChanged && body == n.Body && ret == n.ReturnType {
			return node
		}
		cp := *n
		cp.Parameters, cp.Body, cp.ReturnType = params, body, ret
		return &cp

	case *LocalFunction:
		params, pChanged := rewriteParams(n.Parameters, parents, fn)
		body := rewriteBlock(n.Body, parents, fn)
		ret := rewriteType(n.ReturnType, parents, fn)
		if !pChanged && body == n.Body && ret == n.ReturnType {
			return node
		}
		cp := *n
		cp.Parameters, cp.Body, cp.ReturnType = params, body, ret
		return &cp

	case *ReturnStatement:
		values, changed := rewriteExprSlice(n.Values, parents, fn)
		if !changed {
			return node
		}
		cp := *n
		cp.Values = values
		return &cp

	case *TypeAlias:
		typ := rewriteType(n.Type, parents, fn)
		if typ == n.Type {
			return node
		}
		cp := *n
		cp.Type = typ
		return &cp

	case *MetamethodDef:
		params, pChanged := rewriteParams(n.Parameters, parents, fn)
		body := rewriteBlock(n.Body, parents, fn)
		if !pChanged && body == n.Body {
			return node
		}
		cp := *n
		cp.Parameters, cp.Body = params, body
		return &cp

	case *ExpressionStatement:
		e := rewriteExpr(n.Expr, parents, fn)
		if e == n.Expr {
			return node
		}
		cp := *n
		cp.Expr = e
		return &cp

	case *UnionType:
		l, r := rewriteType(n.Left, parents, fn), rewriteType(n.Right, parents, fn)
		if l == n.Left && r == n.Right {
			return node
		}
		cp := *n
		cp.Left, cp.Right = l, r
		return &cp

	case *OptionalType:
		base := rewriteType(n.BaseType, parents, fn)
		if base == n.BaseType {
			return node
		}
		cp := *n
		cp.BaseType = base
		return &cp

	case *TableType:
		fields, changed := rewriteTableTypeFields(n.Fields, parents, fn)
		if !changed {
			return node
		}
		cp := *n
		cp.Fields = fields
		return &cp

	case *GenericType:
		base := rewriteType(n.BaseType, parents, fn)
		types, changed := rewriteTypeSlice(n.Types, parents, fn)
		if base == n.BaseType && !changed {
			return node
		}
		cp := *n
		cp.BaseType, cp.Types = base, types
		return &cp
	}

	return node
}

func rewriteExprSlice(exprs []Expr, parents []Node, fn EditFunc) ([]Expr, bool) {
	changed := false
	out := make([]Expr, len(exprs))
	for i, e := range exprs {
		ne := rewriteExpr(e, parents, fn)
		out[i] = ne
		if ne != e {
			changed = true
		}
	}
	if !changed {
		return exprs, false
	}
	return out, true
}

func rewriteStmtSlice(stmts []Stmt, parents []Node, fn EditFunc) ([]Stmt, bool) {
	changed := false
	out := make([]Stmt, len(stmts))
	for i, s := range stmts {
		ns := rewriteStmt(s, parents, fn)
		out[i] = ns
		if ns != s {
			changed = true
		}
	}
	if !changed {
		return stmts, false
	}
	return out, true
}

func rewriteTypeSlice(types []TypeNode, parents []Node, fn EditFunc) ([]TypeNode, bool) {
	changed := false
	out := make([]TypeNode, len(types))
	for i, t := range types {
		nt := rewriteType(t, parents, fn)
		out[i] = nt
		if nt != t {
			changed = true
		}
	}
	if !changed {
		return types, false
	}
	return out, true
}

func rewriteParams(params []*Parameter, parents []Node, fn EditFunc) ([]*Parameter, bool) {
	changed := false
	out := make([]*Parameter, len(params))
	for i, p := range params {
		newType := rewriteType(p.Type, parents, fn)
		if newType == p.Type {
			out[i] = p
		} else {
			cp := *p
			cp.Type = newType
			out[i] = &cp
			changed = true
		}
	}
	if !changed {
		return params, false
	}
	return out, true
}

func rewriteTableFields(fields []*TableField, parents []Node, fn EditFunc) ([]*TableField, bool) {
	changed := false
	out := make([]*TableField, len(fields))
	for i, f := range fields {
		nk := rewriteExpr(f.Key, parents, fn)
		nv := rewriteExpr(f.Value, parents, fn)
		if nk == f.Key && nv == f.Value {
			out[i] = f
		} else {
			cp := *f
			cp.Key, cp.Value = nk, nv
			out[i] = &cp
			changed = true
		}
	}
	if !changed {
		return fields, false
	}
	return out, true
}

func rewriteTableTypeFields(fields []*TableTypeField, parents []Node, fn EditFunc) ([]*TableTypeField, bool) {
	changed := false
	out := make([]*TableTypeField, len(fields))
	for i, f := range fields {
		nk := rewriteType(f.Key, parents, fn)
		nv := rewriteType(f.Value, parents, fn)
		if nk == f.Key && nv == f.Value {
			out[i] = f
		} else {
			cp := *f
			cp.Key, cp.Value = nk, nv
			out[i] = &cp
			changed = true
		}
	}
	if !changed {
		return fields, false
	}
	return out, true
}

func rewriteElseIfExprs(clauses []*ElseIfExprClause, parents []Node, fn EditFunc) ([]*ElseIfExprClause, bool) {
	changed := false
	out := make([]*ElseIfExprClause, len(clauses))
	for i, c := range clauses {
		nc := rewriteExpr(c.Condition, parents, fn)
		nt := rewriteExpr(c.Then, parents, fn)
		if nc == c.Condition && nt == c.Then {
			out[i] = c
		} else {
			cp := *c
			cp.Condition, cp.Then = nc, nt
			out[i] = &cp
			changed = true
		}
	}
	if !changed {
		return clauses, false
	}
	return out, true
}

func rewriteElseIfStmts(clauses []*ElseIfClause, parents []Node, fn EditFunc) ([]*ElseIfClause, bool) {
	changed := false
	out := make([]*ElseIfClause, len(clauses))
	for i, c := range clauses {
		nc := rewriteExpr(c.Condition, parents, fn)
		nb := rewriteBlock(c.Body, parents, fn)
		if nc == c.Condition && nb == c.Body {
			out[i] = c
		} else {
			cp := *c
			cp.Condition, cp.Body = nc, nb
			out[i] = &cp
			changed = true
		}
	}
	if !changed {
		return clauses, false
	}
	return out, true
}

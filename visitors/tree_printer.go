package visitors

import (
	"fmt"
	"strings"

	"github.com/Wh1teSlash/luau-parser/ast"
)

type TreePrinter struct {
	builder strings.Builder
	indent  int
}

func NewTreePrinter() *TreePrinter {
	return &TreePrinter{}
}

func (p *TreePrinter) Print(node ast.Node) string {
	p.builder.Reset()
	p.indent = 0
	node.Accept(p)
	return p.builder.String()
}

// Helper to write an indented line with formatting
func (p *TreePrinter) writeLine(format string, args ...any) {
	p.builder.WriteString(strings.Repeat("  ", p.indent))
	p.builder.WriteString(fmt.Sprintf(format, args...))
	p.builder.WriteString("\n")
}

func (p *TreePrinter) printParams(params []*ast.Parameter) {
	if len(params) == 0 {
		return
	}
	p.writeLine("Parameters:")
	p.indent++
	for _, param := range params {
		typeStr := ""
		if param.Type != nil {
			typeStr = fmt.Sprintf(" (Type: %s)", param.Type.Type)
		}
		p.writeLine("- %s%s", param.Name, typeStr)
	}
	p.indent--
}

func (p *TreePrinter) VisitProgram(node *ast.Program) any {
	p.writeLine("Program")
	p.indent++
	for _, stmt := range node.Body {
		stmt.Accept(p)
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitBlock(node *ast.Block) any {
	p.writeLine("Block")
	p.indent++
	for _, stmt := range node.Statements {
		stmt.Accept(p)
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitModule(node *ast.Module) any {
	p.writeLine("Module")
	p.indent++
	if node.Body != nil {
		node.Body.Accept(p)
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitComment(node *ast.Comment) any {
	p.writeLine("Comment: %q", strings.TrimSpace(node.Text))
	return nil
}

func (p *TreePrinter) VisitAssignment(node *ast.Assignment) any {
	p.writeLine("Assignment (Op: %s)", node.Operator)
	p.indent++
	p.writeLine("Targets:")
	p.indent++
	for _, target := range node.Targets {
		target.Accept(p)
	}
	p.indent--
	p.writeLine("Values:")
	p.indent++
	for _, value := range node.Values {
		value.Accept(p)
	}
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitLocalAssignment(node *ast.LocalAssignment) any {
	p.writeLine("LocalAssignment (Names: %s)", strings.Join(node.Names, ", "))
	p.indent++
	if len(node.Values) > 0 {
		p.writeLine("Values:")
		p.indent++
		for _, value := range node.Values {
			value.Accept(p)
		}
		p.indent--
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitIfStatement(node *ast.IfStatement) any {
	p.writeLine("IfStatement")
	p.indent++
	p.writeLine("Condition:")
	p.indent++
	node.Condition.Accept(p)
	p.indent--
	p.writeLine("Then:")
	p.indent++
	node.Then.Accept(p)
	p.indent--

	for _, elif := range node.ElseIfs {
		p.writeLine("ElseIf:")
		p.indent++
		p.writeLine("Condition:")
		p.indent++
		elif.Condition.Accept(p)
		p.indent--
		p.writeLine("Then:")
		p.indent++
		elif.Body.Accept(p)
		p.indent--
		p.indent--
	}

	if node.Else != nil {
		p.writeLine("Else:")
		p.indent++
		node.Else.Accept(p)
		p.indent--
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitWhileLoop(node *ast.WhileLoop) any {
	p.writeLine("WhileLoop")
	p.indent++
	p.writeLine("Condition:")
	p.indent++
	node.Condition.Accept(p)
	p.indent--
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitRepeatLoop(node *ast.RepeatLoop) any {
	p.writeLine("RepeatLoop")
	p.indent++
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.writeLine("Condition:")
	p.indent++
	node.Condition.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitForLoop(node *ast.ForLoop) any {
	p.writeLine("ForLoop (Var: %s)", node.Variable)
	p.indent++
	p.writeLine("Start:")
	p.indent++
	node.Start.Accept(p)
	p.indent--
	p.writeLine("End:")
	p.indent++
	node.End.Accept(p)
	p.indent--
	if node.Step != nil {
		p.writeLine("Step:")
		p.indent++
		node.Step.Accept(p)
		p.indent--
	}
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitForInLoop(node *ast.ForInLoop) any {
	p.writeLine("ForInLoop (Vars: %s)", strings.Join(node.Variables, ", "))
	p.indent++
	p.writeLine("Iterables:")
	p.indent++
	for _, it := range node.Iterables {
		it.Accept(p)
	}
	p.indent--
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitDoBlock(node *ast.DoBlock) any {
	p.writeLine("DoBlock")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	return nil
}

func (p *TreePrinter) VisitFunctionDef(node *ast.FunctionDef) any {
	retType := ""
	if node.ReturnType != nil {
		retType = " -> " + node.ReturnType.Type
	}
	p.writeLine("FunctionDef (Name: %s)%s", node.Name, retType)
	p.indent++
	p.printParams(node.Parameters)
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitLocalFunction(node *ast.LocalFunction) any {
	retType := ""
	if node.ReturnType != nil {
		retType = " -> " + node.ReturnType.Type
	}
	p.writeLine("LocalFunctionDef (Name: %s)%s", node.Name, retType)
	p.indent++
	p.printParams(node.Parameters)
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitReturnStatement(node *ast.ReturnStatement) any {
	p.writeLine("ReturnStatement")
	if len(node.Values) > 0 {
		p.indent++
		for _, val := range node.Values {
			val.Accept(p)
		}
		p.indent--
	}
	return nil
}

func (p *TreePrinter) VisitBreakStatement(node *ast.BreakStatement) any {
	p.writeLine("BreakStatement")
	return nil
}

func (p *TreePrinter) VisitContinueStatement(node *ast.ContinueStatement) any {
	p.writeLine("ContinueStatement")
	return nil
}

func (p *TreePrinter) VisitTypeAlias(node *ast.TypeAlias) any {
	p.writeLine("TypeAlias (Export: %t, Name: %s, Type: %s)", node.IsExport, node.Name, node.Type.Type)
	return nil
}

func (p *TreePrinter) VisitMetamethodDef(node *ast.MetamethodDef) any {
	p.writeLine("MetamethodDef (Name: %s)", node.Name)
	p.indent++
	p.printParams(node.Parameters)
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitEmptyStatement(node *ast.EmptyStatement) any {
	p.writeLine("EmptyStatement")
	return nil
}

func (p *TreePrinter) VisitBinaryOp(node *ast.BinaryOp) any {
	p.writeLine("BinaryOp (Op: %s)", node.Op)
	p.indent++
	node.Left.Accept(p)
	node.Right.Accept(p)
	p.indent--
	return nil
}

func (p *TreePrinter) VisitUnaryOp(node *ast.UnaryOp) any {
	p.writeLine("UnaryOp (Op: %s)", node.Op)
	p.indent++
	node.Operand.Accept(p)
	p.indent--
	return nil
}

func (p *TreePrinter) VisitIdentifier(node *ast.Identifier) any {
	p.writeLine("Identifier: %s", node.Name)
	return nil
}

func (p *TreePrinter) VisitLiteral(node *ast.Literal) any {
	if node.Type == "string" {
		p.writeLine("Literal (Type: %s, Value: %q)", node.Type, node.Value)
	} else {
		p.writeLine("Literal (Type: %s, Value: %v)", node.Type, node.Value)
	}
	return nil
}

func (p *TreePrinter) VisitFunctionCall(node *ast.FunctionCall) any {
	p.writeLine("FunctionCall")
	p.indent++
	p.writeLine("Function:")
	p.indent++
	node.Function.Accept(p)
	p.indent--
	if len(node.Args) > 0 {
		p.writeLine("Arguments:")
		p.indent++
		for _, arg := range node.Args {
			arg.Accept(p)
		}
		p.indent--
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitMethodCall(node *ast.MethodCall) any {
	p.writeLine("MethodCall (Method: %s)", node.Method)
	p.indent++
	p.writeLine("Object:")
	p.indent++
	node.Object.Accept(p)
	p.indent--
	if len(node.Args) > 0 {
		p.writeLine("Arguments:")
		p.indent++
		for _, arg := range node.Args {
			arg.Accept(p)
		}
		p.indent--
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitIndexAccess(node *ast.IndexAccess) any {
	p.writeLine("IndexAccess")
	p.indent++
	p.writeLine("Table:")
	p.indent++
	node.Table.Accept(p)
	p.indent--
	p.writeLine("Index:")
	p.indent++
	node.Index.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitFieldAccess(node *ast.FieldAccess) any {
	p.writeLine("FieldAccess (Field: %s)", node.Field)
	p.indent++
	p.writeLine("Object:")
	p.indent++
	node.Object.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitTableLiteral(node *ast.TableLiteral) any {
	p.writeLine("TableLiteral")
	if len(node.Fields) > 0 {
		p.indent++
		for _, field := range node.Fields {
			if field.Key != nil {
				p.writeLine("Field (Key/Value):")
				p.indent++
				field.Key.Accept(p)
				field.Value.Accept(p)
				p.indent--
			} else {
				p.writeLine("Field (Value only):")
				p.indent++
				field.Value.Accept(p)
				p.indent--
			}
		}
		p.indent--
	}
	return nil
}

func (p *TreePrinter) VisitFunctionExpr(node *ast.FunctionExpr) any {
	retType := ""
	if node.ReturnType != nil {
		retType = " -> " + node.ReturnType.Type
	}
	p.writeLine("FunctionExpr%s", retType)
	p.indent++
	p.printParams(node.Parameters)
	p.writeLine("Body:")
	p.indent++
	node.Body.Accept(p)
	p.indent--
	p.indent--
	return nil
}

func (p *TreePrinter) VisitTypeCast(node *ast.TypeCast) any {
	p.writeLine("TypeCast (To: %s)", node.Type.Type)
	p.indent++
	node.Value.Accept(p)
	p.indent--
	return nil
}

func (p *TreePrinter) VisitIfExpr(node *ast.IfExpr) any {
	p.writeLine("IfExpr")
	p.indent++
	p.writeLine("Condition:")
	p.indent++
	node.Condition.Accept(p)
	p.indent--
	p.writeLine("Then:")
	p.indent++
	node.Then.Accept(p)
	p.indent--
	if node.Else != nil {
		p.writeLine("Else:")
		p.indent++
		node.Else.Accept(p)
		p.indent--
	}
	p.indent--
	return nil
}

func (p *TreePrinter) VisitVarArgs(node *ast.VarArgs) any {
	p.writeLine("VarArgs (...)")
	return nil
}

func (p *TreePrinter) VisitParenExpr(node *ast.ParenExpr) any {
	p.writeLine("ParenExpr")
	p.indent++
	node.Expr.Accept(p)
	p.indent--
	return nil
}

func (p *TreePrinter) VisitExpressionStatement(node *ast.ExpressionStatement) any {
	p.writeLine("ExpressionStatement")
	p.indent++
	node.Expr.Accept(p)
	p.indent--
	return nil
}

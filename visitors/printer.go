package visitors

import (
	"fmt"
	"strings"

	"github.com/Wh1teSlash/luau-parser/ast"
)

type Printer struct {
	builder strings.Builder
	indent  int
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(node ast.Node) string {
	p.builder.Reset()
	p.indent = 0
	node.Accept(p)
	return p.builder.String()
}

func (p *Printer) write(s string) {
	p.builder.WriteString(s)
}

func (p *Printer) writeIndent() {
	p.builder.WriteString(strings.Repeat("\t", p.indent))
}

func (p *Printer) printExprList(exprs []ast.Expr) {
	for i, expr := range exprs {
		expr.Accept(p)
		if i < len(exprs)-1 {
			p.write(", ")
		}
	}
}

func (p *Printer) printParams(params []*ast.Parameter) {
	for i, param := range params {
		p.write(param.Name)
		if param.Type != nil {
			p.write(": ")
			p.write(param.Type.Type)
		}
		if i < len(params)-1 {
			p.write(", ")
		}
	}
}

func (p *Printer) VisitProgram(node *ast.Program) any {
	for _, stmt := range node.Body {
		stmt.Accept(p)
		p.write("\n")
	}
	return nil
}

func (p *Printer) VisitBlock(node *ast.Block) any {
	for _, stmt := range node.Statements {
		stmt.Accept(p)
		p.write("\n")
	}
	return nil
}

func (p *Printer) VisitModule(node *ast.Module) any {
	if node.Body != nil {
		node.Body.Accept(p)
	}
	return nil
}

func (p *Printer) VisitComment(node *ast.Comment) any {
	p.writeIndent()
	p.write("-- ")
	p.write(node.Text)
	return nil
}

func (p *Printer) VisitAssignment(node *ast.Assignment) any {
	p.writeIndent()
	node.Target.Accept(p)
	p.write(" = ")
	node.Value.Accept(p)
	return nil
}

func (p *Printer) VisitLocalAssignment(node *ast.LocalAssignment) any {
	p.writeIndent()
	p.write("local ")
	for i, name := range node.Names {
		p.write(name)
		if i < len(node.Types) && node.Types[i] != nil {
			p.write(": ")
			p.write(node.Types[i].Type)
		}
		if i < len(node.Names)-1 {
			p.write(", ")
		}
	}

	if len(node.Values) > 0 {
		p.write(" = ")
		p.printExprList(node.Values)
	}
	return nil
}

func (p *Printer) VisitIfStatement(node *ast.IfStatement) any {
	p.writeIndent()
	p.write("if ")
	node.Condition.Accept(p)
	p.write(" then\n")

	p.indent++
	node.Then.Accept(p)
	p.indent--

	for _, elif := range node.ElseIfs {
		p.writeIndent()
		p.write("elseif ")
		elif.Condition.Accept(p)
		p.write(" then\n")
		p.indent++
		elif.Body.Accept(p)
		p.indent--
	}

	if node.Else != nil {
		p.writeIndent()
		p.write("else\n")
		p.indent++
		node.Else.Accept(p)
		p.indent--
	}

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitWhileLoop(node *ast.WhileLoop) any {
	p.writeIndent()
	p.write("while ")
	node.Condition.Accept(p)
	p.write(" do\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitRepeatLoop(node *ast.RepeatLoop) any {
	p.writeIndent()
	p.write("repeat\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("until ")
	node.Condition.Accept(p)
	return nil
}

func (p *Printer) VisitForLoop(node *ast.ForLoop) any {
	p.writeIndent()
	p.write("for ")
	p.write(node.Variable)
	p.write(" = ")
	node.Start.Accept(p)
	p.write(", ")
	node.End.Accept(p)
	if node.Step != nil {
		p.write(", ")
		node.Step.Accept(p)
	}
	p.write(" do\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitForInLoop(node *ast.ForInLoop) any {
	p.writeIndent()
	p.write("for ")
	p.write(strings.Join(node.Variables, ", "))
	p.write(" in ")
	p.printExprList(node.Iterables)
	p.write(" do\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitDoBlock(node *ast.DoBlock) any {
	p.writeIndent()
	p.write("do\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitFunctionDef(node *ast.FunctionDef) any {
	p.writeIndent()
	p.write(fmt.Sprintf("function %s(", node.Name))
	p.printParams(node.Parameters)
	p.write(")")

	if node.ReturnType != nil {
		p.write(": ")
		p.write(node.ReturnType.Type)
	}
	p.write("\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitLocalFunction(node *ast.LocalFunction) any {
	p.writeIndent()
	p.write(fmt.Sprintf("local function %s(", node.Name))
	p.printParams(node.Parameters)
	p.write(")")

	if node.ReturnType != nil {
		p.write(": ")
		p.write(node.ReturnType.Type)
	}
	p.write("\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitReturnStatement(node *ast.ReturnStatement) any {
	p.writeIndent()
	p.write("return")
	if len(node.Values) > 0 {
		p.write(" ")
		p.printExprList(node.Values)
	}
	return nil
}

func (p *Printer) VisitBreakStatement(node *ast.BreakStatement) any {
	p.writeIndent()
	p.write("break")
	return nil
}

func (p *Printer) VisitContinueStatement(node *ast.ContinueStatement) any {
	p.writeIndent()
	p.write("continue")
	return nil
}

func (p *Printer) VisitTypeAlias(node *ast.TypeAlias) any {
	p.writeIndent()
	if node.IsExport {
		p.write("export ")
	}
	p.write(fmt.Sprintf("type %s = ", node.Name))
	p.write(node.Type.Type)
	return nil
}

func (p *Printer) VisitMetamethodDef(node *ast.MetamethodDef) any {
	p.writeIndent()

	p.write(fmt.Sprintf("function %s(", node.Name))
	p.printParams(node.Parameters)
	p.write(")\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitEmptyStatement(node *ast.EmptyStatement) any {
	p.writeIndent()
	p.write(";")
	return nil
}

func (p *Printer) VisitBinaryOp(node *ast.BinaryOp) any {
	node.Left.Accept(p)
	p.write(" " + node.Op + " ")
	node.Right.Accept(p)
	return nil
}

func (p *Printer) VisitUnaryOp(node *ast.UnaryOp) any {
	p.write(node.Op)
	if node.Op == "not" {
		p.write(" ")
	}
	node.Operand.Accept(p)
	return nil
}

func (p *Printer) VisitIdentifier(node *ast.Identifier) any {
	p.write(node.Name)
	return nil
}

func (p *Printer) VisitLiteral(node *ast.Literal) any {
	if node.Type == "string" {
		p.write(fmt.Sprintf("%q", node.Value))
	} else {
		p.write(fmt.Sprintf("%v", node.Value))
	}
	return nil
}

func (p *Printer) VisitFunctionCall(node *ast.FunctionCall) any {
	node.Function.Accept(p)
	p.write("(")
	p.printExprList(node.Args)
	p.write(")")
	return nil
}

func (p *Printer) VisitMethodCall(node *ast.MethodCall) any {
	node.Object.Accept(p)
	p.write(":")
	p.write(node.Method)
	p.write("(")
	p.printExprList(node.Args)
	p.write(")")
	return nil
}

func (p *Printer) VisitIndexAccess(node *ast.IndexAccess) any {
	node.Table.Accept(p)
	p.write("[")
	node.Index.Accept(p)
	p.write("]")
	return nil
}

func (p *Printer) VisitFieldAccess(node *ast.FieldAccess) any {
	node.Object.Accept(p)
	p.write(".")
	p.write(node.Field)
	return nil
}

func (p *Printer) VisitTableLiteral(node *ast.TableLiteral) any {
	p.write("{")
	for i, field := range node.Fields {
		if field.Key != nil {
			p.write("[")
			field.Key.Accept(p)
			p.write("] = ")
		}
		field.Value.Accept(p)

		if i < len(node.Fields)-1 {
			p.write(", ")
		}
	}
	p.write("}")
	return nil
}

func (p *Printer) VisitFunctionExpr(node *ast.FunctionExpr) any {
	p.write("function(")
	p.printParams(node.Parameters)
	p.write(")")

	if node.ReturnType != nil {
		p.write(": ")
		p.write(node.ReturnType.Type)
	}
	p.write("\n")

	p.indent++
	node.Body.Accept(p)
	p.indent--

	p.writeIndent()
	p.write("end")
	return nil
}

func (p *Printer) VisitTypeCast(node *ast.TypeCast) any {
	node.Value.Accept(p)
	p.write(" :: ")
	p.write(node.Type.Type)
	return nil
}

func (p *Printer) VisitIfExpr(node *ast.IfExpr) any {
	p.write("if ")
	node.Condition.Accept(p)
	p.write(" then ")
	node.Then.Accept(p)
	if node.Else != nil {
		p.write(" else ")
		node.Else.Accept(p)
	}
	return nil
}

func (p *Printer) VisitVarArgs(node *ast.VarArgs) any {
	p.write("...")
	return nil
}

func (p *Printer) VisitParenExpr(node *ast.ParenExpr) any {
	p.write("(")
	node.Expr.Accept(p)
	p.write(")")
	return nil
}

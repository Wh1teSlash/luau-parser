package ast

import "fmt"

type Block struct {
	BaseNode
	Statements []Stmt
}

func (b *Block) String() string       { return fmt.Sprintf("Block{%d statements}", len(b.Statements)) }
func (b *Block) Accept(v Visitor) any { return v.VisitBlock(b) }
func (b *Block) statementNode()       {}

type Program struct {
	BaseNode
	Body []Stmt
}

func (p *Program) String() string       { return fmt.Sprintf("Program{%d statements}", len(p.Body)) }
func (p *Program) Accept(v Visitor) any { return v.VisitProgram(p) }

type Comment struct {
	BaseNode
	Text string
}

func (c *Comment) String() string       { return fmt.Sprintf("Comment{%s}", c.Text) }
func (c *Comment) Accept(v Visitor) any { return v.VisitComment(c) }

type Module struct {
	BaseNode
	Name string
	Body *Block
}

func (m *Module) String() string       { return fmt.Sprintf("Module{name: %s}", m.Name) }
func (m *Module) Accept(v Visitor) any { return v.VisitModule(m) }

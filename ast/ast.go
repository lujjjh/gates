package ast

import "github.com/lujjjh/gates/token"

type Expr interface {
	exprNode()
}

type Ident struct {
	NamePos token.Pos
	Name    string
}

type Lit struct {
	ValuePos token.Pos
	Kind     token.Token
	Value    string
}

type UnaryExpr struct {
	OpPos token.Pos
	Op    token.Token
	X     Expr
}

type BinaryExpr struct {
	X     Expr
	OpPos token.Pos
	Op    token.Token
	Y     Expr
}

type ParenExpr struct {
	Lparen token.Pos
	X      Expr
	Rparen token.Pos
}

type BadExpr struct {
	From, To token.Pos
}

func (*Ident) exprNode()      {}
func (*Lit) exprNode()        {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*ParenExpr) exprNode()  {}
func (*BadExpr) exprNode()    {}

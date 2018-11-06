package syntax

type Expr interface {
	exprNode()
}

type expr struct{}

func (expr) exprNode() {}

type (
	Ident struct {
		expr
		NamePos Pos
		Name    string
	}

	Lit struct {
		expr
		ValuePos Pos
		Kind     Token
		Value    string
	}

	UnaryExpr struct {
		expr
		OpPos Pos
		Op    Token
		X     Expr
	}

	BinaryExpr struct {
		expr
		X     Expr
		OpPos Pos
		Op    Token
		Y     Expr
	}

	ParenExpr struct {
		expr
		Lparen Pos
		X      Expr
		Rparen Pos
	}

	BadExpr struct {
		expr
		From, To Pos
	}
)

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

	ArrayLit struct {
		expr
		Lbrack   Pos
		ElemList []Expr
		Rbrack   Pos
	}

	MapLitEntry struct {
		Key   Expr
		Value Expr
	}

	MapLit struct {
		expr
		Lbrace  Pos
		Entries []MapLitEntry
		Rbrace  Pos
	}

	FunctionLit struct {
		expr
		Function      Pos
		ParameterList *ParameterList
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

	SelectorExpr struct {
		expr
		X   Expr
		Sel *Ident
	}

	IndexExpr struct {
		expr
		X      Expr
		Lbrack Pos
		Index  Expr
		Rbrack Pos
	}

	CallExpr struct {
		expr
		Fun    Expr
		Lparen Pos
		Args   []Expr
		Rparen Pos
	}

	BadExpr struct {
		expr
		From, To Pos
	}

	ParameterList struct {
		Lparen Pos
		List   []*Ident
		Rparen Pos
	}
)

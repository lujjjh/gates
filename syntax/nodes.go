package syntax

type Expr interface {
	exprNode()
}

type expr struct{}

func (expr) exprNode() {}

type Stmt interface {
	stmtNode()
}

type stmt struct{}

func (stmt) stmtNode() {}

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
		Body          *FunctionBody
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

	AssignStmt struct {
		stmt
		Lhs    Expr
		TokPos Pos
		Tok    Token
		Rhs    Expr
	}

	ExprStmt struct {
		stmt
		X Expr
	}

	LetStmt struct {
		stmt
		Let    Pos
		Name   *Ident
		Assign Pos
		Value  Expr
	}

	ReturnStmt struct {
		stmt
		Return Pos
		Result Expr
	}

	BadExpr struct {
		expr
		From, To Pos
	}

	BadStmt struct {
		stmt
		From, To Pos
	}

	ParameterList struct {
		Lparen Pos
		List   []*Ident
		Rparen Pos
	}

	FunctionBody struct {
		Lbrace   Pos
		StmtList []Stmt
		Rbrace   Pos
	}
)

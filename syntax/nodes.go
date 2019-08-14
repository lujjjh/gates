package syntax

type Node interface {
	node()
}

type Expr interface {
	Node
	exprNode()
}

type expr struct{}

func (expr) node()     {}
func (expr) exprNode() {}

type Stmt interface {
	Node
	stmtNode()
}

type stmt struct{}

func (stmt) node()     {}
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

	ArrayLitEntry struct {
		Expanded bool
		Value    Expr
	}

	ArrayLit struct {
		expr
		Lbrack   Pos
		ElemList []ArrayLitEntry
		Rbrack   Pos
	}

	MapLitEntry struct {
		Expanded bool
		Key      Expr
		Value    Expr
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

	VarDeclExpr struct {
		expr
		Name        string
		NamePos     Pos
		Initializer Expr
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

	BodyStmt struct {
		stmt
		Lbrace   Pos
		StmtList []Stmt
		Rbrace   Pos
	}

	LetStmt struct {
		stmt
		Let  Pos
		List []Expr
	}

	IfStmt struct {
		stmt
		If         Pos
		Test       Expr
		Consequent Stmt
		Alternate  Stmt
	}

	ForStmt struct {
		stmt
		For         Pos
		Initializer Stmt
		Test        Expr
		Update      Stmt
		Body        Stmt
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

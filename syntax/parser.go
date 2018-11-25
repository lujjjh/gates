package syntax

import (
	"strconv"
)

type parser struct {
	file    *File
	errors  ErrorList
	scanner Scanner

	// Next token
	pos Pos    // token position
	tok Token  // one token look-ahead
	lit string // token literal
}

func (p *parser) init(fset *FileSet, filename string, src []byte) {
	p.file = fset.AddFile(filename, -1, len(src))
	eh := func(pos Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.file, src, eh)

	p.next()
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

// A bailout panic is raised to indicate early termination.
type bailout struct{}

func (p *parser) error(pos Pos, msg string) {
	epos := p.file.Position(pos)

	// If AllErrors is not set, discard errors reported on the same line
	// as the last recorded error and stop parsing if there are more than
	// 10 errors.
	n := len(p.errors)
	if n > 0 && p.errors[n-1].Pos.Line == epos.Line {
		return // discard - likely a spurious error
	}

	p.errors.Add(epos, msg)
	panic(bailout{})
}

func (p *parser) errorExpected(pos Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		switch {
		case p.tok.IsLiteral():
			// print 123 rather than 'INT', etc.
			msg += ", found " + p.lit
		default:
			msg += ", found '" + p.tok.String() + "'"
		}
	}
	p.error(pos, msg)
}

func (p *parser) expect(tok Token) Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

func (p *parser) parseIdent() *Ident {
	pos := p.pos
	name := ""
	if p.tok == IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(IDENT) // use expect() error handling
	}
	return &Ident{NamePos: pos, Name: name}
}

func (p *parser) parseArrayLit() *ArrayLit {
	lbrack := p.expect(LBRACK)
	var elemList []Expr
	if p.tok != RBRACK {
		for {
			elem := p.parseExpr()
			elemList = append(elemList, elem)
			if p.tok != COMMA {
				break
			}
			p.next()
		}
	}
	rbrack := p.expect(RBRACK)
	return &ArrayLit{
		Lbrack:   lbrack,
		ElemList: elemList,
		Rbrack:   rbrack,
	}
}

func (p *parser) parseMapLit() *MapLit {
	lbrace := p.expect(LBRACE)
	var entries []MapLitEntry
	if p.tok != RBRACE {
		for {
			var key Expr
			if p.tok == LBRACK {
				p.next()
				key = p.parseExpr()
				p.expect(RBRACK)
			} else if p.tok == IDENT {
				ident := p.parseIdent()
				key = &Lit{
					ValuePos: ident.NamePos,
					Kind:     STRING,
					Value:    strconv.Quote(ident.Name),
				}
			} else {
				key = p.parseOperand()
			}
			p.expect(COLON)
			value := p.parseExpr()
			entries = append(entries, MapLitEntry{Key: key, Value: value})
			if p.tok != COMMA {
				break
			}
			p.next()
		}
	}
	rbrace := p.expect(RBRACE)
	return &MapLit{Lbrace: lbrace, Entries: entries, Rbrace: rbrace}
}

func (p *parser) parseFunction() *FunctionLit {
	function := p.expect(FUNCTION)
	parameterList := p.parseFunctionParameterList()
	body := p.parseFunctionBody()

	return &FunctionLit{
		Function:      function,
		ParameterList: parameterList,
		Body:          body,
	}
}

func (p *parser) parseFunctionParameterList() *ParameterList {
	lparen := p.expect(LPAREN)

	var list []*Ident
	if p.tok != RPAREN {
		for {
			ident := p.parseIdent()
			list = append(list, ident)
			if p.tok != COMMA {
				break
			}
			p.next()
		}
	}

	rparen := p.expect(RPAREN)

	return &ParameterList{
		Lparen: lparen,
		List:   list,
		Rparen: rparen,
	}
}

func (p *parser) parseFunctionBody() *FunctionBody {
	lbrace := p.expect(LBRACE)

	stmtList := p.parseStmtList()

	rbrace := p.expect(RBRACE)

	return &FunctionBody{
		Lbrace:   lbrace,
		StmtList: stmtList,
		Rbrace:   rbrace,
	}
}

func (p *parser) parseOperand() Expr {
	switch p.tok {
	case IDENT:
		x := p.parseIdent()
		return x

	case NUMBER, STRING, BOOL, NULL:
		x := &Lit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
		p.next()
		return x

	case LPAREN:
		lparen := p.pos
		p.next()
		x := p.parseExpr()
		rparen := p.expect(RPAREN)
		return &ParenExpr{Lparen: lparen, X: x, Rparen: rparen}

	case LBRACK:
		return p.parseArrayLit()

	case LBRACE:
		return p.parseMapLit()

	case FUNCTION:
		return p.parseFunction()
	}

	// we have an error
	pos := p.pos
	p.errorExpected(pos, "operand")
	return &BadExpr{From: pos, To: p.pos}
}

func (p *parser) parseSelector(x Expr) Expr {
	sel := p.parseIdent()
	return &SelectorExpr{X: x, Sel: sel}
}

func (p *parser) parseIndex(x Expr) Expr {
	lbrack := p.expect(LBRACK)
	index := p.parseExpr()
	rbrack := p.expect(RBRACK)

	return &IndexExpr{X: x, Lbrack: lbrack, Index: index, Rbrack: rbrack}
}

func (p *parser) parseCall(fun Expr) *CallExpr {
	lparen := p.expect(LPAREN)
	var list []Expr
	for p.tok != RPAREN && p.tok != EOF {
		list = append(list, p.parseExpr())
		if p.tok != RPAREN {
			p.expect(COMMA)
		}
	}
	rparen := p.expect(RPAREN)

	return &CallExpr{Fun: fun, Lparen: lparen, Args: list, Rparen: rparen}
}

func (p *parser) parsePrimaryExpr() Expr {
	x := p.parseOperand()

L:
	for {
		switch p.tok {
		case PERIOD:
			p.next()
			if p.tok == IDENT {
				x = p.parseSelector(x)
			} else {
				p.errorExpected(p.pos, "selector")
				p.next()
			}
		case LBRACK:
			x = p.parseIndex(x)
		case LPAREN:
			x = p.parseCall(x)
		default:
			break L
		}
	}

	return x
}

func (p *parser) parseUnaryExpr() Expr {
	switch p.tok {
	case ADD, SUB, NOT:
		pos, op := p.pos, p.tok
		p.next()
		x := p.parseUnaryExpr()
		return &UnaryExpr{OpPos: pos, Op: op, X: x}
	}

	return p.parsePrimaryExpr()
}

func (p *parser) tokPrec() (Token, int) {
	tok := p.tok
	return tok, tok.Precedence()
}

func (p *parser) parseBinaryExpr(prec1 int) Expr {
	x := p.parseUnaryExpr()
	for {
		op, oprec := p.tokPrec()
		if oprec < prec1 {
			return x
		}
		pos := p.expect(op)
		y := p.parseBinaryExpr(oprec + 1)
		x = &BinaryExpr{X: x, OpPos: pos, Op: op, Y: y}
	}
}

func (p *parser) parseExpr() Expr {
	return p.parseBinaryExpr(LowestPrec + 1)
}

func (p *parser) parseSimpleStmt() Stmt {
	x := p.parseExpr()

	switch p.tok {
	case ASSIGN:
		pos, tok := p.pos, p.tok
		p.next()
		y := p.parseExpr()
		as := &AssignStmt{Lhs: x, TokPos: pos, Tok: tok, Rhs: y}
		return as
	}

	return &ExprStmt{X: x}
}

func (p *parser) parseVarDecl() Expr {
	namePos := p.pos
	name := p.parseIdent().Name
	var initializer Expr

	if p.tok == ASSIGN {
		p.next()
		initializer = p.parseExpr()
	}

	return &VarDeclExpr{
		Name:        name,
		NamePos:     namePos,
		Initializer: initializer,
	}
}

func (p *parser) parseVarDeclList() []Expr {
	var list []Expr

	for {
		list = append(list, p.parseVarDecl())
		if p.tok != COMMA {
			break
		}
		p.next()
	}

	return list
}

func (p *parser) parseLetStmt() Stmt {
	let := p.expect(LET)

	list := p.parseVarDeclList()
	p.expect(SEMICOLON)
	return &LetStmt{
		Let:  let,
		List: list,
	}
}

func (p *parser) parseBodyStmt() Stmt {
	lbrace := p.expect(LBRACE)
	stmtList := p.parseStmtList()
	rbrace := p.expect(RBRACE)
	return &BodyStmt{
		Lbrace:   lbrace,
		StmtList: stmtList,
		Rbrace:   rbrace,
	}
}

func (p *parser) parseIfStmt() Stmt {
	ifPos := p.expect(IF)
	p.expect(LPAREN)
	test := p.parseExpr()
	p.expect(RPAREN)
	consequent := p.parseStmt()

	var alternate Stmt
	if p.tok == ELSE {
		p.next()
		alternate = p.parseStmt()
	}

	return &IfStmt{
		If:         ifPos,
		Test:       test,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func (p *parser) parseForStmt() Stmt {
	forPos := p.expect(FOR)
	p.expect(LPAREN)
	var initializer Stmt
	if p.tok == LET {
		initializer = p.parseLetStmt()
	}
	p.expect(SEMICOLON)
	var test Expr
	if p.tok != SEMICOLON && p.tok != EOF {
		test = p.parseExpr()
	}
	var update Stmt
	if p.tok != RPAREN && p.tok != EOF {
		update = p.parseSimpleStmt()
	}
	p.expect(RPAREN)

	body := p.parseStmt()

	return &ForStmt{
		For:         forPos,
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        body,
	}
}

func (p *parser) parseReturnStmt() Stmt {
	pos := p.expect(RETURN)
	var result Expr
	if p.tok != SEMICOLON {
		result = p.parseExpr()
	}
	p.expect(SEMICOLON)
	return &ReturnStmt{
		Return: pos,
		Result: result,
	}
}

func (p *parser) parseStmt() Stmt {
	switch p.tok {
	case LBRACE:
		return p.parseBodyStmt()
	case LET:
		return p.parseLetStmt()
	case IDENT:
		s := p.parseSimpleStmt()
		p.expect(SEMICOLON)
		return s
	case IF:
		return p.parseIfStmt()
	case FOR:
		return p.parseForStmt()
	case RETURN:
		return p.parseReturnStmt()
	default:
		pos := p.pos
		p.errorExpected(pos, "statement")
		p.next()
		return &BadStmt{From: pos, To: p.pos}
	}
}

func (p *parser) parseStmtList() (list []Stmt) {
	for p.tok != RBRACE && p.tok != EOF {
		list = append(list, p.parseStmt())
	}

	return
}

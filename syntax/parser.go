package syntax

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

func (p *parser) parseOperand() Expr {
	switch p.tok {
	case IDENT:
		x := p.parseIdent()
		return x

	case NUMBER, STRING, BOOL:
		x := &Lit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
		p.next()
		return x

	case LPAREN:
		lparen := p.pos
		p.next()
		x := p.parseExpr()
		rparen := p.expect(RPAREN)
		return &ParenExpr{Lparen: lparen, X: x, Rparen: rparen}
	}

	// we have an error
	pos := p.pos
	p.errorExpected(pos, "operand")
	return &BadExpr{From: pos, To: p.pos}
}

func (p *parser) parsePrimaryExpr() Expr {
	x := p.parseOperand()
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

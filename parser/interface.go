package parser

import (
	"github.com/lujjjh/gates/ast"
	"github.com/lujjjh/gates/token"
)

func ParseExpr(x string) (e ast.Expr, err error) {
	var p parser

	defer func() {
		if e := recover(); e != nil {
			// resume same panic if it's not a bailout
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}
		p.errors.Sort()
		err = p.errors.Err()
	}()

	// parse expr
	p.init(token.NewFileSet(), "", []byte(x))
	e = p.parseExpr()
	p.expect(token.EOF)

	if p.errors.Len() > 0 {
		p.errors.Sort()
		return nil, p.errors.Err()
	}

	return e, nil
}

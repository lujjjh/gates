package gates

import (
	"fmt"
	"strconv"

	"github.com/lujjjh/gates/syntax"
)

type compiler struct {
	program *Program
	scope   *scope
}

type CompilerError struct {
	Message string
	File    *syntax.File
	Pos     syntax.Pos
}

type CompilerSyntaxError struct {
	CompilerError
}

func (e *CompilerSyntaxError) Error() string {
	if e.File != nil {
		return fmt.Sprintf("SyntaxError: %s at %s", e.Message, e.File.Position(e.Pos))
	}
	return fmt.Sprintf("SyntaxError: %s", e.Message)
}

func (c *compiler) emit(instructions ...instruction) {
	c.program.code = append(c.program.code, instructions...)
}

func (c *compiler) throwSyntaxError(pos syntax.Pos, format string, args ...interface{}) {
	panic(&CompilerSyntaxError{
		CompilerError: CompilerError{
			File:    c.program.src,
			Pos:     pos,
			Message: fmt.Sprintf(format, args...),
		},
	})
}

func (c *compiler) openScope() {
	c.scope = newScope(c.scope)
}

func (c *compiler) closeScope() {
	c.scope = c.scope.outer
}

func (c *compiler) compileAssignStmt(s *syntax.AssignStmt) {
	switch s.Tok {
	case syntax.ASSIGN:
		c.compileExpr(s.Lhs).emitSetter(c.compileExpr(s.Rhs))
	default:
		panic(fmt.Errorf("unknown assign operator: %s", s.Tok.String()))
	}
	return
}

func (c *compiler) compileLetStmt(s *syntax.LetStmt) {
	for _, expr := range s.List {
		c.compileExpr(expr).emitGetter()
	}
}

func (c *compiler) compileIfStmt(s *syntax.IfStmt) {
	c.compileExpr(s.Test).emitGetter()
	jmp := len(c.program.code)
	c.emit(nil)
	c.compileStmt(s.Consequent)
	if s.Alternate != nil {
		jmp2 := len(c.program.code)
		c.emit(nil)
		c.program.code[jmp] = jne(len(c.program.code) - jmp)
		c.compileStmt(s.Alternate)
		c.program.code[jmp2] = jmp1(len(c.program.code) - jmp2)
		return
	}
	c.program.code[jmp] = jne(len(c.program.code) - jmp)
}

func (c *compiler) compileReturnStmt(s *syntax.ReturnStmt) {
	if s.Result == nil {
		c.emit(loadNull)
	} else {
		c.compileExpr(s.Result).emitGetter()
	}
	c.emit(ret)
}

func (c *compiler) compileStmt(s syntax.Stmt) {
	switch s := s.(type) {
	case *syntax.ExprStmt:
		c.compileExpr(s.X).emitGetter()
		c.emit(pop)
	case *syntax.BodyStmt:
		c.openScope()
		c.emit(newStash)
		for _, stmt := range s.StmtList {
			c.compileStmt(stmt)
		}
		c.emit(popStash)
		c.closeScope()
	case *syntax.IfStmt:
		c.compileIfStmt(s)
	case *syntax.AssignStmt:
		c.compileAssignStmt(s)
	case *syntax.LetStmt:
		c.compileLetStmt(s)
	case *syntax.ReturnStmt:
		c.compileReturnStmt(s)
	default:
		panic(fmt.Errorf("unknown statement type: %T", s))
	}
}

func (c *compiler) compileIdent(l *syntax.Ident) compiledExpr {
	idExpr := &compiledIdentExpr{
		name: l.Name,
	}
	idExpr.init(c, l.NamePos)
	return idExpr
}

func (c *compiler) compileLit(l *syntax.Lit) compiledExpr {
	var v Value
	switch l.Kind {
	case syntax.NUMBER:
		i, err := strconv.ParseInt(l.Value, 0, 64)
		if err == nil {
			v = Int(i)
		} else {
			f, _ := strconv.ParseFloat(l.Value, 64)
			v = Float(f)
		}
	case syntax.STRING:
		s, _ := strconv.Unquote(l.Value)
		v = String(s)
	case syntax.BOOL:
		v = Bool(l.Value == "true")
	case syntax.NULL:
		v = Null
	default:
		panic(fmt.Errorf("unknown token type %v", l.Kind))
	}
	lit := &compiledLit{
		value: v,
	}
	lit.init(c, l.ValuePos)
	return lit
}

func (c *compiler) compileArrayLit(e *syntax.ArrayLit) compiledExpr {
	r := &compiledArrayLit{
		expr: e,
	}
	r.init(c, e.Lbrack)
	return r
}

func (c *compiler) compileMapLit(e *syntax.MapLit) compiledExpr {
	r := &compiledMapLit{
		expr: e,
	}
	r.init(c, e.Lbrace)
	return r
}

func (c *compiler) compileFunctionLit(e *syntax.FunctionLit) compiledExpr {
	r := &compiledFunctionLit{
		expr: e,
	}
	r.init(c, e.Function)
	return r
}

func (c *compiler) toStashlessFunction(code []instruction) {
	code[0] = noop
	for i, ins := range code {
		switch ins := ins.(type) {
		case loadLocal:
			level := int(ins >> 24)
			idx := uint32(ins & 0x00FFFFFF)
			level--
			if level < 0 {
				code[i] = loadStack(idx)
				continue
			}
			code[i] = loadLocal(uint32(level)<<24 | idx)
		case storeLocal:
			level := int(ins >> 24)
			idx := uint32(ins & 0x00FFFFFF)
			level--
			if level < 0 {
				code[i] = storeStack(idx)
				continue
			}
			code[i] = storeLocal(uint32(level)<<24 | idx)
		}
	}
}

func (c *compiler) compileUnaryExpr(e *syntax.UnaryExpr) compiledExpr {
	r := &compiledUnaryExpr{
		op: e.Op,
		x:  c.compileExpr(e.X),
	}
	r.init(c, e.OpPos)
	return r
}

func (c *compiler) compileLogicalAnd(x, y syntax.Expr, pos syntax.Pos) compiledExpr {
	r := &compiledLogicalAnd{
		x: c.compileExpr(x),
		y: c.compileExpr(y),
	}
	r.init(c, pos)
	return r
}

func (c *compiler) compileLogicalOr(x, y syntax.Expr, pos syntax.Pos) compiledExpr {
	r := &compiledLogicalOr{
		x: c.compileExpr(x),
		y: c.compileExpr(y),
	}
	r.init(c, pos)
	return r
}

func (c *compiler) compileBinaryExpr(e *syntax.BinaryExpr) compiledExpr {
	switch e.Op {
	case syntax.LAND:
		return c.compileLogicalAnd(e.X, e.Y, e.OpPos)
	case syntax.LOR:
		return c.compileLogicalOr(e.X, e.Y, e.OpPos)
	}

	r := &compiledBinaryExpr{
		x:  c.compileExpr(e.X),
		op: e.Op,
		y:  c.compileExpr(e.Y),
	}
	r.init(c, e.OpPos)
	return r
}

func (c *compiler) compileSelectorExpr(e syntax.Expr, key Value, pos syntax.Pos) compiledExpr {
	lit := &compiledLit{
		value: key,
	}
	lit.init(c, pos)
	r := &compiledSelectorExpr{
		expr: c.compileExpr(e),
		key:  lit,
	}
	r.init(c, pos)
	return r
}

func (c *compiler) compileIndexExpr(e, index syntax.Expr, pos syntax.Pos) compiledExpr {
	r := &compiledIndexExpr{
		expr:  c.compileExpr(e),
		index: c.compileExpr(index),
	}
	r.init(c, pos)
	return r
}

func (c *compiler) compileCallExpr(e *syntax.CallExpr) compiledExpr {
	args := make([]compiledExpr, len(e.Args))
	for i, argExpr := range e.Args {
		args[i] = c.compileExpr(argExpr)
	}
	r := &compiledCallExpr{
		fun:  c.compileExpr(e.Fun),
		args: args,
	}
	r.init(c, e.Lparen)
	return r
}

func (c *compiler) compileVarDeclExpr(e *syntax.VarDeclExpr) compiledExpr {
	var initializer compiledExpr
	if e.Initializer != nil {
		initializer = c.compileExpr(e.Initializer)
	}
	r := &compiledVarDeclExpr{
		name:        e.Name,
		initializer: initializer,
	}
	r.init(c, e.NamePos)
	return r
}

func (c *compiler) compileExpr(e syntax.Expr) compiledExpr {
	switch e := e.(type) {
	case *syntax.Ident:
		return c.compileIdent(e)
	case *syntax.Lit:
		return c.compileLit(e)
	case *syntax.ArrayLit:
		return c.compileArrayLit(e)
	case *syntax.MapLit:
		return c.compileMapLit(e)
	case *syntax.FunctionLit:
		return c.compileFunctionLit(e)
	case *syntax.UnaryExpr:
		return c.compileUnaryExpr(e)
	case *syntax.BinaryExpr:
		return c.compileBinaryExpr(e)
	case *syntax.ParenExpr:
		return c.compileExpr(e.X)
	case *syntax.SelectorExpr:
		return c.compileSelectorExpr(e.X, String(e.Sel.Name), e.Sel.NamePos)
	case *syntax.IndexExpr:
		return c.compileIndexExpr(e.X, e.Index, e.Lbrack)
	case *syntax.CallExpr:
		return c.compileCallExpr(e)
	case *syntax.VarDeclExpr:
		return c.compileVarDeclExpr(e)
	default:
		panic(fmt.Errorf("unknown expression type: %T", e))
	}
}

func (c *compiler) compile(e syntax.Expr) {
	c.compileExpr(e).emitGetter()
	c.emit(halt)
}

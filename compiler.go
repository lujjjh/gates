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

func (c *compiler) emit(instructions ...instruction) {
	c.program.code = append(c.program.code, instructions...)
}

func (c *compiler) compileLetStmt(s *syntax.LetStmt) {
	idx := c.scope.bindName(s.Name.Name)
	if s.Value != nil {
		c.compileExpr(s.Value)
		c.emit(storeLocal(idx))
	}
}

func (c *compiler) compileReturnStmt(s *syntax.ReturnStmt) {
	if s.Result == nil {
		c.emit(loadNull)
	} else {
		c.compileExpr(s.Result)
	}
	c.emit(ret)
}

func (c *compiler) compileStmt(s syntax.Stmt) {
	switch s := s.(type) {
	case *syntax.LetStmt:
		c.compileLetStmt(s)
	case *syntax.ReturnStmt:
		c.compileReturnStmt(s)
	default:
		panic(fmt.Errorf("unknown statement type: %T", s))
	}
}

func (c *compiler) compileIdent(l *syntax.Ident) {
	name := l.Name
	if c.scope != nil {
		idx, ok := c.scope.lookupName(name)
		if ok {
			c.emit(loadLocal(idx))
			return
		}
	}
	c.emit(load(c.program.defineLit(String(l.Name))), loadGlobal, get)
}

func (c *compiler) compileLit(l *syntax.Lit) {
	switch l.Kind {
	case syntax.NUMBER:
		var v Value
		i, err := strconv.ParseInt(l.Value, 0, 64)
		if err == nil {
			v = Int(i)
		} else {
			f, _ := strconv.ParseFloat(l.Value, 64)
			v = Float(f)
		}
		c.emit(load(c.program.defineLit(v)))
	case syntax.STRING:
		s, _ := strconv.Unquote(l.Value)
		c.emit(load(c.program.defineLit(String(s))))
	case syntax.BOOL:
		c.emit(load(c.program.defineLit(Bool(l.Value == "true"))))
	case syntax.NULL:
		c.emit(load(c.program.defineLit(Null)))
	default:
		panic(fmt.Errorf("unknown token type %v", l.Kind))
	}
}

func (c *compiler) compileArrayLit(e *syntax.ArrayLit) {
	for _, elem := range e.ElemList {
		c.compileExpr(elem)
	}
	c.emit(newArray(len(e.ElemList)))
}

func (c *compiler) compileMapLit(e *syntax.MapLit) {
	for _, entry := range e.Entries {
		c.compileExpr(entry.Key)
		c.compileExpr(entry.Value)
	}
	c.emit(newMap(len(e.Entries)))
}

func (c *compiler) compileFunctionLit(e *syntax.FunctionLit) {
	j := len(c.program.code)
	c.emit(nil, nil)
	c.scope = newScope(c.scope)
	for i, ident := range e.ParameterList.List {
		idx := c.scope.bindName(ident.Name)
		c.emit(loadStack(-(i + 1)), storeLocal(idx))
	}
	for _, stmt := range e.Body.StmtList {
		c.compileStmt(stmt)
	}
	c.emit(loadNull, ret)
	c.program.code[j] = newFunc(len(c.scope.names)<<24 | j + 2)
	c.scope = c.scope.outer
	c.program.code[j+1] = jmp1(len(c.program.code) - (j + 1))
}

func (c *compiler) compileUnaryExpr(e *syntax.UnaryExpr) {
	c.compileExpr(e.X)
	switch e.Op {
	case syntax.ADD:
		c.emit(plus)
	case syntax.SUB:
		c.emit(neg)
	case syntax.NOT:
		c.emit(not)
	default:
		panic(fmt.Errorf("unknown unary operator: %s", e.Op))
	}
}

func (c *compiler) compileBinaryExpr(e *syntax.BinaryExpr) {
	switch e.Op {
	case syntax.ADD:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(add)
	case syntax.SUB:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(sub)
	case syntax.MUL:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(mul)
	case syntax.QUO:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(div)
	case syntax.REM:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(mod)
	case syntax.AND:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(and)
	case syntax.OR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(or)
	case syntax.XOR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(xor)
	case syntax.SHL:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(shl)
	case syntax.SHR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(shr)
	case syntax.LAND:
		c.compileExpr(e.X)
		j := len(c.program.code)
		c.emit(nil, pop)
		c.compileExpr(e.Y)
		c.program.code[j] = jneq1(len(c.program.code) - j)
	case syntax.LOR:
		c.compileExpr(e.X)
		j := len(c.program.code)
		c.emit(nil, pop)
		c.compileExpr(e.Y)
		c.program.code[j] = jeq1(len(c.program.code) - j)
	case syntax.EQL:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(eq)
	case syntax.LSS:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(lt)
	case syntax.GTR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(gt)
	case syntax.NEQ:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(neq)
	case syntax.LEQ:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(lte)
	case syntax.GEQ:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.emit(gte)
	default:
		panic(fmt.Errorf("unknown binary operator: %s", e.Op))
	}
}

func (c *compiler) compileSelectorExpr(e syntax.Expr, key Value) {
	c.emit(load(c.program.defineLit(key)))
	c.compileExpr(e)
	c.emit(get)
}

func (c *compiler) compileIndexExpr(e, index syntax.Expr) {
	c.compileExpr(index)
	c.compileExpr(e)
	c.emit(get)
}

func (c *compiler) compileCallExpr(fun syntax.Expr, args []syntax.Expr) {
	for _, arg := range args {
		c.compileExpr(arg)
	}
	c.emit(load(c.program.defineLit(Int(len(args)))))
	c.compileExpr(fun)
	c.emit(call)
}

func (c *compiler) compileExpr(e syntax.Expr) {
	switch e := e.(type) {
	case *syntax.Ident:
		c.compileIdent(e)
	case *syntax.Lit:
		c.compileLit(e)
	case *syntax.ArrayLit:
		c.compileArrayLit(e)
	case *syntax.MapLit:
		c.compileMapLit(e)
	case *syntax.FunctionLit:
		c.compileFunctionLit(e)
	case *syntax.UnaryExpr:
		c.compileUnaryExpr(e)
	case *syntax.BinaryExpr:
		c.compileBinaryExpr(e)
	case *syntax.ParenExpr:
		c.compileExpr(e.X)
	case *syntax.SelectorExpr:
		c.compileSelectorExpr(e.X, String(e.Sel.Name))
	case *syntax.IndexExpr:
		c.compileIndexExpr(e.X, e.Index)
	case *syntax.CallExpr:
		c.compileCallExpr(e.Fun, e.Args)
	default:
		panic(fmt.Errorf("unknown expression type: %T", e))
	}
}

func (c *compiler) compile(e syntax.Expr) {
	c.compileExpr(e)
	c.emit(halt)
}

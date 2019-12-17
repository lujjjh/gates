package gates

import (
	"fmt"

	"github.com/lujjjh/gates/syntax"
)

type compiledExpr interface {
	emitGetter()
	emitSetter(compiledExpr)
}

type baseCompiledExpr struct {
	c   *compiler
	pos syntax.Pos
}

type compiledIdentExpr struct {
	baseCompiledExpr
	name string
}

type compiledLit struct {
	baseCompiledExpr
	value Value
}

type compiledArrayLit struct {
	baseCompiledExpr
	expr *syntax.ArrayLit
}

type compiledMapLit struct {
	baseCompiledExpr
	expr *syntax.MapLit
}

type compiledFunctionLit struct {
	baseCompiledExpr
	expr *syntax.FunctionLit
}

type compiledUnaryExpr struct {
	baseCompiledExpr
	op syntax.Token
	x  compiledExpr
}

type compiledBinaryExpr struct {
	baseCompiledExpr
	x  compiledExpr
	op syntax.Token
	y  compiledExpr
}

type compiledLogicalAnd struct {
	baseCompiledExpr
	x, y compiledExpr
}

type compiledLogicalOr struct {
	baseCompiledExpr
	x, y compiledExpr
}

type compiledSelectorExpr struct {
	baseCompiledExpr
	expr compiledExpr
	key  compiledExpr
}

type compiledIndexExpr struct {
	baseCompiledExpr
	expr  compiledExpr
	index compiledExpr
}

type compiledCallExpr struct {
	baseCompiledExpr
	fun  compiledExpr
	args []compiledExpr
}

type compiledVarDeclExpr struct {
	baseCompiledExpr
	name        string
	initializer compiledExpr
}

func (e *baseCompiledExpr) init(c *compiler, pos syntax.Pos) {
	e.c = c
	e.pos = pos
}

func (e *baseCompiledExpr) emitSetter(valueExpr compiledExpr) {
	e.c.throwSyntaxError(e.pos, "not a valid left-value expression")
}

func (e *compiledIdentExpr) emitGetter() {
	if e.c.scope != nil {
		idx, ok := e.c.scope.lookupName(e.name)
		if ok {
			e.c.emit(loadLocal(idx))
			return
		}
	}
	e.c.emit(load(e.c.program.defineLit(String(e.name))), loadGlobal, get)
}

func (e *compiledIdentExpr) emitSetter(valueExpr compiledExpr) {
	valueExpr.emitGetter()
	if e.c.scope != nil {
		idx, ok := e.c.scope.lookupName(e.name)
		if ok {
			e.c.emit(storeLocal(idx))
			return
		}
	}
	e.c.emit(load(e.c.program.defineLit(String(e.name))), loadGlobal, set)
}

func (e *compiledLit) emitGetter() {
	e.c.emit(load(e.c.program.defineLit(e.value)))
}

func (e *compiledArrayLit) emitGetter() {
	e.c.emit(newArray(0))
	for _, elem := range e.expr.ElemList {
		e.c.compileExpr(elem.Value).emitGetter()
		if elem.Expanded {
			e.c.emit(arrayConcat)
		} else {
			e.c.emit(arrayPush)
		}
	}
}

func (e *compiledMapLit) emitGetter() {
	e.c.emit(newMap(0))
	for _, entry := range e.expr.Entries {
		if entry.Expanded {
			e.c.compileExpr(entry.Value).emitGetter()
			e.c.emit(mapConcat)
		} else {
			e.c.compileExpr(entry.Key).emitGetter()
			e.c.compileExpr(entry.Value).emitGetter()
			e.c.emit(mapSet)
		}
	}
}

func (e *compiledFunctionLit) emitGetter() {
	savedProgram := e.c.program
	p := &Program{
		src: e.c.program.src,
	}
	e.c.program = p
	e.c.emit(newStash)
	e.c.scope = newScope(e.c.scope)
	for i, ident := range e.expr.ParameterList.List {
		idx := e.c.scope.bindName(ident.Name)
		e.c.emit(loadStack(-(i + 1)), storeLocal(idx))
	}
	for _, stmt := range e.expr.Body.StmtList {
		e.c.compileStmt(stmt)
	}
	e.c.emit(loadNull, ret)
	if !e.c.scope.visited {
		e.c.toStashlessFunction(e.c.program.code)
	}
	stackSize := len(e.c.scope.names)
	e.c.scope = e.c.scope.outer
	e.c.program = savedProgram
	e.c.emit(&newFunc{
		program:   p,
		stackSize: stackSize,
	})
}

func (e *compiledUnaryExpr) emitGetter() {
	e.x.emitGetter()
	switch e.op {
	case syntax.ADD:
		e.c.emit(plus)
	case syntax.SUB:
		e.c.emit(neg)
	case syntax.NOT:
		e.c.emit(not)
	default:
		panic(fmt.Errorf("unknown unary operator: %s", e.op))
	}
}

func (e *compiledBinaryExpr) emitGetter() {
	if e.op == syntax.PIPE {
		e.x.emitGetter()
		e.c.emit(load(e.c.program.defineLit(Int(1))))
		e.y.emitGetter()
		e.c.emit(call)
		return
	}
	e.x.emitGetter()
	e.y.emitGetter()
	switch e.op {
	case syntax.ADD:
		e.c.emit(add)
	case syntax.SUB:
		e.c.emit(sub)
	case syntax.MUL:
		e.c.emit(mul)
	case syntax.QUO:
		e.c.emit(div)
	case syntax.REM:
		e.c.emit(mod)
	case syntax.XOR:
		e.c.emit(xor)
	case syntax.SHL:
		e.c.emit(shl)
	case syntax.SHR:
		e.c.emit(shr)
	case syntax.EQL:
		e.c.emit(eq)
	case syntax.LSS:
		e.c.emit(lt)
	case syntax.GTR:
		e.c.emit(gt)
	case syntax.NEQ:
		e.c.emit(neq)
	case syntax.LEQ:
		e.c.emit(lte)
	case syntax.GEQ:
		e.c.emit(gte)
	default:
		panic(fmt.Errorf("unknown binary operator: %s", e.op))
	}
}

func (e *compiledLogicalAnd) emitGetter() {
	e.x.emitGetter()
	j := len(e.c.program.code)
	e.c.emit(nil, pop)
	e.y.emitGetter()
	e.c.program.code[j] = jneq1(len(e.c.program.code) - j)
}

func (e *compiledLogicalOr) emitGetter() {
	e.x.emitGetter()
	j := len(e.c.program.code)
	e.c.emit(nil, pop)
	e.y.emitGetter()
	e.c.program.code[j] = jeq1(len(e.c.program.code) - j)
}

func (e *compiledSelectorExpr) emitGetter() {
	e.key.emitGetter()
	e.expr.emitGetter()
	e.c.emit(get)
}

func (e *compiledSelectorExpr) emitSetter(valueExpr compiledExpr) {
	valueExpr.emitGetter()
	e.key.emitGetter()
	e.expr.emitGetter()
	e.c.emit(set)
}

func (e *compiledIndexExpr) emitGetter() {
	e.index.emitGetter()
	e.expr.emitGetter()
	e.c.emit(get)
}

func (e *compiledIndexExpr) emitSetter(valueExpr compiledExpr) {
	valueExpr.emitGetter()
	e.index.emitGetter()
	e.expr.emitGetter()
	e.c.emit(set)
}

func (e *compiledCallExpr) emitGetter() {
	for _, arg := range e.args {
		arg.emitGetter()
	}
	e.c.emit(load(e.c.program.defineLit(Int(len(e.args)))))
	e.fun.emitGetter()
	e.c.emit(call)
}

func (e *compiledVarDeclExpr) emitGetter() {
	idx := e.c.scope.bindName(e.name)
	if e.initializer != nil {
		e.initializer.emitGetter()
		e.c.emit(storeLocal(idx))
	}
}

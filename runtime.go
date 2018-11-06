package gates

import (
	"github.com/lujjjh/gates/syntax"
)

func Compile(x string) (program *Program, err error) {
	e, err := syntax.ParseExpr(x)
	if err != nil {
		return nil, err
	}

	compiler := &compiler{
		program: &Program{},
	}
	compiler.compile(e)

	return compiler.program, nil
}

type Runtime struct {
	vm *vm
}

func New() *Runtime {
	r := &Runtime{}
	r.init()
	return r
}

func (r *Runtime) init() {
	r.vm = &vm{}
}

func (r *Runtime) RunProgram(program *Program) Value {
	r.vm.program = program
	r.vm.pc = 0
	r.vm.run()
	return r.vm.stack.Pop()
}

func (r *Runtime) RunString(s string) (Value, error) {
	program, err := Compile(s)
	if err != nil {
		return nil, err
	}
	return r.RunProgram(program), nil
}

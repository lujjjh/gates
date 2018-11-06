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
	vm     *vm
	global Ref
}

func New() *Runtime {
	r := &Runtime{}
	r.init()
	return r
}

func (r *Runtime) init() {
	r.vm = &vm{r: r}
}

func (r *Runtime) SetGlobal(global interface{}) {
	r.global = Ref{global}
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

func (r *Runtime) ToValue(i interface{}) Value {
	switch i := i.(type) {
	case nil:
		return Null
	case Value:
		return i
	case string:
		return String(i)
	case bool:
		return Bool(i)
	case int:
		return intNumber(int64(i))
	case int8:
		return intNumber(int64(i))
	case int16:
		return intNumber(int64(i))
	case int32:
		return intNumber(int64(i))
	case int64:
		return intNumber(i)
	case uint:
		return intNumber(int64(i))
	case uint8:
		return intNumber(int64(i))
	case uint16:
		return intNumber(int64(i))
	case uint32:
		return intNumber(int64(i))
	case uint64:
		return intNumber(int64(i))
	case float32:
		return floatNumber(float64(i))
	case float64:
		return floatNumber(i)
	case map[string]interface{}:
		return Map(i)
	default:
		return Ref{i}
	}
}

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
	r.vm.init()
}

func (r *Runtime) Reset() {
	r.vm.init()
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
		return Int(int64(i))
	case int8:
		return Int(int64(i))
	case int16:
		return Int(int64(i))
	case int32:
		return Int(int64(i))
	case int64:
		return Int(i)
	case uint:
		return Int(int64(i))
	case uint8:
		return Int(int64(i))
	case uint16:
		return Int(int64(i))
	case uint32:
		return Int(int64(i))
	case uint64:
		return Int(int64(i))
	case float32:
		return Float(float64(i))
	case float64:
		return Float(i)
	case map[string]interface{}:
		return Map(i)
	case []interface{}:
		return Array(i)
	default:
		return Ref{i}
	}
}

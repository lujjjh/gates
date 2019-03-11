package gates

import (
	"context"
	"unsafe"

	"github.com/gates/gates/syntax"
)

func Compile(x string) (program *Program, err error) {
	defer func() {
		if x := recover(); x != nil {
			program = nil
			switch x1 := x.(type) {
			case *CompilerSyntaxError:
				err = x1
			default:
				panic(x)
			}
		}
	}()

	e, err := syntax.ParseExpr(x)
	if err != nil {
		return nil, err
	}

	compiler := &compiler{
		program: &Program{
			src: syntax.NewFileSet().AddFile("", -1, len(x)),
		},
	}
	compiler.compile(e)

	return compiler.program, nil
}

type Runtime struct {
	vm     *vm
	global *Global
}

func New() *Runtime {
	r := &Runtime{}
	r.init()
	return r
}

func (r *Runtime) init() {
	r.vm = &vm{r: r}
	r.vm.init()
	r.global = NewGlobal()
	r.global.initBuiltIns()

	initPackageStrings(r)
}

func (r *Runtime) Global() *Global {
	return r.global
}

func (r *Runtime) RunProgram(ctx context.Context, program *Program) (Value, error) {
	r.vm.program = program
	r.vm.pc = 0
	r.vm.ctx = ctx
	if err := r.vm.run(); err != nil {
		return nil, err
	}
	return r.vm.stack.Pop(), nil
}

func (r *Runtime) RunString(s string) (Value, error) {
	program, err := Compile(s)
	if err != nil {
		return nil, err
	}
	return r.RunProgram(context.Background(), program)
}

func (r *Runtime) Call(f Function, args ...Value) Value {
	switch f := f.(type) {
	case *nativeFunction:
		return f.fun(&functionCall{vm: r.vm, args: args})
	case *literalFunction:
		vm := r.vm
		for i := range args {
			vm.stack.Push(args[i])
		}
		vm.stack.Push(Int(len(args)))
		pc := vm.pc
		vm.pc = -1
		vm.pushCtx()
		vm.bp = vm.stack.sp
		for i := 0; i < f.stackSize; i++ {
			vm.stack.Push(Null)
		}
		vm.stash = f.stash
		vm.program = f.program
		vm.pc = 0
		if err := vm.run(); err != nil {
			panic(err)
		}
		vm.halt = false
		vm.pc = pc
		return vm.stack.Pop()
	}
	return Null
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
	case map[string]Value:
		return Map(i)
	case []Value:
		return Array(i)
	default:
		return Ref{i}
	}
}

func (r *Runtime) Context() context.Context { return r.vm.ctx }

type toNativer interface {
	toNative(seen map[unsafe.Pointer]interface{}) interface{}
}

func toNative(seen map[unsafe.Pointer]interface{}, v Value) (result interface{}) {
	if seen == nil {
		seen = make(map[unsafe.Pointer]interface{})
	}
	if toNativer, haveToNativer := v.(toNativer); haveToNativer {
		return toNativer.toNative(seen)
	}
	return v.ToNative()
}

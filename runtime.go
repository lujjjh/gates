package gates

import (
	"context"
	"unsafe"

	"github.com/gates/gates/syntax"
)

type ToNativeOption int

const (
	SkipCircularReference ToNativeOption = 1 << iota
)

func checkToNativeOption(desiredOption ToNativeOption, options int) bool {
	return options&int(desiredOption) == int(desiredOption)
}

func convertToNativeOption2BinaryOptions(options []ToNativeOption) int {
	ops := 0
	for _, op := range options {
		ops |= int(op)
	}
	return ops
}

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
	r.global.initBuiltInFunctions()

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
	return ToValue(i)
}

func (r *Runtime) Context() context.Context { return r.vm.ctx }

type toNativer interface {
	toNative(seen map[unsafe.Pointer]interface{}, options int) interface{}
}

func toNative(seen map[unsafe.Pointer]interface{}, v Value, options int) (result interface{}) {
	if seen == nil {
		seen = make(map[unsafe.Pointer]interface{})
	}
	if toNativer, haveToNativer := v.(toNativer); haveToNativer {
		return toNativer.toNative(seen, options)
	}
	return v.ToNative()
}

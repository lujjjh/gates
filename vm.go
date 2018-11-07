package gates

import (
	"math"
	"strings"
)

type valueStack struct {
	l  []Value
	sp int
}

func (v *valueStack) init() {
	v.l = v.l[:0]
	v.sp = 0
}

func (v *valueStack) expand(idx int) {
	if idx < len(v.l) {
		return
	}

	if idx < cap(v.l) {
		v.l = v.l[:idx+1]
	} else {
		n := make([]Value, idx+1, (idx+1)<<1)
		copy(n, v.l)
		v.l = n
	}
}

func (v *valueStack) Push(value Value) {
	v.expand(v.sp)
	v.l[v.sp] = value
	v.sp++
}

func (v *valueStack) Peek() Value {
	return v.l[v.sp-1]
}

func (v *valueStack) Pop() Value {
	v.sp--
	return v.l[v.sp]
}

type vm struct {
	r       *Runtime
	halt    bool
	pc      int
	stack   valueStack
	program *Program
}

func (vm *vm) init() {
	vm.stack.init()
}

func (vm *vm) run() {
	vm.halt = false
	for !vm.halt {
		vm.program.code[vm.pc].exec(vm)
	}
}

type instruction interface {
	exec(*vm)
}

type _halt struct{}

var halt _halt

func (_halt) exec(vm *vm) {
	vm.halt = true
	vm.pc++
}

type load uint

func (index load) exec(vm *vm) {
	vm.stack.Push(vm.program.values[index])
	vm.pc++
}

type _loadGlobal struct{}

var loadGlobal _loadGlobal

func (_loadGlobal) exec(vm *vm) {
	vm.stack.Push(vm.r.global)
	vm.pc++
}

type _pop struct{}

var pop _pop

func (_pop) exec(vm *vm) {
	vm.stack.Pop()
	vm.pc++
}

type _get struct{}

var get _get

func (_get) exec(vm *vm) {
	base := vm.stack.Pop()
	key := vm.stack.Pop()
	vm.stack.Push(objectGet(vm.r, base, key))
	vm.pc++
}

type jeq1 int64

func (j jeq1) exec(vm *vm) {
	if vm.stack.Peek().ToBool() {
		vm.pc += int(j)
	} else {
		vm.pc++
	}
}

type jneq1 int64

func (j jneq1) exec(vm *vm) {
	if !vm.stack.Peek().ToBool() {
		vm.pc += int(j)
	} else {
		vm.pc++
	}
}

type _plus struct{}

var plus _plus

func (_plus) exec(vm *vm) {
	vm.stack.Push(vm.stack.Pop().ToNumber())
	vm.pc++
}

type _neg struct{}

var neg _neg

func (_neg) exec(vm *vm) {
	n := vm.stack.Pop().ToNumber()
	if n.IsInt() {
		vm.stack.Push(Int(-n.ToInt()))
	} else {
		vm.stack.Push(Float(-n.ToFloat()))
	}
	vm.pc++
}

type _not struct{}

var not _not

func (_not) exec(vm *vm) {
	v := vm.stack.Pop().ToBool()
	vm.stack.Push(Bool(!v))
	vm.pc++
}

type _add struct{}

var add _add

func (_add) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()

	switch {
	case x.IsString() || y.IsString():
		xStr, yStr := x.ToString(), y.ToString()
		vm.stack.Push(String(xStr + yStr))
	case x.IsInt() && y.IsInt():
		vm.stack.Push(Int(x.ToInt() + y.ToInt()))
	default:
		vm.stack.Push(Float(x.ToFloat() + y.ToFloat()))
	}

	vm.pc++
}

type _sub struct{}

var sub _sub

func (_sub) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()

	switch {
	case x.IsInt() && y.IsInt():
		vm.stack.Push(Int(x.ToInt() - y.ToInt()))
	default:
		vm.stack.Push(Float(x.ToFloat() - y.ToFloat()))
	}

	vm.pc++
}

type _mul struct{}

var mul _mul

func (_mul) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()

	switch {
	case x.IsInt() && y.IsInt():
		xI, yI := x.ToInt(), y.ToInt()
		res := xI * yI
		// overflow
		if xI != 0 && res/xI != yI {
			vm.stack.Push(Float(x.ToFloat() * y.ToFloat()))
			vm.pc++
			return
		}
		vm.stack.Push(Int(x.ToInt() * y.ToInt()))
	default:
		vm.stack.Push(Float(x.ToFloat() * y.ToFloat()))
	}

	vm.pc++
}

type _div struct{}

var div _div

func (_div) exec(vm *vm) {
	y := vm.stack.Pop().ToFloat()
	x := vm.stack.Pop().ToFloat()

	vm.stack.Push(Float(x / y))

	vm.pc++
}

type _mod struct{}

var mod _mod

func (_mod) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()

	if x.IsInt() && y.IsInt() {
		xI, yI := x.ToInt(), y.ToInt()
		if yI != 0 {
			vm.stack.Push(Int(xI % yI))
			vm.pc++
			return
		}
	}

	vm.stack.Push(Float(math.Mod(x.ToFloat(), y.ToFloat())))
	vm.pc++
}

type _and struct{}

var and _and

func (_and) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(Int(x & y))
	vm.pc++
}

type _or struct{}

var or _or

func (_or) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(Int(x | y))
	vm.pc++
}

type _xor struct{}

var xor _xor

func (_xor) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(Int(x ^ y))
	vm.pc++
}

type _shl struct{}

var shl _shl

func (_shl) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(Int(x << uint64(y)))
	vm.pc++
}

type _shr struct{}

var shr _shr

func (_shr) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(Int(x >> uint64(y)))
	vm.pc++
}

type _eq struct{}

var eq _eq

func (_eq) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	vm.stack.Push(Bool(x.Equals(y)))
	vm.pc++
}

type _neq struct{}

var neq _neq

func (_neq) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	vm.stack.Push(Bool(!x.Equals(y)))
	vm.pc++
}

func less(x, y Value) bool {
	switch {
	case x.IsString() && y.IsString():
		xs, ys := x.ToString(), y.ToString()
		return strings.Compare(xs, ys) == -1
	case x.IsInt() && y.IsInt():
		return x.ToInt() < y.ToInt()
	default:
		return x.ToFloat() < y.ToFloat()
	}
}

type _lt struct{}

var lt _lt

func (_lt) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	vm.stack.Push(Bool(less(x, y)))
	vm.pc++
}

type _lte struct{}

var lte _lte

func (_lte) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	vm.stack.Push(Bool(!less(y, x)))
	vm.pc++
}

type _gt struct{}

var gt _gt

func (_gt) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	vm.stack.Push(Bool(less(y, x)))
	vm.pc++
}

type _gte struct{}

var gte _gte

func (_gte) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	vm.stack.Push(Bool(!less(x, y)))
	vm.pc++
}

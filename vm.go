package gates

import (
	"container/list"
	"math"
	"strings"
)

type valueStack struct{ l list.List }

func (v *valueStack) Push(value Value) {
	v.l.PushBack(value)
}

func (v *valueStack) Peek() Value {
	e := v.l.Back()
	if e == nil {
		return nil
	}
	return e.Value.(Value)
}

func (v *valueStack) Pop() Value {
	e := v.l.Back()
	if e == nil {
		return nil
	}
	v.l.Remove(e)
	return e.Value.(Value)
}

type vm struct {
	halt    bool
	pc      int
	stack   valueStack
	program *Program
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

type _pop struct{}

var pop _pop

func (_pop) exec(vm *vm) {
	vm.stack.Pop()
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
	if n.isInt() {
		vm.stack.Push(intNumber(-n.ToInt()))
	} else {
		vm.stack.Push(floatNumber(-n.ToFloat()))
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
	case x.isString() || y.isString():
		xStr, yStr := x.ToString(), y.ToString()
		vm.stack.Push(String(xStr + yStr))
	case x.isInt() && y.isInt():
		vm.stack.Push(intNumber(x.ToInt() + y.ToInt()))
	default:
		vm.stack.Push(floatNumber(x.ToFloat() + y.ToFloat()))
	}

	vm.pc++
}

type _sub struct{}

var sub _sub

func (_sub) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()

	switch {
	case x.isInt() && y.isInt():
		vm.stack.Push(intNumber(x.ToInt() - y.ToInt()))
	default:
		vm.stack.Push(floatNumber(x.ToFloat() - y.ToFloat()))
	}

	vm.pc++
}

type _mul struct{}

var mul _mul

func (_mul) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()

	switch {
	case x.isInt() && y.isInt():
		xI, yI := x.ToInt(), y.ToInt()
		res := xI * yI
		// overflow
		if xI != 0 && res/xI != yI {
			vm.stack.Push(floatNumber(x.ToFloat() * y.ToFloat()))
			vm.pc++
			return
		}
		vm.stack.Push(intNumber(x.ToInt() * y.ToInt()))
	default:
		vm.stack.Push(floatNumber(x.ToFloat() * y.ToFloat()))
	}

	vm.pc++
}

type _div struct{}

var div _div

func (_div) exec(vm *vm) {
	y := vm.stack.Pop().ToFloat()
	x := vm.stack.Pop().ToFloat()

	vm.stack.Push(floatNumber(x / y))

	vm.pc++
}

type _mod struct{}

var mod _mod

func (_mod) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()

	if x.isInt() && y.isInt() {
		xI, yI := x.ToInt(), y.ToInt()
		if yI != 0 {
			vm.stack.Push(intNumber(xI % yI))
			vm.pc++
			return
		}
	}

	vm.stack.Push(floatNumber(math.Mod(x.ToFloat(), y.ToFloat())))
	vm.pc++
}

type _and struct{}

var and _and

func (_and) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intNumber(x & y))
	vm.pc++
}

type _or struct{}

var or _or

func (_or) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intNumber(x | y))
	vm.pc++
}

type _xor struct{}

var xor _xor

func (_xor) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intNumber(x ^ y))
	vm.pc++
}

type _shl struct{}

var shl _shl

func (_shl) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intNumber(x << uint64(y)))
	vm.pc++
}

type _shr struct{}

var shr _shr

func (_shr) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intNumber(x >> uint64(y)))
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
	case x.isString() && y.isString():
		xs, ys := x.ToString(), y.ToString()
		return strings.Compare(xs, ys) == -1
	case x.isInt() && y.isInt():
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

func (vm *vm) run() {
	vm.halt = false
	for !vm.halt {
		vm.program.code[vm.pc].exec(vm)
	}
}

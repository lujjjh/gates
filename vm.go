package gates

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
)

type valueStack struct {
	l  []Value
	sp int
}

type stash struct {
	values valueStack
	names  map[string]uint32

	outer *stash
}

type ctx struct {
	program *Program
	stash   *stash
	pc, bp  int
}

var ErrStackOverflow = errors.New("stack overflow")

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

func (v *valueStack) PopN(n int) []Value {
	values := v.l[v.sp-n : v.sp]
	v.sp -= n
	return values
}

func (s *stash) putByName(name string, v Value) bool {
	if idx, ok := s.names[name]; ok {
		s.values.expand(int(idx))
		s.values.l[idx] = v
		return true
	}
	return false
}

func (s *stash) putByIdx(idx uint32, v Value) {
	s.values.expand(int(idx))
	s.values.l[idx] = v
}

func (s *stash) getByName(name string) (Value, bool) {
	if idx, ok := s.names[name]; ok {
		return s.values.l[idx], true
	}
	return nil, false
}

func (s *stash) getByIdx(idx uint32) Value {
	if int(idx) < len(s.values.l) {
		return s.values.l[idx]
	}
	return Null
}

type vm struct {
	r         *Runtime
	ctx       context.Context
	halt      bool
	pc        int
	stack     valueStack
	stash     *stash
	callStack []ctx
	bp        int
	program   *Program
}

func (vm *vm) newStash() {
	vm.stash = &stash{
		outer: vm.stash,
	}
}

func (vm *vm) init() {
	vm.stack.init()
	vm.stash = nil
	vm.callStack = nil
}

func (vm *vm) run() (err error) {
	defer func() {
		r := recover()
		if r != nil {
			if rErr, ok := r.(error); ok {
				if rErr == ErrStackOverflow {
					err = rErr
					return
				}
			}
			panic(r)
		}
	}()

	vm.halt = false
	ctx := vm.ctx
	for !vm.halt {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		vm.program.code[vm.pc].exec(vm)
	}
	return nil
}

func (vm *vm) pushCtx() {
	if len(vm.callStack) > 1<<10 {
		panic(ErrStackOverflow)
	}
	vm.callStack = append(vm.callStack, ctx{
		program: vm.program,
		stash:   vm.stash,
		pc:      vm.pc,
		bp:      vm.bp,
	})
}

func (vm *vm) popCtx() {
	l := len(vm.callStack) - 1
	vm.program = vm.callStack[l].program
	vm.stash = vm.callStack[l].stash
	vm.pc = vm.callStack[l].pc
	vm.bp = vm.callStack[l].bp
	vm.callStack = vm.callStack[:l]
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

type _noop struct{}

var noop _noop

func (_noop) exec(vm *vm) {
	vm.pc++
}

type load uint

func (index load) exec(vm *vm) {
	vm.stack.Push(vm.program.values[index])
	vm.pc++
}

type _loadNull struct{}

var loadNull _loadNull

func (_loadNull) exec(vm *vm) {
	vm.stack.Push(Null)
	vm.pc++
}

type _loadGlobal struct{}

var loadGlobal _loadGlobal

func (_loadGlobal) exec(vm *vm) {
	vm.stack.Push(vm.r.global.m)
	vm.pc++
}

type loadStack int

func (l loadStack) exec(vm *vm) {
	idx := int(l)
	bp := vm.bp
	if l < 0 {
		argc := int(vm.stack.l[bp-1].ToInt())
		argn := -idx - 1
		if argn >= argc {
			vm.stack.Push(Null)
		} else {
			vm.stack.Push(vm.stack.l[bp-1-argc+argn])
		}
	} else {
		vm.stack.Push(vm.stack.l[bp+idx])
	}
	vm.pc++
}

type storeStack uint32

func (s storeStack) exec(vm *vm) {
	idx := int(s)
	bp := vm.bp
	vm.stack.l[bp+idx] = vm.stack.Pop()
	vm.pc++
}

type loadLocal uint32

func (l loadLocal) exec(vm *vm) {
	level := l >> 24
	idx := uint32(l & 0x00FFFFFF)
	stash := vm.stash
	for ; level > 0; level-- {
		stash = stash.outer
	}
	vm.stack.Push(stash.getByIdx(idx))
	vm.pc++
}

type storeLocal uint32

func (s storeLocal) exec(vm *vm) {
	v := vm.stack.Pop()
	level := s >> 24
	idx := uint32(s & 0x00FFFFFF)
	stash := vm.stash
	for ; level > 0; level-- {
		stash = stash.outer
	}
	stash.putByIdx(idx, v)
	vm.pc++
}

type _pop struct{}

var pop _pop

func (_pop) exec(vm *vm) {
	vm.stack.Pop()
	vm.pc++
}

type newArray uint

func (l newArray) exec(vm *vm) {
	values := make([]Value, l)
	copy(values, vm.stack.PopN(int(l)))
	vm.stack.Push(Array(values))
	vm.pc++
}

type newMap uint

func (l newMap) exec(vm *vm) {
	m := make(Map, l)
	ll := int(l) * 2
	kvs := vm.stack.PopN(ll)
	for i := 0; i < ll; i += 2 {
		key := kvs[i]
		value := kvs[i+1]
		m[key.ToString()] = value
	}
	vm.stack.Push(m)
	vm.pc++
}

type newFunc struct {
	program   *Program
	stackSize int
}

func (l newFunc) exec(vm *vm) {
	f := &literalFunction{
		program:   l.program,
		stackSize: l.stackSize,
		stash:     vm.stash,
	}
	vm.stack.Push(f)
	vm.pc++
}

type _newStash struct{}

var newStash _newStash

func (_newStash) exec(vm *vm) {
	vm.newStash()
	vm.pc++
}

type _popStash struct{}

var popStash _popStash

func (_popStash) exec(vm *vm) {
	vm.stash = vm.stash.outer
	vm.pc++
}

type _set struct{}

var set _set

func (_set) exec(vm *vm) {
	base := vm.stack.Pop()
	key := vm.stack.Pop()
	value := vm.stack.Pop()
	objectSet(vm.r, base, key, value)
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

type jmp1 int64

func (j jmp1) exec(vm *vm) {
	vm.pc += int(j)
}

type jne int32

func (j jne) exec(vm *vm) {
	if !vm.stack.Pop().ToBool() {
		vm.pc += int(j)
	} else {
		vm.pc++
	}
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
		vm.stack.Push(intToValue(-n.ToInt()))
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
		vm.stack.Push(intToValue(x.ToInt() + y.ToInt()))
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
		vm.stack.Push(intToValue(x.ToInt() - y.ToInt()))
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
		vm.stack.Push(intToValue(x.ToInt() * y.ToInt()))
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
			vm.stack.Push(intToValue(xI % yI))
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
	vm.stack.Push(intToValue(x & y))
	vm.pc++
}

type _or struct{}

var or _or

func (_or) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intToValue(x | y))
	vm.pc++
}

type _xor struct{}

var xor _xor

func (_xor) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intToValue(x ^ y))
	vm.pc++
}

type _shl struct{}

var shl _shl

func (_shl) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intToValue(x << uint64(y)))
	vm.pc++
}

type _shr struct{}

var shr _shr

func (_shr) exec(vm *vm) {
	y := vm.stack.Pop().ToInt()
	x := vm.stack.Pop().ToInt()
	vm.stack.Push(intToValue(x >> uint64(y)))
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

func less(x, y Value) Value {
	switch {
	case x.IsString() && y.IsString():
		xs, ys := x.ToString(), y.ToString()
		return Bool(strings.Compare(xs, ys) == -1)
	case x.IsInt() && y.IsInt():
		return Bool(x.ToInt() < y.ToInt())
	}

	nx := x.ToFloat()
	ny := y.ToFloat()

	if math.IsNaN(nx) || math.IsNaN(ny) {
		return Null
	}
	return Bool(nx < ny)
}

type _lt struct{}

var lt _lt

func (_lt) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	ret := less(x, y)
	if ret == Null {
		vm.stack.Push(False)
	} else {
		vm.stack.Push(ret)
	}
	vm.pc++
}

type _lte struct{}

var lte _lte

func (_lte) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	ret := less(y, x)
	if ret == Null || ret == True {
		vm.stack.Push(False)
	} else {
		vm.stack.Push(True)
	}
	vm.pc++
}

type _gt struct{}

var gt _gt

func (_gt) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	ret := less(y, x)
	if ret == Null {
		vm.stack.Push(False)
	} else {
		vm.stack.Push(ret)
	}
	vm.pc++
}

type _gte struct{}

var gte _gte

func (_gte) exec(vm *vm) {
	y := vm.stack.Pop()
	x := vm.stack.Pop()
	ret := less(x, y)
	if ret == Null || ret == True {
		vm.stack.Push(False)
	} else {
		vm.stack.Push(True)
	}
	vm.pc++
}

type _call struct{}

var call _call

func (_call) exec(vm *vm) {
	fun := vm.stack.Pop().ToFunction()
	switch f := fun.(type) {
	case *nativeFunction:
		argc := int(vm.stack.Pop().ToInt())
		args := vm.stack.PopN(argc)
		fc := &functionCall{vm: vm, args: args}
		vm.stack.Push(f.fun(fc))
		vm.pc++
	case *literalFunction:
		vm.pc++
		vm.pushCtx()
		vm.bp = vm.stack.sp
		for i := 0; i < f.stackSize; i++ {
			vm.stack.Push(Null)
		}
		vm.stash = f.stash
		vm.program = f.program
		vm.pc = 0
	default:
		panic(fmt.Errorf("unsupported function type: %T", fun))
	}
}

type _ret struct{}

var ret _ret

func (_ret) exec(vm *vm) {
	argc := int(vm.stack.l[vm.bp-1].ToInt())
	returnValue := vm.stack.Pop()
	vm.stack.sp = vm.bp - 1 - argc
	vm.stack.l = vm.stack.l[:vm.stack.sp]
	vm.stack.Push(returnValue)
	vm.popCtx()
	if vm.pc < 0 {
		vm.halt = true
	}
}

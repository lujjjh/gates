package vm

import (
	"fmt"
	"math"
	"strings"

	"github.com/gates/gates"
)

const (
	StackSize = 4096
	MaxFrames = 1024
)

type VM struct {
	constants []gates.Value
	globals   []gates.Value
	stack     [StackSize]gates.Value
	sp        int

	frames      [MaxFrames]Frame
	framesIndex int
	curFrame    *Frame
	curInsts    []byte

	ip int

	err error
}

func New(fn *gates.CompiledFunction) *VM {
	v := &VM{
		framesIndex: -1,
		ip:          -1,
	}
	v.frames[0].fn = fn
	v.frames[0].ip = -1
	v.curFrame = &v.frames[0]
	v.curInsts = fn.Instructions
	return v
}

func (v *VM) Run() error {
	v.sp = 0
	v.curFrame = &v.frames[0]
	v.curInsts = v.curFrame.fn.Instructions
	v.framesIndex = 1
	v.ip = -1

	v.run()
	return v.err
}

func (v *VM) run() {
	defer func() {
		if r := recover(); r != nil {
			if v.sp < 0 || v.sp >= StackSize {
				v.err = gates.ErrStackOverflow
				return
			}

			if v.ip < len(v.curInsts)-1 {
				if err, ok := r.(error); ok {
					v.err = err
				} else {
					v.err = fmt.Errorf("panic: %v", r)
				}
			}
		}
	}()

	for {
		v.ip++
		switch v.curInsts[v.ip] {
		case OpLoadConst:
			v.ip += 2
			idx := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
			v.stack[v.sp] = v.constants[idx]
			v.sp++

		case OpLoadNull:
			v.stack[v.sp] = gates.Null
			v.sp++

		case OpLoadGlobal:
			v.ip += 2
			idx := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
			v.stack[v.sp] = v.globals[idx]
			v.sp++

		case OpStoreGlobal:
			v.ip += 2
			idx := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
			v.sp--
			v.globals[idx] = v.stack[v.sp]

		case OpLoadLocal:
			v.ip++
			idx := int(v.curInsts[v.ip])

			val := v.stack[v.curFrame.bp+idx]

			v.stack[v.sp] = val
			v.sp++

		case OpStoreLocal:
			v.ip++
			idx := int(v.curInsts[v.ip])

			val := v.stack[v.sp-1]
			v.sp--

			v.stack[v.curFrame.bp+idx] = val

		case OpUnaryPlus:
			v.stack[v.sp-1] = v.stack[v.sp-1].ToNumber()

		case OpUnaryMinus:
			n := v.stack[v.sp-1].ToNumber()
			if n.IsInt() {
				v.stack[v.sp-1] = gates.Int(-n.ToInt())
			} else {
				v.stack[v.sp-1] = gates.Float(-n.ToFloat())
			}

		case OpUnaryNot:
			v.stack[v.sp-1] = gates.Bool(!v.stack[v.sp-1].ToBool())

		case OpBinaryAdd:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			if x.IsString() || y.IsString() {
				v.stack[v.sp-2] = gates.String(x.ToString() + y.ToString())
			} else if x.IsInt() && y.IsInt() {
				v.stack[v.sp-2] = gates.Int(x.ToInt() + y.ToInt())
			} else {
				v.stack[v.sp-2] = gates.Float(x.ToFloat() + y.ToFloat())
			}
			v.sp--

		case OpBinarySub:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			if x.IsInt() && y.IsInt() {
				v.stack[v.sp-2] = gates.Int(x.ToInt() - y.ToInt())
			} else {
				v.stack[v.sp-2] = gates.Float(x.ToFloat() - y.ToFloat())
			}
			v.sp--

		case OpBinaryMul:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			if x.IsInt() && y.IsInt() {
				xI, yI := x.ToInt(), y.ToInt()
				res := xI * yI
				// overflow
				if xI != 0 && res/xI != yI {
					v.stack[v.sp-2] = gates.Float(x.ToFloat() * y.ToFloat())
				} else {
					v.stack[v.sp-2] = gates.Int(x.ToInt() * y.ToInt())
				}
			} else {
				v.stack[v.sp-2] = gates.Float(x.ToFloat() * y.ToFloat())
			}
			v.sp--

		case OpBinaryDiv:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = gates.Float(x.ToFloat() / y.ToFloat())
			v.sp--

		case OpBinaryMod:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			if x.IsInt() && y.IsInt() && y.ToInt() != 0 {
				v.stack[v.sp-2] = gates.Int(x.ToInt() % y.ToInt())
			} else {
				v.stack[v.sp-2] = gates.Float(math.Mod(x.ToFloat(), y.ToFloat()))
			}
			v.sp--

		case OpBinaryEq:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = gates.Bool(x.Equals(y))
			v.sp--

		case OpBinaryNEq:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = gates.Bool(!x.Equals(y))
			v.sp--

		case OpBinaryLT:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = gates.Bool(less(x, y, false))
			v.sp--

		case OpBinaryLTE:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = gates.Bool(!less(y, x, true))
			v.sp--

		case OpBinaryGT:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = gates.Bool(less(y, x, false))
			v.sp--

		case OpBinaryGTE:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = gates.Bool(!less(x, y, true))
			v.sp--

		case OpCall:
			v.ip++
			numArgs := int(v.curInsts[v.ip])
			callee := v.stack[v.sp-1-numArgs]
			switch callee := callee.(type) {
			case *gates.CompiledFunction:
				v.curFrame.ip = v.ip
				v.curFrame = &v.frames[v.framesIndex]
				v.curFrame.fn = callee
				v.curFrame.bp = v.sp - numArgs
				v.curInsts = callee.Instructions
				v.ip = -1
				v.framesIndex++
				v.sp -= numArgs - callee.NumLocals
			}

		case OpReturn:
			retVal := v.stack[v.sp-1]
			v.framesIndex--
			v.curFrame = &v.frames[v.framesIndex-1]
			v.curInsts = v.curFrame.fn.Instructions
			v.ip = v.curFrame.ip

			v.sp = v.frames[v.framesIndex].bp

			v.stack[v.sp-1] = retVal
		}
	}
}

func less(x, y gates.Value, defaults bool) bool {
	switch {
	case x.IsString() && y.IsString():
		xs, ys := x.ToString(), y.ToString()
		return strings.Compare(xs, ys) == -1
	case x.IsInt() && y.IsInt():
		return x.ToInt() < y.ToInt()
	}

	nx := x.ToFloat()
	ny := y.ToFloat()

	if math.IsNaN(nx) || math.IsNaN(ny) {
		return defaults
	}
	return nx < ny
}

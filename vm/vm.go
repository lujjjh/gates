package vm

import (
	"errors"
	"fmt"

	"github.com/lujjjh/gates"
)

const (
	StackSize = 4096
	MaxFrames = 1024
)

var (
	ErrStackOverflow = errors.New("stack overflow")
)

type VM struct {
	constants []interface{}
	globals   []interface{}
	stack     [StackSize]interface{}
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
				v.err = ErrStackOverflow
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
			v.stack[v.sp] = nil
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

		case OpArray:
			v.ip += 2
			size := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
			result := make([]interface{}, size)
			if size > 0 {
				copy(result, v.stack[v.sp-size:v.sp])
			}
			v.sp -= size
			v.stack[v.sp] = result
			v.sp++

		case OpMergeArray:
			v.ip++
			numSegments := int(v.curInsts[v.ip])
			result := make([]interface{}, 0)
			for _, iterable := range v.stack[v.sp-numSegments : v.sp] {
				it, ok := iterable.([]interface{})
				if !ok {
					continue
				}
				for _, value := range it {
					result = append(result, value)
				}
			}
			v.sp -= numSegments
			v.stack[v.sp] = result
			v.sp++

		case OpMap:
			v.ip += 2
			size := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
			result := make(map[string]interface{}, size)
			offset := size * 2
			for i := v.sp - offset; i < v.sp; i += 2 {
				key := fmt.Sprint(v.stack[i]) // TODO implement toString
				value := v.stack[i+1]
				result[key] = value
			}
			v.sp -= offset
			v.stack[v.sp] = result
			v.sp++

		case OpMergeMap:
			v.ip++
			numSegments := int(v.curInsts[v.ip])
			result := make(map[string]interface{})
			for _, iterable := range v.stack[v.sp-numSegments : v.sp] {
				it, ok := iterable.(map[string]interface{})
				if !ok {
					continue
				}
				for k, v := range it {
					result[k] = v
				}
			}
			v.sp -= numSegments
			v.stack[v.sp] = result
			v.sp++

		case OpUnaryPlus:
			v.stack[v.sp-1] = toNumber(v.stack[v.sp-1])

		case OpUnaryMinus:
			v.stack[v.sp-1] = negate(v.stack[v.sp-1])

		case OpUnaryNot:
			v.stack[v.sp-1] = !toBool(v.stack[v.sp-1])

		case OpBinaryAdd:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = add(x, y)
			v.sp--

		case OpBinarySub:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = sub(x, y)
			v.sp--

		case OpBinaryMul:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = mul(x, y)
			v.sp--

		case OpBinaryDiv:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = div(x, y)
			v.sp--

		case OpBinaryMod:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = mod(x, y)
			v.sp--

		case OpBinaryEq:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = x == y
			v.sp--

		case OpBinaryNEq:
			x, y := v.stack[v.sp-2], v.stack[v.sp-1]
			v.stack[v.sp-2] = x != y
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

package gates

import (
	"math"
)

type Function interface {
	Value
	function()
}

type FunctionCall interface {
	Args() []Value
}

type functionCall struct {
	args []Value
}

func (fc *functionCall) Args() []Value { return fc.args }

func FunctionFunc(fun func(FunctionCall) Value) Function {
	return &nativeFunction{fun: fun}
}

var _EmptyFunction = FunctionFunc(func(FunctionCall) Value { return Null })

type nativeFunction struct {
	fun func(FunctionCall) Value
}

func (*nativeFunction) function() {}

func (*nativeFunction) IsString() bool   { return false }
func (*nativeFunction) IsInt() bool      { return false }
func (*nativeFunction) IsFloat() bool    { return false }
func (*nativeFunction) IsBool() bool     { return false }
func (*nativeFunction) IsFunction() bool { return true }

func (*nativeFunction) ToString() string        { return "function () { [ native code ] }" }
func (*nativeFunction) ToInt() int64            { return 0 }
func (*nativeFunction) ToFloat() float64        { return math.NaN() }
func (*nativeFunction) ToNumber() Number        { return Int(0) }
func (*nativeFunction) ToBool() bool            { return true }
func (f *nativeFunction) ToFunction() Function  { return f }
func (f *nativeFunction) ToNative() interface{} { return f.fun }

func (f *nativeFunction) Equals(other Value) bool {
	return (interface{})(f) == (interface{})(other)
}

func (f *nativeFunction) SameAs(other Value) bool { return f.Equals(other) }

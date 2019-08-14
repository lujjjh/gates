package gates

import (
	"math"
	"strconv"
)

type Int int64
type Float float64

type Number interface {
	Value
	number()
}

func (Int) number()   {}
func (Float) number() {}

func (Int) IsString() bool     { return false }
func (Int) IsInt() bool        { return true }
func (Int) IsFloat() bool      { return false }
func (i Int) ToNumber() Number { return i }
func (Int) IsBool() bool       { return false }
func (Int) IsFunction() bool   { return false }

func (i Int) ToString() string                       { return strconv.FormatInt(int64(i), 10) }
func (i Int) ToInt() int64                           { return int64(i) }
func (i Int) ToFloat() float64                       { return float64(i) }
func (i Int) ToBool() bool                           { return int64(i) != 0 }
func (i Int) ToFunction() Function                   { return _EmptyFunction }
func (i Int) ToNative(...ToNativeOption) interface{} { return i.ToInt() }

func (i Int) Equals(other Value) bool {
	switch {
	case other.IsInt():
		return i.ToInt() == other.ToInt()
	case other.IsFloat():
		return i.ToFloat() == other.ToFloat()
	case other.IsString():
		return other.Equals(i)
	case other.IsBool():
		return i.ToInt() == other.ToInt()
	}
	return false
}

func (i Int) SameAs(b Value) bool {
	ib, ok := b.(Int)
	if !ok {
		return false
	}
	return i == ib
}

func (Float) IsString() bool     { return false }
func (Float) IsInt() bool        { return false }
func (Float) IsFloat() bool      { return true }
func (f Float) ToNumber() Number { return f }
func (Float) IsBool() bool       { return false }
func (Float) IsFunction() bool   { return false }

func (f Float) ToString() string                       { return strconv.FormatFloat(float64(f), 'g', -1, 64) }
func (f Float) ToInt() int64                           { return int64(f) }
func (f Float) ToFloat() float64                       { return float64(f) }
func (f Float) ToBool() bool                           { return float64(f) != 0 && !math.IsNaN(float64(f)) }
func (f Float) ToFunction() Function                   { return _EmptyFunction }
func (f Float) ToNative(...ToNativeOption) interface{} { return f.ToFloat() }

func (f Float) Equals(other Value) bool {
	switch {
	case other.IsInt():
		return f.ToFloat() == other.ToFloat()
	case other.IsFloat():
		return f.ToFloat() == other.ToFloat()
	case other.IsString():
		return other.Equals(f)
	case other.IsBool():
		return f.ToInt() == other.ToInt()
	default:
		return false
	}
}

func (f Float) SameAs(b Value) bool {
	fb, ok := b.(Float)
	if !ok {
		return false
	}
	return f == fb
}

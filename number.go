package gates

import (
	"strconv"
)

type Int int64
type Float float64

type Number interface {
	Value
	number()
}

func (i Int) number()   {}
func (f Float) number() {}

func (i Int) IsString() bool   { return false }
func (i Int) IsInt() bool      { return true }
func (i Int) IsFloat() bool    { return false }
func (i Int) ToNumber() Number { return i }
func (i Int) IsBool() bool     { return false }

func (i Int) ToString() string { return strconv.FormatInt(int64(i), 10) }
func (i Int) ToInt() int64     { return int64(i) }
func (i Int) ToFloat() float64 { return float64(i) }
func (i Int) ToBool() bool     { return int64(i) != 0 }

func (i Int) Equals(other Value) bool {
	switch {
	case other.IsInt():
		return i.ToInt() == other.ToInt()
	case other.IsFloat():
		return i.ToFloat() == other.ToFloat()
	case other.IsString():
		return other.ToNumber().Equals(i)
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

func (f Float) IsString() bool   { return false }
func (f Float) IsInt() bool      { return false }
func (f Float) IsFloat() bool    { return true }
func (f Float) ToNumber() Number { return f }
func (f Float) IsBool() bool     { return false }

func (f Float) ToString() string { return strconv.FormatFloat(float64(f), 'g', -1, 64) }
func (f Float) ToInt() int64     { return int64(f) }
func (f Float) ToFloat() float64 { return float64(f) }
func (f Float) ToBool() bool     { return float64(f) != 0 }

func (f Float) Equals(other Value) bool {
	switch {
	case other.IsFloat():
		return f.ToFloat() == other.ToFloat()
	case other.IsString():
		return other.ToNumber().Equals(f)
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

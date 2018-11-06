package gates

import (
	"strconv"
)

type Number struct {
	isI bool
	i   int64
	f   float64
}

func intNumber(i int64) Number     { return Number{true, i, 0} }
func floatNumber(f float64) Number { return Number{false, 0, f} }

func (n Number) IsString() bool { return false }
func (n Number) IsInt() bool    { return n.isI }
func (n Number) IsFloat() bool  { return !n.isI }
func (n Number) IsBool() bool   { return false }

func (n Number) ToString() string {
	if n.isI {
		return strconv.FormatInt(n.i, 10)
	}
	return strconv.FormatFloat(n.f, 'g', -1, 64)
}

func (n Number) ToInt() int64 {
	if n.isI {
		return n.i
	}
	return int64(n.f)
}

func (n Number) ToFloat() float64 {
	if !n.isI {
		return n.f
	}
	return float64(n.i)
}

func (n Number) ToNumber() Number { return n }

func (n Number) ToBool() bool {
	if n.isI {
		return n.i != 0
	}
	return n.f != 0
}

func (n Number) Equals(other Value) bool {
	switch {
	case n.IsInt() && other.IsInt():
		return n.ToInt() == other.ToInt()
	case other.IsFloat():
		return n.ToFloat() == other.ToFloat()
	case other.IsString():
		return other.ToNumber().Equals(n)
	case other.IsBool():
		return n.ToInt() == other.ToInt()
	default:
		return false
	}
}

func (n Number) SameAs(b Value) bool {
	nb, ok := b.(Number)
	if !ok {
		return false
	}
	return n == nb
}

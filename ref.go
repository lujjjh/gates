package gates

import (
	"math"
)

type Ref struct {
	v interface{}
}

func (Ref) IsString() bool { return false }
func (Ref) IsInt() bool    { return false }
func (Ref) IsFloat() bool  { return false }
func (Ref) IsBool() bool   { return false }

func (Ref) ToString() string     { return "[object Ref]" }
func (Ref) ToInt() int64         { return 0 }
func (Ref) ToFloat() float64     { return math.NaN() }
func (ref Ref) ToNumber() Number { return floatNumber(ref.ToFloat()) }
func (Ref) ToBool() bool         { return true }

func (ref Ref) Equals(other Value) bool {
	if o, ok := other.(Ref); ok {
		return ref.v == o.v
	}
	return false
}

func (ref Ref) SameAs(other Value) bool { return ref.Equals(other) }

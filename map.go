package gates

import (
	"math"
)

type Map map[string]interface{}

func (Map) IsString() bool   { return false }
func (Map) IsInt() bool      { return false }
func (Map) IsFloat() bool    { return false }
func (Map) IsBool() bool     { return false }
func (Map) IsFunction() bool { return false }

func (Map) ToString() string     { return "[object Map]" }
func (Map) ToInt() int64         { return 0 }
func (Map) ToFloat() float64     { return math.NaN() }
func (m Map) ToNumber() Number   { return Float(m.ToFloat()) }
func (Map) ToBool() bool         { return true }
func (Map) ToFunction() Function { return _EmptyFunction }

func (m Map) Equals(other Value) bool { return Value(m) == other }
func (m Map) SameAs(other Value) bool { return m.Equals(other) }

func (m Map) Get(r *Runtime, key Value) Value {
	if m == nil {
		return Null
	}
	return r.ToValue(m[key.ToString()])
}

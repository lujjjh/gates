package gates

import (
	"math"
	"reflect"
)

type Map map[string]Value

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

func (m Map) Equals(other Value) bool {
	o, ok := other.(Map)
	if !ok {
		return false
	}
	return reflect.DeepEqual(m, o)
}

func (m Map) SameAs(other Value) bool { return false }

func (m Map) Get(r *Runtime, key Value) Value {
	if m == nil {
		return Null
	}
	return r.ToValue(m[key.ToString()])
}

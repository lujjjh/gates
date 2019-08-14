package gates

import (
	"math"
	"reflect"
	"strings"
	"unsafe"
)

type Array struct {
	Values []Value
}

type arrayIter struct {
	i int
	a *Array
}

func NewArray(values []Value) Array {
	return Array{
		Values: values,
	}
}

func NewArrayFromStringSlice(a []string) Array {
	values := make([]Value, len(a))
	for i := range a {
		values[i] = String(a[i])
	}
	return NewArray(values)
}

func (Array) Type() string { return "array" }

func (Array) IsString() bool   { return false }
func (Array) IsInt() bool      { return false }
func (Array) IsFloat() bool    { return false }
func (Array) IsBool() bool     { return false }
func (Array) IsFunction() bool { return false }

func (a Array) ToString() string {
	stringSl := make([]string, 0, len(a.Values))
	for _, v := range a.Values {
		stringSl = append(stringSl, v.ToString())
	}
	return strings.Join(stringSl, ",")
}

func (Array) ToInt() int64         { return 0 }
func (Array) ToFloat() float64     { return math.NaN() }
func (a Array) ToNumber() Number   { return Float(a.ToFloat()) }
func (Array) ToBool() bool         { return true }
func (Array) ToFunction() Function { return _EmptyFunction }

func (a Array) ToNative(ops ...ToNativeOption) interface{} {
	return toNative(nil, a, convertToNativeOption2BinaryOptions(ops))
}

func (a Array) toNative(seen map[unsafe.Pointer]interface{}, ops int) interface{} {
	if a.Values == nil {
		return []interface{}(nil)
	}
	addr := unsafe.Pointer(reflect.ValueOf(a.Values).Pointer())
	if v, ok := seen[addr]; ok && !checkToNativeOption(SkipCircularReference, ops) {
		return v
	} else if ok {
		return nil
	}
	result := make([]interface{}, len(a.Values))
	seen[addr] = result
	for i := range a.Values {
		result[i] = toNative(seen, a.Values[i], ops)
	}
	delete(seen, addr)
	return result
}

func (a Array) Equals(other Value) bool {
	o, ok := other.(Array)
	if !ok {
		return false
	}
	return reflect.DeepEqual(a.Values, o.Values)
}

func (a Array) SameAs(other Value) bool { return false }

func (a Array) Get(r *Runtime, key Value) Value {
	if key.IsInt() {
		i := key.ToInt()
		if i < 0 || i >= int64(len(a.Values)) {
			return Null
		}
		return a.Values[i]
	}

	switch key.ToString() {
	case "length":
		return Int(len(a.Values))
	}

	return Null
}

func (a Array) Set(r *Runtime, key, value Value) {
	if !key.IsInt() {
		return
	}
	i := key.ToInt()
	if i < 0 || i >= int64(len(a.Values)) {
		return
	}
	a.Values[i] = value
}

func (a Array) Iterator() Iterator {
	return &arrayIter{i: 0, a: &a}
}

func (a *arrayIter) Next() (Value, bool) {
	i := a.i
	if i >= 0 && i < len(a.a.Values) {
		a.i++
		return a.a.Values[i], true
	}
	return Null, false
}

func (a *Array) push(value Value) {
	a.Values = append(a.Values, value)
}

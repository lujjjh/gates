package gates

import (
	"math"
	"reflect"
	"strings"
)

type Array struct {
	values []Value
}

type arrayIter struct {
	i int
	a *Array
}

func NewArray(values []Value) Array {
	return Array{
		values: values,
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
	stringSl := make([]string, 0, len(a.values))
	for _, v := range a.values {
		stringSl = append(stringSl, ToValue(v).ToString())
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

func (a Array) toNative(seen map[interface{}]interface{}, ops int) interface{} {
	if a.values == nil {
		return []interface{}(nil)
	}
	v := reflect.ValueOf(a.values)
	ptr := struct {
		ptr uintptr
		len int
	}{v.Pointer(), v.Len()}
	if v, ok := seen[ptr]; ok && !checkToNativeOption(SkipCircularReference, ops) {
		return v
	} else if ok {
		return nil
	}
	result := make([]interface{}, len(a.values))
	seen[ptr] = result
	for i := range a.values {
		result[i] = toNative(seen, a.values[i], ops)
	}
	delete(seen, ptr)
	return result
}

func (a Array) Equals(other Value) bool {
	o, ok := other.(Array)
	if !ok {
		return false
	}
	return reflect.DeepEqual(a.values, o.values)
}

func (a Array) SameAs(other Value) bool { return false }

func (a Array) Get(r *Runtime, key Value) Value {
	i := key.ToNumber()
	if i.IsInt() {
		ii := i.ToInt()
		if ii < 0 || ii >= int64(len(a.values)) {
			return Null
		}
		return a.values[ii]
	}

	switch key.ToString() {
	case "length":
		return Int(len(a.values))
	}

	return Null
}

func (a Array) Set(r *Runtime, key, value Value) {
	if !key.IsInt() {
		return
	}
	i := key.ToInt()
	if i < 0 || i >= int64(len(a.values)) {
		return
	}
	a.values[i] = value
}

func (a Array) Iterator() Iterator {
	return &arrayIter{i: 0, a: &a}
}

func (a *arrayIter) Next() (Value, bool) {
	i := a.i
	if i >= 0 && i < len(a.a.values) {
		a.i++
		return a.a.values[i], true
	}
	return Null, false
}

func (a *Array) push(value Value) {
	a.values = append(a.values, value)
}

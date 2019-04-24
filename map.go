package gates

import (
	"math"
	"reflect"
	"sort"
	"unsafe"
)

type Map map[string]Value

type mapIter struct {
	m    Map
	i    int
	keys []string
}

func (Map) Type() string { return "map" }

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

func (m Map) ToNative(ops ...ToNativeOption) interface{} {
	return toNative(nil, m, convertToNativeOption2BinaryOptions(ops))
}

func (m Map) toNative(seen map[unsafe.Pointer]interface{}, options int) interface{} {
	if m == nil {
		return map[string]interface{}(nil)
	}
	addr := unsafe.Pointer(reflect.ValueOf(m).Pointer())
	if v, ok := seen[addr]; ok && !checkToNativeOption(SkipCircularReference, options) {
		return v
	} else if ok {
		return nil
	}
	result := make(map[string]interface{}, len(m))
	seen[addr] = result
	for k, v := range m {
		result[k] = toNative(seen, v, options)
	}
	return result
}

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

func (m Map) Set(r *Runtime, key, value Value) {
	if m == nil {
		return
	}
	m[key.ToString()] = value
}

func (m Map) Iterator() Iterator {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return &mapIter{m: m, i: 0, keys: keys}
}

func (m *mapIter) Next() (Value, bool) {
SkipEmpty:
	i := m.i
	if i >= 0 && i < len(m.keys) {
		m.i++
		k := m.keys[i]
		if _, ok := m.m[k]; !ok {
			goto SkipEmpty
		}
		v := m.m[k]
		return Map(map[string]Value{
			"key":   String(k),
			"value": v,
		}), true
	}
	return Null, false
}

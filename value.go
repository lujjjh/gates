package gates

import (
	"fmt"
)

var intCache [256]Value

type Value interface {
	IsString() bool
	IsInt() bool
	IsFloat() bool
	IsBool() bool
	IsFunction() bool

	ToString() string
	ToInt() int64
	ToFloat() float64
	ToNumber() Number
	ToBool() bool
	ToFunction() Function

	ToNative(...ToNativeOption) interface{}

	Equals(Value) bool
	SameAs(Value) bool
}

type Iterable interface {
	Iterator() Iterator
}

type Iterator interface {
	Next() (value Value, ok bool)
}

func intToValue(i int64) Value {
	if i >= -128 && i <= 127 {
		return intCache[i+128]
	}
	return Int(i)
}

func init() {
	for i := 0; i < 256; i++ {
		intCache[i] = Int(i - 128)
	}
}

func ToValue(i interface{}) Value {
	switch i := i.(type) {
	case nil:
		return Null
	case Value:
		return i
	case string:
		return String(i)
	case bool:
		return Bool(i)
	case int:
		return Int(int64(i))
	case int8:
		return Int(int64(i))
	case int16:
		return Int(int64(i))
	case int32:
		return Int(int64(i))
	case int64:
		return Int(i)
	case uint:
		return Int(int64(i))
	case uint8:
		return Int(int64(i))
	case uint16:
		return Int(int64(i))
	case uint32:
		return Int(int64(i))
	case uint64:
		return Int(int64(i))
	case float32:
		return Float(float64(i))
	case float64:
		return Float(i)
	case map[string]Value:
		return Map(i)
	case []Value:
		return NewArray(i)
	default:
		return Ref{i}
	}
}

func GetIterable(v Value) (Iterable, bool) {
	iter, ok := unref(v).(Iterable)
	return iter, ok
}

func GetIterator(v Value) (Iterator, bool) {
	iterable, ok := GetIterable(v)
	if !ok {
		return nil, false
	}
	return iterable.Iterator(), true
}

type typer interface {
	Type() string
}

// Type returns the type tag of the given value.
func Type(v Value) string {
	switch {
	case v == Null:
		return "null"
	case v.IsBool():
		return "bool"
	case v.IsFloat() || v.IsInt():
		return "number"
	case v.IsFunction():
		return "function"
	case v.IsString():
		return "string"
	}
	if t, haveTyper := unref(v).(typer); haveTyper {
		return t.Type()
	}
	return ""
}

type ErrTypeMismatch struct {
	expected Value
	actual   Value
}

func (e *ErrTypeMismatch) Error() string {
	return fmt.Sprint(Type(e.expected), " expected, got ", Type(e.actual))
}

type ErrTypeNotSupported struct {
	v interface{}
}

func (e *ErrTypeNotSupported) Error() string {
	return fmt.Sprintf("type %T not supported", e.v)
}

func convertValue(r *Runtime, dst interface{}, src Value) error {
	convertArray := func() (result []Value, err error) {
		if src == Null {
			return make([]Value, 0), nil
		}
		if Type(src) != Type(Array{}) {
			return nil, &ErrTypeMismatch{
				expected: Array{},
				actual:   src,
			}
		}
		length := objectGet(r, src, String("length")).ToInt()
		result = make([]Value, 0, length)
		for i := int64(0); i < length; i++ {
			result = append(result, objectGet(r, src, Int(i)))
		}
		return
	}

	convertMap := func() (result map[string]Value, err error) {
		if src == Null {
			return make(map[string]Value), nil
		}
		if Type(src) != Type(Map{}) {
			return nil, &ErrTypeMismatch{
				expected: Map{},
				actual:   src,
			}
		}
		it, ok := GetIterator(src)
		if !ok {
			return nil, &ErrTypeMismatch{
				expected: Map{},
				actual:   src,
			}
		}
		result = make(map[string]Value)
		for {
			elem, ok := it.Next()
			if !ok {
				break
			}
			key := objectGet(r, elem, String("key")).ToString()
			value := objectGet(r, elem, String("value"))
			result[key] = value
		}
		return
	}

	switch dst := dst.(type) {
	case *Value:
		*dst = src
	case *bool:
		*dst = src.ToBool()
	case *Bool:
		*dst = Bool(src.ToBool())
	case *float64:
		*dst = src.ToFloat()
	case *Float:
		*dst = Float(src.ToFloat())
	case *int64:
		*dst = src.ToInt()
	case *Int:
		*dst = Int(src.ToInt())
	case *Callback:
		f := src.ToFunction()
		*dst = func(args ...Value) Value {
			return r.Call(f, args...)
		}
	case *string:
		*dst = src.ToString()
	case *String:
		*dst = String(src.ToString())
	case *[]Value:
		result, err := convertArray()
		if err != nil {
			return err
		}
		*dst = result
	case *Array:
		result, err := convertArray()
		if err != nil {
			return err
		}
		*dst = NewArray(result)
	case *map[string]Value:
		result, err := convertMap()
		if err != nil {
			return err
		}
		*dst = result
	case *Map:
		result, err := convertMap()
		if err != nil {
			return err
		}
		*dst = Map(result)
	default:
		return &ErrTypeNotSupported{v: dst}
	}
	return nil
}

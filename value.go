package gates

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

	ToNative() interface{}

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
		return Array(i)
	default:
		return Ref{i}
	}
}

func GetIterable(v Value) (Iterable, bool) {
	iter, ok := unref(v).(Iterable)
	return iter, ok
}

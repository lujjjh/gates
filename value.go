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

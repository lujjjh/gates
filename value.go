package gates

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

package gates

type Value interface {
	isString() bool
	isInt() bool
	isFloat() bool
	isBool() bool

	ToString() string
	ToInt() int64
	ToFloat() float64
	ToNumber() Number
	ToBool() bool

	SameAs(Value) bool
}

package gates

type Value interface {
	IsString() bool
	IsInt() bool
	IsFloat() bool
	IsBool() bool

	ToString() string
	ToInt() int64
	ToFloat() float64
	ToNumber() Number
	ToBool() bool

	Equals(Value) bool

	SameAs(Value) bool
}

package gates

type _Null struct{}

var Null _Null

func (_Null) IsString() bool   { return false }
func (_Null) IsInt() bool      { return false }
func (_Null) IsFloat() bool    { return false }
func (_Null) IsBool() bool     { return false }
func (_Null) IsFunction() bool { return false }

func (_Null) ToString() string                       { return "" }
func (_Null) ToInt() int64                           { return 0 }
func (_Null) ToFloat() float64                       { return 0 }
func (_Null) ToNumber() Number                       { return Int(0) }
func (_Null) ToBool() bool                           { return false }
func (_Null) ToFunction() Function                   { return _EmptyFunction }
func (_Null) ToNative(...ToNativeOption) interface{} { return nil }

func (_Null) Equals(other Value) bool {
	return other == Null
}

func (_Null) SameAs(other Value) bool {
	return other == Null
}

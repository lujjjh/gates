package gates

type _Null struct{}

var Null _Null

func (n _Null) IsString() bool { return false }
func (n _Null) IsInt() bool    { return false }
func (n _Null) IsFloat() bool  { return false }
func (n _Null) IsBool() bool   { return false }

func (n _Null) ToString() string { return "null" }
func (n _Null) ToInt() int64     { return 0 }
func (n _Null) ToFloat() float64 { return 0 }
func (n _Null) ToNumber() Number { return intNumber(0) }
func (n _Null) ToBool() bool     { return false }

func (n _Null) Equals(other Value) bool {
	if other == Null {
		return true
	}
	return other.ToNumber().Equals(intNumber(0))
}

func (n _Null) SameAs(other Value) bool {
	return other == Null
}

package gates

var (
	True  = Bool(true)
	False = Bool(false)
)

type Bool bool

func (Bool) IsString() bool   { return false }
func (Bool) IsInt() bool      { return false }
func (Bool) IsFloat() bool    { return false }
func (Bool) IsBool() bool     { return true }
func (Bool) IsFunction() bool { return false }

func (b Bool) ToString() string {
	if bool(b) {
		return "true"
	}
	return "false"
}

func (b Bool) ToInt() int64 {
	if bool(b) {
		return 1
	}
	return 0
}

func (b Bool) ToFloat() float64                       { return float64(b.ToInt()) }
func (b Bool) ToNumber() Number                       { return Int(b.ToInt()) }
func (b Bool) ToBool() bool                           { return bool(b) }
func (b Bool) ToFunction() Function                   { return _EmptyFunction }
func (b Bool) ToNative(...ToNativeOption) interface{} { return b.ToBool() }

func (b Bool) Equals(other Value) bool {
	if other.IsBool() {
		return b.SameAs(other)
	}
	return other.Equals(Int(b.ToInt()))
}

func (b Bool) SameAs(bv Value) bool {
	bb, ok := bv.(Bool)
	if !ok {
		return false
	}
	return bool(bb) == bool(b)
}

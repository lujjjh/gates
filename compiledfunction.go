package gates

type CompiledFunction struct {
	Instructions  []byte
	NumLocals     int
	NumParameters int
}

func (f *CompiledFunction) function() {}

func (f *CompiledFunction) IsString() bool   { return false }
func (f *CompiledFunction) IsInt() bool      { return false }
func (f *CompiledFunction) IsFloat() bool    { return false }
func (f *CompiledFunction) IsBool() bool     { return false }
func (f *CompiledFunction) IsFunction() bool { return true }

func (f *CompiledFunction) ToString() string     { return "<function>" }
func (f *CompiledFunction) ToInt() int64         { return 0 }
func (f *CompiledFunction) ToFloat() float64     { return 0 }
func (f *CompiledFunction) ToNumber() Number     { return Int(0) }
func (f *CompiledFunction) ToBool() bool         { return true }
func (f *CompiledFunction) ToFunction() Function { return f }

func (f *CompiledFunction) ToNative(...ToNativeOption) interface{} { return nil }

func (f *CompiledFunction) Equals(Value) bool { return false }
func (f *CompiledFunction) SameAs(Value) bool { return false }

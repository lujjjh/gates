package gates

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

type String string

func (String) IsString() bool   { return true }
func (String) IsInt() bool      { return false }
func (String) IsFloat() bool    { return false }
func (String) IsBool() bool     { return false }
func (String) IsFunction() bool { return false }

func (s String) ToString() string { return string(s) }

func (s String) ToInt() int64 {
	i, _ := strconv.ParseInt(string(s), 0, 64)
	return i
}

func (s String) ToFloat() float64 {
	f, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func (s String) ToNumber() Number {
	t := strings.TrimSpace(string(s))
	i, err := strconv.ParseInt(t, 0, 64)
	if err == nil {
		return Int(i)
	}
	return Float(s.ToFloat())
}

func (s String) ToBool() bool                           { return string(s) != "" }
func (s String) ToFunction() Function                   { return _EmptyFunction }
func (s String) ToNative(...ToNativeOption) interface{} { return string(s) }

func (s String) Equals(other Value) bool {
	switch {
	case other.IsString():
		return s.SameAs(other)
	case other.IsInt(), other.IsFloat(), other.IsBool():
		return other.Equals(other)
	default:
		return false
	}
}

func (s String) SameAs(b Value) bool {
	bs, ok := b.(String)
	if !ok {
		return false
	}
	return string(bs) == string(s)
}

func (s String) Get(r *Runtime, key Value) Value {
	switch {
	case key.IsInt():
		index := int(key.ToInt())
		if index < 0 {
			return Null
		}
		i := 0
		start := -1
		for j := range string(s) {
			if i == index {
				start = j
			}
			if i == index+1 {
				return String(string(s)[start:j])
			}
			i++
		}
		if start == -1 {
			return Null
		}
		return String(string(s)[start:])
	case key.IsString():
		switch key.ToString() {
		case "length":
			return Int(int64(utf8.RuneCountInString(string(s))))
		}
	}

	return Null
}

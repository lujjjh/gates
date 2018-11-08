package gates

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

type _String struct{ s string }

func String(s string) _String { return _String{s} }

func (s _String) IsString() bool { return true }
func (s _String) IsInt() bool    { return false }
func (s _String) IsFloat() bool  { return false }
func (s _String) IsBool() bool   { return false }

func (s _String) ToString() string { return s.s }

func (s _String) ToInt() int64 {
	i, _ := strconv.ParseInt(s.s, 0, 64)
	return i
}

func (s _String) ToFloat() float64 {
	f, err := strconv.ParseFloat(s.s, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func (s _String) ToNumber() Number {
	t := strings.TrimSpace(s.s)
	i, err := strconv.ParseInt(t, 0, 64)
	if err == nil {
		return Int(i)
	}
	return Float(s.ToFloat())
}

func (s _String) ToBool() bool { return s.s != "" }

func (s _String) Equals(other Value) bool {
	switch {
	case other.IsString():
		return s.SameAs(other)
	case other.IsInt(), other.IsFloat(), other.IsBool():
		return s.ToNumber().Equals(other)
	default:
		return false
	}
}

func (s _String) SameAs(b Value) bool {
	bs, ok := b.(_String)
	if !ok {
		return false
	}
	return bs.s == s.s
}

func (s _String) Get(r *Runtime, key Value) Value {
	switch {
	case key.IsInt():
		index := int(key.ToInt())
		if index < 0 {
			return Null
		}
		i := 0
		start := -1
		for j := range s.s {
			if i == index {
				start = j
			}
			if i == index+1 {
				return String(string(s.s[start:j]))
			}
			i++
		}
		if start == -1 {
			return Null
		}
		return String(string(s.s[start:]))
	case key.IsString():
		switch key.ToString() {
		case "length":
			return Int(int64(utf8.RuneCountInString(string(s.s))))
		}
	}

	return Null
}

package gates

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

type String string

func (s String) IsString() bool { return true }
func (s String) IsInt() bool    { return false }
func (s String) IsFloat() bool  { return false }
func (s String) IsBool() bool   { return false }

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

func (s String) ToBool() bool { return string(s) != "" }

func (s String) Equals(other Value) bool {
	switch {
	case other.IsString():
		return s.SameAs(other)
	case other.IsInt(), other.IsFloat(), other.IsBool():
		return s.ToNumber().Equals(other)
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
		index := key.ToInt()
		if index < 0 {
			return Null
		}
		for _, r := range string(s) {
			if index == 0 {
				return String(string(r))
			}
			index--
		}
		return Null
	case key.IsString():
		switch key.ToString() {
		case "length":
			return Int(int64(utf8.RuneCountInString(string(s))))
		}
	}

	return Null
}

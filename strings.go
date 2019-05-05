package gates

import (
	"regexp"
	"strings"
	"unicode"
)

func packageStrings() Map {
	return Map{
		"has_prefix": FunctionFunc(func(fc FunctionCall) Value {
			var s, prefix string
			if err := NewArgumentScanner(fc).Scan(&s, &prefix); err != nil {
				return False
			}
			return Bool(strings.HasPrefix(s, prefix))
		}),

		"has_suffix": FunctionFunc(func(fc FunctionCall) Value {
			var s, suffix string
			if err := NewArgumentScanner(fc).Scan(&s, &suffix); err != nil {
				return False
			}
			return Bool(strings.HasSuffix(s, suffix))
		}),

		"to_lower": FunctionFunc(func(fc FunctionCall) Value {
			var s string
			if err := NewArgumentScanner(fc).Scan(&s); err != nil {
				return Null
			}
			return String(strings.ToLower(s))
		}),

		"to_upper": FunctionFunc(func(fc FunctionCall) Value {
			var s string
			if err := NewArgumentScanner(fc).Scan(&s); err != nil {
				return Null
			}
			return String(strings.ToUpper(s))
		}),

		"trim": FunctionFunc(func(fc FunctionCall) Value {
			var s, cutset string
			scanner := NewArgumentScanner(fc)
			if err := scanner.Scan(&s); err != nil {
				return Null
			}
			if err := scanner.Scan(&cutset); err != nil {
				return String(strings.TrimSpace(s))
			}
			return String(strings.Trim(s, cutset))
		}),

		"trim_left": FunctionFunc(func(fc FunctionCall) Value {
			var s, cutset string
			scanner := NewArgumentScanner(fc)
			if err := scanner.Scan(&s); err != nil {
				return Null
			}
			if err := scanner.Scan(&cutset); err != nil {
				return String(strings.TrimLeftFunc(s, unicode.IsSpace))
			}
			return String(strings.TrimLeft(s, cutset))
		}),

		"trim_right": FunctionFunc(func(fc FunctionCall) Value {
			var s, cutset string
			scanner := NewArgumentScanner(fc)
			if err := scanner.Scan(&s); err != nil {
				return Null
			}
			if err := scanner.Scan(&cutset); err != nil {
				return String(strings.TrimRightFunc(s, unicode.IsSpace))
			}
			return String(strings.TrimRight(s, cutset))
		}),

		"split": FunctionFunc(func(fc FunctionCall) Value {
			var s, sep string
			scanner := NewArgumentScanner(fc)
			if err := scanner.Scan(&s); err != nil {
				return Null
			}
			if err := scanner.Scan(&sep); err != nil {
				return NewArray([]Value{String(s)})
			}
			return NewArrayFromStringSlice(strings.Split(s, sep))
		}),

		"join": FunctionFunc(func(fc FunctionCall) Value {
			var a []Value
			var sep string
			if err := NewArgumentScanner(fc).Scan(&a, &sep); err != nil {
				return Null
			}
			as := make([]string, len(a))
			for i := range a {
				as[i] = a[i].ToString()
			}
			return String(strings.Join(as, sep))
		}),

		"match": FunctionFunc(func(fc FunctionCall) Value {
			var expr, s string
			if err := NewArgumentScanner(fc).Scan(&expr, &s); err != nil {
				return Null
			}
			re, err := regexp.Compile(expr)
			if err != nil {
				return Null
			}
			result := re.FindStringSubmatch(s)
			if len(result) == 0 {
				return Null
			}
			m := newMatcher(result, re.SubexpNames())
			return Map{"group": FunctionFunc(m.group)}
		}),

		"find_all": FunctionFunc(func(fc FunctionCall) Value {
			var expr, s string
			if err := NewArgumentScanner(fc).Scan(&expr, &s); err != nil {
				return Null
			}
			re, err := regexp.Compile(expr)
			if err != nil {
				return Null
			}
			return NewArrayFromStringSlice(re.FindAllString(s, -1))
		}),

		"contains": FunctionFunc(func(fc FunctionCall) Value {
			var s, substr string
			if err := NewArgumentScanner(fc).Scan(&s, &substr); err != nil {
				return Null
			}
			return Bool(strings.Contains(s, substr))
		}),

		"contains_any": FunctionFunc(func(fc FunctionCall) Value {
			var s, chars string
			if err := NewArgumentScanner(fc).Scan(&s, &chars); err != nil {
				return Null
			}
			return Bool(strings.ContainsAny(s, chars))
		}),

		"index": FunctionFunc(func(fc FunctionCall) Value {
			var s, substr string
			if err := NewArgumentScanner(fc).Scan(&s, &substr); err != nil {
				return Null
			}
			return Int(strings.Index(s, substr))
		}),

		"index_any": FunctionFunc(func(fc FunctionCall) Value {
			var s, chars string
			if err := NewArgumentScanner(fc).Scan(&s, &chars); err != nil {
				return Null
			}
			return Int(strings.IndexAny(s, chars))
		}),

		"last_index": FunctionFunc(func(fc FunctionCall) Value {
			var s, substr string
			if err := NewArgumentScanner(fc).Scan(&s, &substr); err != nil {
				return Null
			}
			return Int(strings.LastIndex(s, substr))
		}),

		"last_index_any": FunctionFunc(func(fc FunctionCall) Value {
			var s, chars string
			if err := NewArgumentScanner(fc).Scan(&s, &chars); err != nil {
				return Null
			}
			return Int(strings.LastIndexAny(s, chars))
		}),

		"repeat": FunctionFunc(func(fc FunctionCall) Value {
			var s string
			var count64 int64
			if err := NewArgumentScanner(fc).Scan(&s, &count64); err != nil {
				return Null
			}
			count := int(count64)
			if count < 0 {
				return Null
			}
			if count > 0 && len(s)*count/count != len(s) {
				return Null
			}
			return String(strings.Repeat(s, count))
		}),
	}
}

type matcher struct {
	result     []string
	nameResult map[string]string
}

func (m matcher) group(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	}
	if args[0].IsInt() {
		index := int(args[0].ToInt())
		if index >= len(m.result) || index < 0 {
			return Null
		}
		return String(m.result[index])
	}
	res, exist := m.nameResult[args[0].ToString()]
	if !exist {
		return Null
	}
	return String(res)
}

func newMatcher(result []string, expNames []string) matcher {
	m := matcher{
		result:     result,
		nameResult: make(map[string]string),
	}
	if len(expNames) > 0 {
		for i, name := range expNames {
			if i != 0 && name != "" {
				m.nameResult[name] = result[i]
			}
		}
	}
	return m
}

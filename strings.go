package gates

import (
	"regexp"
	"strings"
	"unicode"
)

type packageStrings struct {
	r *Runtime
}

func (s packageStrings) export() Map {
	ps := map[string]func(FunctionCall) Value{
		"has_prefix":     s.hasPrefix,
		"has_suffix":     s.hasSuffix,
		"to_lower":       s.toLower,
		"to_upper":       s.toUpper,
		"trim":           s.trim,
		"trim_left":      s.trimLeft,
		"trim_right":     s.trimRight,
		"split":          s.split,
		"join":           s.join,
		"match":          s.match,
		"find_all":       s.findAll,
		"contains":       s.contains,
		"contains_any":   s.containsAny,
		"index":          s.index,
		"index_any":      s.indexAny,
		"last_index":     s.lastIndex,
		"last_index_any": s.lastIndexAny,
		"repeat":         s.repeat,
	}
	m := make(Map, len(ps))
	for name, fun := range ps {
		m[name] = FunctionFunc(fun)
	}
	return m
}

func initPackageStrings(r *Runtime) {
	ps := packageStrings{r}
	r.Global().Set("strings", ps.export())
}

func (s packageStrings) hasPrefix(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return False
	}
	return Bool(strings.HasPrefix(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) hasSuffix(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return False
	}
	return Bool(strings.HasSuffix(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) toLower(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	}
	return String(strings.ToLower(args[0].ToString()))
}

func (s packageStrings) toUpper(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	}
	return String(strings.ToUpper(args[0].ToString()))
}

func (s packageStrings) trim(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	} else if len(args) == 1 {
		return String(strings.TrimSpace(args[0].ToString()))
	}
	return String(strings.Trim(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) trimLeft(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	} else if len(args) == 1 {
		return String(strings.TrimLeftFunc(args[0].ToString(), unicode.IsSpace))
	}
	return String(strings.TrimLeft(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) trimRight(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	} else if len(args) == 1 {
		return String(strings.TrimRightFunc(args[0].ToString(), unicode.IsSpace))
	}
	return String(strings.TrimRight(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) split(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	} else if len(args) == 1 {
		return NewArray([]Value{args[0]})
	}

	values := strings.Split(args[0].ToString(), args[1].ToString())
	result := make([]Value, len(values))
	for i, value := range values {
		result[i] = String(value)
	}
	return NewArray(result)
}

func (s packageStrings) join(fc FunctionCall) Value {
	args := fc.Args()
	var sep string
	if len(args) == 0 {
		return Null
	} else if len(args) > 1 {
		sep = args[1].ToString()
	}

	get, ok := args[0].(getter)
	if !ok {
		return Null
	}

	length := get.Get(s.r, String("length"))

	a := make([]string, 0, length.ToInt())

	var i int64
	for ; i < length.ToInt(); i++ {
		a = append(a, get.Get(s.r, Int(i)).ToString())
	}

	return String(strings.Join(a, sep))
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

func (s packageStrings) match(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	pattern, err := regexp.Compile(args[0].ToString())
	if err != nil {
		return Null
	}
	str := args[1].ToString()
	result := pattern.FindStringSubmatch(str)
	if len(result) == 0 {
		return Null
	}
	m := newMatcher(result, pattern.SubexpNames())
	return Map{"group": FunctionFunc(m.group)}
}

func (packageStrings) findAll(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}
	pattern, err := regexp.Compile(args[0].ToString())
	if err != nil {
		return Null
	}
	str := args[1].ToString()
	results := pattern.FindAllString(str, -1)

	result := make([]Value, len(results))
	for i, value := range results {
		result[i] = String(value)
	}
	return NewArray(result)
}

func (packageStrings) contains(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	s := args[0].ToString()
	substr := args[1].ToString()
	return Bool(strings.Contains(s, substr))
}

func (packageStrings) containsAny(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	s := args[0].ToString()
	substr := args[1].ToString()
	return Bool(strings.ContainsAny(s, substr))
}

func (packageStrings) index(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	s := args[0].ToString()
	substr := args[1].ToString()
	return Int(strings.Index(s, substr))
}

func (packageStrings) indexAny(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	s := args[0].ToString()
	substr := args[1].ToString()
	return Int(strings.IndexAny(s, substr))
}

func (packageStrings) lastIndex(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	s := args[0].ToString()
	substr := args[1].ToString()
	return Int(strings.LastIndex(s, substr))
}

func (packageStrings) lastIndexAny(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	s := args[0].ToString()
	substr := args[1].ToString()
	return Int(strings.LastIndexAny(s, substr))
}

func (packageStrings) repeat(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}

	s := args[0].ToString()
	count := int(args[1].ToInt())
	if count < 0 {
		return Null
	}
	if count > 0 && len(s)*count/count != len(s) {
		return Null
	}
	return String(strings.Repeat(s, count))
}

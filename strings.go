package gates

import (
	"strings"
)

type packageStrings struct {
	r *Runtime
}

func (s packageStrings) export() Map {
	ps := map[string]func(FunctionCall) Value{
		"has_prefix": s.hasPrefix,
		"has_suffix": s.hasSuffix,
		"to_lower":   s.toLower,
		"to_upper":   s.toUpper,
		"trim":       s.trim,
		"trim_left":  s.trimLeft,
		"trim_right": s.trimRight,
		"trim_space": s.trimSpace,
		"split":      s.split,
		"join":       s.join,
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
		return args[0]
	}
	return String(strings.Trim(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) trimLeft(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	} else if len(args) == 1 {
		return args[0]
	}
	return String(strings.TrimLeft(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) trimRight(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	} else if len(args) == 1 {
		return args[0]
	}
	return String(strings.TrimRight(args[0].ToString(), args[1].ToString()))
}

func (s packageStrings) trimSpace(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	}
	return String(strings.TrimSpace(args[0].ToString()))
}

func (s packageStrings) split(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Null
	} else if len(args) == 1 {
		return Array([]Value{args[0]})
	}

	values := strings.Split(args[0].ToString(), args[1].ToString())
	result := make([]Value, len(values))
	for i, value := range values {
		result[i] = String(value)
	}
	return Array(result)
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

package gates

var builtInFunctions = map[string]Function{
	"bool": FunctionFunc(func(fc FunctionCall) Value {
		var v Bool
		if NewArgumentScanner(fc).Scan(&v) != nil {
			return False
		}
		return v
	}),

	"int": FunctionFunc(func(fc FunctionCall) Value {
		var v Int
		if NewArgumentScanner(fc).Scan(&v) != nil {
			return Int(0)
		}
		return v
	}),

	"number": FunctionFunc(func(fc FunctionCall) Value {
		var v Value
		if NewArgumentScanner(fc).Scan(&v) != nil {
			return Int(0)
		}
		return v.ToNumber()
	}),

	"string": FunctionFunc(func(fc FunctionCall) Value {
		var v String
		if NewArgumentScanner(fc).Scan(&v) != nil {
			return String("")
		}
		return v
	}),

	"type": FunctionFunc(func(fc FunctionCall) Value {
		var v Value
		if NewArgumentScanner(fc).Scan(&v) != nil {
			return String("")
		}
		return String(Type(v))
	}),

	"curry": CurriedFunctionFunc(2, func(fc FunctionCall) Value {
		var n int64
		var f Value
		if NewArgumentScanner(fc).Scan(&n, &f) != nil {
			return Null
		}
		return Curry(f.ToFunction(), int(n))
	}),

	"map": CurriedFunctionFunc(2, func(fc FunctionCall) Value {
		var f Callback
		var base []Value
		if NewArgumentScanner(fc).Scan(&f, &base) != nil {
			return Null
		}
		result := make([]Value, len(base))
		for i := 0; i < len(base); i++ {
			result[i] = f(base[i], Int(i))
		}
		return NewArray(result)
	}),

	"filter": CurriedFunctionFunc(2, func(fc FunctionCall) Value {
		var f Callback
		var base []Value
		if NewArgumentScanner(fc).Scan(&f, &base) != nil {
			return Null
		}
		result := make([]Value, 0)
		for i := 0; i < len(base); i++ {
			if f(base[i], Int(i)).ToBool() {
				result = append(result, base[i])
			}
		}
		return NewArray(result)
	}),

	"reduce": CurriedFunctionFunc(3, func(fc FunctionCall) Value {
		var f Callback
		var initial Value
		var base Value
		if NewArgumentScanner(fc).Scan(&f, &initial, &base) != nil {
			return Null
		}
		var baseArray []Value
		if convertValue(fc.Runtime(), &baseArray, base) != nil {
			return Null
		}
		acc := initial
		for i := 0; i < len(baseArray); i++ {
			acc = f(acc, baseArray[i], Int(i), base)
		}
		return acc
	}),

	"find": CurriedFunctionFunc(2, func(fc FunctionCall) Value {
		var f Callback
		var base []Value
		if NewArgumentScanner(fc).Scan(&f, &base) != nil {
			return Null
		}
		for i := 0; i < len(base); i++ {
			if f(base[i], Int(i)).ToBool() {
				return base[i]
			}
		}
		return Null
	}),

	"find_index": CurriedFunctionFunc(2, func(fc FunctionCall) Value {
		var f Callback
		var base []Value
		if NewArgumentScanner(fc).Scan(&f, &base) != nil {
			return Int(-1)
		}
		for i := 0; i < len(base); i++ {
			if f(base[i], Int(i)).ToBool() {
				return Int(i)
			}
		}
		return Int(-1)
	}),

	"find_last": CurriedFunctionFunc(2, func(fc FunctionCall) Value {
		var f Callback
		var base []Value
		if NewArgumentScanner(fc).Scan(&f, &base) != nil {
			return Null
		}
		for i := len(base) - 1; i >= 0; i-- {
			if f(base[i], Int(i)).ToBool() {
				return base[i]
			}
		}
		return Null
	}),

	"find_last_index": CurriedFunctionFunc(2, func(fc FunctionCall) Value {
		var f Callback
		var base []Value
		if NewArgumentScanner(fc).Scan(&f, &base) != nil {
			return Int(-1)
		}
		for i := len(base) - 1; i >= 0; i-- {
			if f(base[i], Int(i)).ToBool() {
				return Int(i)
			}
		}
		return Int(-1)
	}),

	"to_entries": FunctionFunc(func(fc FunctionCall) Value {
		var v Value
		if NewArgumentScanner(fc).Scan(&v) != nil {
			return Null
		}
		it, ok := GetIterator(v)
		if !ok {
			return Null
		}
		entries := make([]Value, 0)
		for {
			value, ok := it.Next()
			if !ok {
				break
			}
			entries = append(entries, value)
		}
		return NewArray(entries)
	}),

	"from_entries": FunctionFunc(func(fc FunctionCall) Value {
		var v Value
		if NewArgumentScanner(fc).Scan(&v) != nil {
			return Null
		}
		it, ok := GetIterator(v)
		if !ok {
			return Null
		}
		result := make(map[string]Value)
		r := fc.Runtime()
		for {
			entry, ok := it.Next()
			if !ok {
				break
			}
			k := objectGet(r, entry, String("key"))
			v := objectGet(r, entry, String("value"))
			result[k.ToString()] = v
		}
		return Map(result)
	}),
}

type Global struct {
	m Map
}

func NewGlobal() *Global {
	return &Global{
		m: make(Map),
	}
}

func (g *Global) Set(name string, value Value) {
	g.m[name] = value
}

func (g *Global) Get(name string) Value {
	return g.m[name]
}

func Curry(f Function, n int) Function {
	return FunctionFunc(func(fc FunctionCall) Value {
		var curriedF Function
		curriedF = FunctionFunc(func(fc FunctionCall) Value {
			args := fc.Args()
			argc := len(args)
			if argc >= n {
				return fc.Runtime().Call(f, args...)
			}
			return Curry(FunctionFunc(func(fc FunctionCall) Value {
				return fc.Runtime().Call(f, append(args, fc.Args()...)...)
			}), n-argc)
		})
		return curriedF
	})
}

func CurriedFunctionFunc(n int, f func(FunctionCall) Value) Function {
	return Curry(FunctionFunc(f), n)
}

func (g *Global) initBuiltInFunctions() {
	for name, f := range builtInFunctions {
		g.Set(name, f)
	}
}

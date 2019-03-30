package gates

func curry(f Function, n int) Function {
	return FunctionFunc(func(fc FunctionCall) Value {
		curriedArgs := make([]Value, 0, n)
		var curriedF Function
		curriedF = FunctionFunc(func(fc FunctionCall) Value {
			args := fc.Args()
			curriedArgs = append(curriedArgs, args...)
			if len(curriedArgs) < n {
				return curriedF
			}
			return fc.Runtime().Call(f, curriedArgs...)
		})
		return fc.Runtime().Call(curriedF, fc.Args()...)
	})
}

func builtInBool(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return False
	}
	return Bool(args[0].ToBool())
}

func builtInInt(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Int(0)
	}
	return Int(args[0].ToInt())
}

func builtInNumber(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Int(0)
	}
	return args[0].ToNumber()
}

func builtInString(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return String("")
	}
	return String(args[0].ToString())
}

func builtInCurry(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) != 2 {
		return Null
	}
	n, f := int(args[0].ToInt()), args[1].ToFunction()
	return curry(f, n)
}

func builtInMap(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}
	r := fc.Runtime()
	f, base := args[0].ToFunction(), args[1]
	length := int(objectGet(r, base, String("length")).ToInt())
	result := make([]Value, length)
	for i := 0; i < length; i++ {
		result[i] = r.Call(f, objectGet(r, base, Int(i)), Int(i))
	}
	return Array(result)
}

func builtInFilter(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}
	r := fc.Runtime()
	f, base := args[0].ToFunction(), args[1]
	length := int(objectGet(r, base, String("length")).ToInt())
	result := make([]Value, 0, length)
	for i := 0; i < length; i++ {
		v := objectGet(r, base, Int(i))
		if r.Call(f, v, Int(i)).ToBool() {
			result = append(result, v)
		}
	}
	return Array(result)
}

func builtInReduce(fc FunctionCall) Value {
	args := fc.Args()
	argc := len(args)
	if argc < 3 {
		return Null
	}
	r := fc.Runtime()
	f, initial, base := args[0].ToFunction(), args[1], args[2]
	length := int(objectGet(r, base, String("length")).ToInt())
	acc := initial
	for i := 0; i < length; i++ {
		v := objectGet(r, base, Int(i))
		acc = r.Call(f, acc, v, Int(i), base)
	}
	return acc
}

func builtInFind(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}
	r := fc.Runtime()
	f, base := args[0].ToFunction(), args[1]
	length := int(objectGet(r, base, String("length")).ToInt())
	for i := 0; i < length; i++ {
		v := objectGet(r, base, Int(i))
		if r.Call(f, v, Int(i)).ToBool() {
			return v
		}
	}
	return Null
}

func builtInFindIndex(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}
	r := fc.Runtime()
	f, base := args[0].ToFunction(), args[1]
	length := int(objectGet(r, base, String("length")).ToInt())
	for i := 0; i < length; i++ {
		v := objectGet(r, base, Int(i))
		if r.Call(f, v, Int(i)).ToBool() {
			return Int(i)
		}
	}
	return Int(-1)
}

func builtInFindLast(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}
	r := fc.Runtime()
	f, base := args[0].ToFunction(), args[1]
	length := int(objectGet(r, base, String("length")).ToInt())
	for i := length - 1; i >= 0; i-- {
		v := objectGet(r, base, Int(i))
		if r.Call(f, v, Int(i)).ToBool() {
			return v
		}
	}
	return Null
}

func builtInFindLastIndex(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 2 {
		return Null
	}
	r := fc.Runtime()
	f, base := args[0].ToFunction(), args[1]
	length := int(objectGet(r, base, String("length")).ToInt())
	for i := length - 1; i >= 0; i-- {
		v := objectGet(r, base, Int(i))
		if r.Call(f, v, Int(i)).ToBool() {
			return Int(i)
		}
	}
	return Int(-1)
}

func builtInToEntries(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 1 {
		return Null
	}
	iterable, ok := GetIterable(args[0])
	if !ok {
		return Null
	}
	entries := make([]Value, 0)
	it := iterable.Iterator()
	for {
		value, ok := it.Next()
		if !ok {
			break
		}
		entries = append(entries, value)
	}
	return Array(entries)
}

func builtInFromEntries(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) < 1 {
		return Null
	}
	r := fc.Runtime()
	iterable, ok := GetIterable(args[0])
	if !ok {
		return Null
	}
	result := make(map[string]Value)
	it := iterable.Iterator()
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
}

package gates

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
	if argc < 2 {
		return Null
	}
	r := fc.Runtime()
	f, base := args[0].ToFunction(), args[1]
	length := int(objectGet(r, base, String("length")).ToInt())
	initial := Value(Null)
	i := 1
	if argc >= 3 {
		initial = args[2]
		i = 0
	} else if length > 0 {
		initial = objectGet(r, base, Int(0))
	}
	acc := initial
	for ; i < length; i++ {
		v := objectGet(r, base, Int(i))
		acc = r.Call(f, acc, v, Int(i), base)
	}
	return acc
}

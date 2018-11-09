package gates

func builtInBool(fc FunctionCall) Value {
	args := fc.Args()
	if len(args) == 0 {
		return Bool(false)
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

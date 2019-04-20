package gates

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

func (g *Global) initBuiltIns() {
	g.Set("bool", FunctionFunc(builtInBool))
	g.Set("int", FunctionFunc(builtInInt))
	g.Set("number", FunctionFunc(builtInNumber))
	g.Set("string", FunctionFunc(builtInString))
	g.Set("type", FunctionFunc(builtInType))

	g.Set("curry", curry(FunctionFunc(builtInCurry), 2))

	g.Set("map", curry(FunctionFunc(builtInMap), 2))
	g.Set("filter", curry(FunctionFunc(builtInFilter), 2))
	g.Set("reduce", curry(FunctionFunc(builtInReduce), 3))
	g.Set("find", curry(FunctionFunc(builtInFind), 2))
	g.Set("find_index", curry(FunctionFunc(builtInFindIndex), 2))
	g.Set("find_last", curry(FunctionFunc(builtInFindLast), 2))
	g.Set("find_last_index", curry(FunctionFunc(builtInFindLastIndex), 2))

	g.Set("to_entries", FunctionFunc(builtInToEntries))
	g.Set("from_entries", FunctionFunc(builtInFromEntries))
}

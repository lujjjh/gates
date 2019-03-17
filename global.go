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

	g.Set("map", FunctionFunc(builtInMap))
	g.Set("filter", FunctionFunc(builtInFilter))
	g.Set("reduce", FunctionFunc(builtInReduce))
	g.Set("find", FunctionFunc(builtInFind))
	g.Set("find_index", FunctionFunc(builtInFindIndex))
	g.Set("find_last", FunctionFunc(builtInFindLast))
	g.Set("find_last_index", FunctionFunc(builtInFindLastIndex))

	g.Set("to_entries", FunctionFunc(builtInToEntries))
	g.Set("from_entries", FunctionFunc(builtInFromEntries))
}

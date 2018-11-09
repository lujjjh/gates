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

func (g *Global) InitBuiltIns() {
	g.Set("bool", FunctionFunc(builtInBool))
	g.Set("int", FunctionFunc(builtInInt))
	g.Set("number", FunctionFunc(builtInNumber))
	g.Set("string", FunctionFunc(builtInString))
}

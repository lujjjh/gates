package gates

type getter interface {
	Get(*Runtime, Value) Value
}

type setter interface {
	Set(*Runtime, Value, Value)
}

type getterFunc func(*Runtime, Value) Value

func (g getterFunc) Get(r *Runtime, v Value) Value { return g(r, v) }

func ObjectGet(r *Runtime, base interface{}, key Value) Value {
	base = unref(base)
	g, ok := base.(getter)
	if !ok {
		m, ok := base.(map[string]Value)
		if !ok {
			return Null
		}
		g = Map(m)
	}
	return g.Get(r, key)
}

func ObjectSet(r *Runtime, base interface{}, key, value Value) {
	base = unref(base)
	g, ok := base.(setter)
	if !ok {
		return
	}
	g.Set(r, key, value)
	return
}

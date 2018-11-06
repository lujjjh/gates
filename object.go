package gates

type getter interface {
	Get(*Runtime, Value) Value
}

type getterFunc func(*Runtime, Value) Value

func (g getterFunc) Get(r *Runtime, v Value) Value { return g(r, v) }

func objectGet(r *Runtime, base interface{}, key Value) Value {
	base = unref(base)
	g, ok := base.(getter)
	if !ok {
		m, ok := base.(map[string]interface{})
		if !ok {
			return Null
		}
		g = Map(m)
	}
	return g.Get(r, key)
}

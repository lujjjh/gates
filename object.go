package gates

type getter interface {
	Get(*Runtime, Value) Value
}

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

package gates

type Program struct {
	code   []instruction
	values []Value
}

func (p *Program) defineLit(v Value) uint {
	for index, value := range p.values {
		if value.SameAs(v) {
			return uint(index)
		}
	}
	index := uint(len(p.values))
	p.values = append(p.values, v)
	return index
}

func (p *Program) emit(instructions ...instruction) {
	p.code = append(p.code, instructions...)
}

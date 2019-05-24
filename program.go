package gates

import "github.com/gates/gates/syntax"

type Program struct {
	src    *syntax.File
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

func (p *Program) InstructionNumber() int {
	return len(p.code)
}

package compiler

type scope struct {
	names   map[string]uint32
	visited bool

	outer *scope
}

func newScope(outer *scope) *scope {
	s := &scope{}
	s.init(outer)
	return s
}

func (s *scope) init(outer *scope) {
	s.names = make(map[string]uint32)
	s.outer = outer
}

func (s *scope) lookupName(name string) (uint32, bool) {
	level := uint32(0)
	for current := s; current != nil; current = current.outer {
		if current != s {
			current.visited = true
		}
		if i, ok := current.names[name]; ok {
			return i | (level << 24), true
		}
		level++
	}
	return 0, false
}

func (s *scope) bindName(name string) uint32 {
	if idx, ok := s.names[name]; ok {
		return idx
	}
	idx := uint32(len(s.names))
	s.names[name] = idx
	return idx
}

package vm

// State is an opaque structure representing Gates state.
type State struct {
	top   int       // first free slot in the stack
	ci    *CallInfo // call info for current function
	oldPC int
}

package vm

// CallInfo is a struct representing the call information.
type CallInfo struct {
	function, top  int
	previous, next *CallInfo
}

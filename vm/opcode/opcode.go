package opcode

// Opcode is the portion of Gates VM instruction that
// specifies the operation to be performed.
type Opcode int

// Opcode could be one of the values below:
const (
	Move Opcode = iota
	LoadK
	LoadBool
	LoadNull
	GetUpVal

	GetGlobal
	GetTable

	SetGlobal
	SetUpVal
	SetTable

	NewMap
	NewArray

	Add
	Sub
	Mul
	Div
	Pow
	Unm
	Not
)

package opcode

const (
	sizeC = 20
	sizeB = 20
	sizeA = 16

	sizeOp = 8

	posOp = 0
	posA  = posOp + sizeOp
	posC  = posA + sizeA
	posB  = posC + sizeC

	maxArgA = 1<<sizeA - 1
	maxArgB = 1<<sizeB - 1
	maxArgC = 1<<sizeC - 1

	bitRK      = 1 << (sizeB - 1)
	maxIndexRK = bitRK - 1
)

func isK(x int) bool   { return x&bitRK != 0 }
func indexK(x int) int { return x & ^bitRK }
func rkAsK(x int) int  { return x | bitRK }

// creates a mask with 'n' 1 bits at position 'p'
func mask1(n, p uint) Instruction { return ^(^Instruction(0) << n) << p }

// creates a mask with 'n' 0 bits at position 'p'
func mask0(n, p uint) Instruction { return ^mask1(n, p) }

// Instruction is a Gates VM instruction. The instruction
// layout is similar to Lua; however, a Gates VM instruction
// occupies 64 bits instead of 32 bits:
//
//     +--------------+----------+-----------+-----------+
//     | opcode (0-7) | A (8-23) | B (24-43) | C (44-63) |
//     +--------------+----------+-----------+-----------+
type Instruction uint64

func (i Instruction) arg(pos, size uint) int { return int(i >> pos & mask1(size, 0)) }

func (i *Instruction) setArg(pos, size uint, arg int) {
	*i = *i&mask0(size, pos) | Instruction(arg)<<pos&mask1(size, pos)
}

// Opcode returns the opcode of the instruction.
func (i Instruction) Opcode() Opcode { return Opcode(i >> posOp & (1<<sizeOp - 1)) }

// SetOpcode updates the opcode of the instruction.
func (i *Instruction) SetOpcode(op Opcode) { i.setArg(posOp, sizeOp, int(op)) }

// A returns the A of the instruction.
func (i Instruction) A() int { return int(i >> posA & maxArgA) }

// SetA updates the A of the instruction.
func (i *Instruction) SetA(x int) { i.setArg(posA, sizeA, x) }

// B returns the B of the instruction.
func (i Instruction) B() int { return int(i >> posB & maxArgB) }

// SetB updates the B of the instruction.
func (i *Instruction) SetB(x int) { i.setArg(posB, sizeB, x) }

// C returns the C of the instruction.
func (i Instruction) C() int { return int(i >> posC & maxArgC) }

// SetC updates the C of the instruction.
func (i *Instruction) SetC(x int) { i.setArg(posC, sizeC, x) }

// NewABC creates an instruction in ABC-form.
func NewABC(op Opcode, a, b, c int) Instruction {
	return Instruction(op)<<posOp |
		Instruction(a)<<posA |
		Instruction(b)<<posB |
		Instruction(c)<<posC
}

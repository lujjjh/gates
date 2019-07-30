package vm

const (
	OpLoadConst byte = iota
	OpLoadNull
	OpLoadGlobal
	OpStoreGlobal
	OpLoadLocal
	OpStoreLocal
	OpUnaryPlus
	OpUnaryMinus
	OpUnaryNot
	OpBinaryAdd
	OpBinarySub
	OpBinaryMul
	OpBinaryDiv
	OpBinaryMod
	OpBinaryEq
	OpBinaryNEq
	OpBinaryLT
	OpBinaryLTE
	OpBinaryGT
	OpBinaryGTE
	OpCall
	OpReturn
)

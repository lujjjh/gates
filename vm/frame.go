package vm

import (
	"github.com/lujjjh/gates"
)

type Frame struct {
	fn *gates.CompiledFunction
	ip int
	bp int
}

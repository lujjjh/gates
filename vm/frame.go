package vm

import (
	"github.com/gates/gates"
)

type Frame struct {
	fn *gates.CompiledFunction
	ip int
	bp int
}

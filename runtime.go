package gates

import (
	"unsafe"
)

type ToNativeOption int

const (
	SkipCircularReference ToNativeOption = 1 << iota
)

func checkToNativeOption(desiredOption ToNativeOption, options int) bool {
	return options&int(desiredOption) == int(desiredOption)
}

func convertToNativeOption2BinaryOptions(options []ToNativeOption) int {
	ops := 0
	for _, op := range options {
		ops |= int(op)
	}
	return ops
}

type Runtime struct {
}

func New() *Runtime {
	r := &Runtime{}
	return r
}

type toNativer interface {
	toNative(seen map[unsafe.Pointer]interface{}, options int) interface{}
}

func toNative(seen map[unsafe.Pointer]interface{}, v Value, options int) (result interface{}) {
	if seen == nil {
		seen = make(map[unsafe.Pointer]interface{})
	}
	if toNativer, haveToNativer := v.(toNativer); haveToNativer {
		return toNativer.toNative(seen, options)
	}
	return v.ToNative()
}

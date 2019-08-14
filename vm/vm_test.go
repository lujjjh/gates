package vm

import (
	"testing"

	"github.com/gates/gates"

	"github.com/stretchr/testify/assert"
)

func TestVM_Load(t *testing.T) {
	v := New(&gates.CompiledFunction{
		Instructions: []byte{
			OpLoadNull,
			OpLoadConst, 0, 1,
			OpLoadGlobal, 0, 0,
		},
	})
	v.constants = []gates.Value{
		gates.Int(0),
		gates.Int(1),
	}
	v.globals = []gates.Value{
		gates.String("foo"),
		gates.String("bar"),
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, gates.Null, v.stack[0])
	assert.Equal(t, gates.Int(1), v.stack[1])
	assert.Equal(t, gates.String("foo"), v.stack[2])
}

func TestVM_Store(t *testing.T) {
	v := New(&gates.CompiledFunction{
		Instructions: []byte{
			OpLoadConst, 0, 0,
			OpStoreGlobal, 0, 1,
		},
	})
	v.constants = []gates.Value{
		gates.Int(42),
	}
	v.globals = make([]gates.Value, 2)
	assert.NoError(t, v.Run())
	assert.Equal(t, gates.Int(42), v.globals[1])
}

func TestVM_Array(t *testing.T) {
	v := New(&gates.CompiledFunction{
		Instructions: []byte{
			OpLoadConst, 0, 0,
			OpLoadConst, 0, 1,
			OpArray, 0, 2,
			OpLoadConst, 0, 2,
			OpMergeArray, 2,
		},
	})
	v.constants = []gates.Value{
		gates.String("foo"),
		gates.String("bar"),
		gates.Array{
			Values: []gates.Value{
				gates.Int(42),
				gates.Null,
			},
		},
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, gates.Array{
		Values: []gates.Value{
			gates.String("foo"),
			gates.String("bar"),
			gates.Int(42),
			gates.Null,
		},
	}, v.stack[0])
}

func TestVM_Map(t *testing.T) {
	v := New(&gates.CompiledFunction{
		Instructions: []byte{
			OpLoadConst, 0, 0,
			OpLoadConst, 0, 1,
			OpLoadConst, 0, 2,
			OpLoadConst, 0, 3,
			OpMap, 0, 2,
			OpLoadConst, 0, 4,
			OpMergeMap, 2,
		},
	})
	v.constants = []gates.Value{
		gates.String("foo"),
		gates.Int(1),
		gates.String("bar"),
		gates.Int(2),
		gates.Map{
			"baz": gates.Int(3),
		},
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, gates.Map{
		"foo": gates.Int(1),
		"bar": gates.Int(2),
		"baz": gates.Int(3),
	}, v.stack[0])
}

func TestVM_Unary(t *testing.T) {
	v := New(&gates.CompiledFunction{
		Instructions: []byte{
			OpLoadConst, 0, 0,
			OpUnaryPlus,
			OpUnaryMinus,
			OpLoadConst, 0, 1,
			OpUnaryNot,
		},
	})
	v.constants = []gates.Value{
		gates.String("-42.0"),
		gates.Null,
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, 2, v.sp)
	assert.Equal(t, gates.Float(42), v.stack[0])
	assert.Equal(t, gates.Bool(true), v.stack[1])
}

func TestVM_Call_Return(t *testing.T) {
	fn := &gates.CompiledFunction{
		Instructions: []byte{
			OpLoadLocal, 0,
			OpLoadLocal, 1,
			OpBinaryAdd,
			OpReturn,
		},
	}

	v := New(&gates.CompiledFunction{
		Instructions: []byte{
			OpLoadGlobal, 0, 0, // fn
			OpLoadConst, 0, 0, // 40
			OpLoadConst, 0, 1, // 2
			OpLoadConst, 0, 1, // 2
			OpCall, 3, // #arguments
		},
	})
	v.constants = []gates.Value{
		gates.Int(40),
		gates.Int(2),
	}
	v.globals = []gates.Value{
		fn,
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, gates.Int(42), v.stack[v.sp-1])
}

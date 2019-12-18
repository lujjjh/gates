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
	v.constants = []interface{}{
		int64(0),
		int64(1),
	}
	v.globals = []interface{}{
		"foo",
		"bar",
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, nil, v.stack[0])
	assert.Equal(t, int64(1), v.stack[1])
	assert.Equal(t, "foo", v.stack[2])
}

func TestVM_Store(t *testing.T) {
	v := New(&gates.CompiledFunction{
		Instructions: []byte{
			OpLoadConst, 0, 0,
			OpStoreGlobal, 0, 1,
		},
	})
	v.constants = []interface{}{
		int64(42),
	}
	v.globals = make([]interface{}, 2)
	assert.NoError(t, v.Run())
	assert.Equal(t, int64(42), v.globals[1])
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
	v.constants = []interface{}{
		"foo",
		"bar",
		[]interface{}{
			int64(42),
			nil,
		},
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, []interface{}{
		"foo",
		"bar",
		int64(42),
		nil,
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
	v.constants = []interface{}{
		"foo",
		int64(1),
		"bar",
		int64(2),
		map[string]interface{}{
			"baz": int64(3),
		},
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, map[string]interface{}{
		"foo": int64(1),
		"bar": int64(2),
		"baz": int64(3),
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
	v.constants = []interface{}{
		"-42.0",
		nil,
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, 2, v.sp)
	assert.Equal(t, float64(42), v.stack[0])
	assert.Equal(t, true, v.stack[1])
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
	v.constants = []interface{}{
		int64(40),
		int64(2),
	}
	v.globals = []interface{}{
		fn,
	}
	assert.NoError(t, v.Run())
	assert.Equal(t, int64(42), v.stack[v.sp-1])
}

package gates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayToNative(t *testing.T) {
	a := []Value{nil, Int(1), nil}
	b := Array(a)
	a[0] = b
	a[2] = b

	x := b.ToNative().([]interface{})
	assert.Equal(t, x, x[0])
	assert.Equal(t, int64(1), x[1])
	assert.Equal(t, x, x[2])

	assert.Equal(t, []interface{}(nil), Array(nil).ToNative())
}

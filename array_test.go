package gates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayToNative(t *testing.T) {
	a := []Value{nil, Int(1), nil}
	b := NewArray(a)
	a[0] = b
	a[2] = b

	x := b.ToNative().([]interface{})
	assert.Equal(t, x, x[0])
	assert.Equal(t, int64(1), x[1])
	assert.Equal(t, x, x[2])

	assert.Equal(t, []interface{}(nil), NewArray(nil).ToNative())
}

func TestArrayToNativeCircular(t *testing.T) {
	assert := assert.New(t)
	c := NewArray([]Value{Int(1)})
	a := []Value{nil, Int(1), nil, c, c}
	b := NewArray(a)
	a[0] = b
	a[2] = b

	x := b.ToNative(SkipCircularReference).([]interface{})
	assert.EqualValues(5, len(x))
	assert.Nil(x[0])
	assert.EqualValues(1, x[1])
	assert.Nil(x[2])
	assert.EqualValues(1, x[3].([]interface{})[0])
	assert.EqualValues(1, x[4].([]interface{})[0])
}

func TestArrayIterator(t *testing.T) {
	a := NewArray([]Value{})
	it := a.Iterator()
	value, ok := it.Next()
	assert.Equal(t, Null, value)
	assert.False(t, ok)

	a = NewArray([]Value{Int(42), String("foo")})
	it = a.Iterator()
	value, ok = it.Next()
	assert.Equal(t, Int(42), value)
	assert.True(t, ok)
	value, ok = it.Next()
	assert.Equal(t, String("foo"), value)
	assert.True(t, ok)
	value, ok = it.Next()
	assert.Equal(t, Null, value)
	assert.False(t, ok)
}

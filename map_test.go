package gates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToNative(t *testing.T) {
	a := map[string]Value{
		"foo": Int(42),
	}
	b := Map(a)
	a["bar"] = b

	x := b.ToNative().(map[string]interface{})
	assert.Equal(t, int64(42), x["foo"])
	assert.Equal(t, x, x["bar"])

	assert.Equal(t, map[string]interface{}(nil), Map(nil).ToNative())
}

func TestMapIterator(t *testing.T) {
	m := Map(map[string]Value{})
	it := m.Iterator()
	value, ok := it.Next()
	assert.Equal(t, Null, value)
	assert.False(t, ok)

	m = Map(map[string]Value{"foo": Int(42), "bar": String("baz"), "deleted": Bool(true)})
	it = m.Iterator()
	value, ok = it.Next()
	assert.Equal(t, Map(map[string]Value{
		"key":   String("bar"),
		"value": String("baz"),
	}), value)
	assert.True(t, ok)
	delete(m, "deleted")
	value, ok = it.Next()
	assert.Equal(t, Map(map[string]Value{
		"key":   String("foo"),
		"value": Int(42),
	}), value)
	assert.True(t, ok)
	value, ok = it.Next()
	assert.Equal(t, Null, value)
	assert.False(t, ok)
}

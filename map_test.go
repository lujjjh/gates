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

package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepresentString(t *testing.T) {
	assert.Equal(t, `"hello world"`, Represent("hello world"))
}

func TestRepresentInt(t *testing.T) {
	assert.Equal(t, "-42", Represent(int64(-42)))
}

func TestRepresentFloat(t *testing.T) {
	assert.Equal(t, "3.14", Represent(float64(3.14)))
}

func TestRepresentNull(t *testing.T) {
	assert.Equal(t, "null", Represent(nil))
}

func TestRepresentBool(t *testing.T) {
	assert.Equal(t, "true", Represent(true))
	assert.Equal(t, "false", Represent(false))
}

func TestRepresentArrayAndMap(t *testing.T) {
	assert.Equal(t, `[3.14,true,null,"hello world",{"ans":42}]`, Represent([]interface{}{
		float64(3.14),
		true,
		nil,
		"hello world",
		map[string]interface{}{"ans": int64(42)},
	}))
}

func TestRepresentCircularStructure(t *testing.T) {
	v := make(map[string]interface{})
	v["foo"] = []interface{}{v}
	v["zoo"] = int64(42)
	assert.Equal(t, `{"foo":[[Circular]],"zoo":42}`, Represent(v))
}

func TestRepresentUnknownType(t *testing.T) {
	assert.Equal(t, `unknown 0 int32`, Represent(int32(0)))
}

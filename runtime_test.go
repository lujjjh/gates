package gates

import (
	"testing"
)

func mustRunStringWithGlobal(s string, global interface{}) Value {
	r := New()
	r.SetGlobal(global)
	value, err := r.RunString(s)
	if err != nil {
		panic(err)
	}
	return value
}

func mustRunString(s string) Value {
	return mustRunStringWithGlobal(s, nil)
}

func assertValue(t *testing.T, expected, actual Value) {
	if !expected.SameAs(actual) {
		t.Errorf("%#v != %#v", actual, expected)
	}
}

func TestRunString(t *testing.T) {
	assertValue(t, intNumber(34), mustRunString("4 + 5 * 6"))
	assertValue(t, floatNumber(0.5), mustRunString("1 / 2"))
	assertValue(t, String("he he"), mustRunString(`"he\x20" + "he"`))
	assertValue(t, floatNumber(1.5), mustRunString(`0 && true || 1.5`))
	assertValue(t, Bool(true), mustRunString(`!(0 && true)`))
	assertValue(t, Bool(true), mustRunString(`1 == "1"`))
	assertValue(t, Bool(true), mustRunString(`"hehe" != ("1" == true)`))
	assertValue(t, Bool(true), mustRunString("1.1 >= 1"))
	assertValue(t, Bool(true), mustRunString(`"abc" > "aba"`))
	assertValue(t, String("nullhehe"), mustRunString(`null + "hehe"`))

	assertValue(t, intNumber(42), mustRunStringWithGlobal("a", map[string]interface{}{
		"a": 42,
	}))
}

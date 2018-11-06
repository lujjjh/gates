package gates

import (
	"testing"
)

func mustRunStribng(s string) Value {
	value, err := New().RunString(s)
	if err != nil {
		panic(err)
	}
	return value
}

func assertValue(t *testing.T, expected, actual Value) {
	if !expected.SameAs(actual) {
		t.Errorf("%#v != %#v", actual, expected)
	}
}

func TestRunString(t *testing.T) {
	assertValue(t, intNumber(34), mustRunStribng("4 + 5 * 6"))
	assertValue(t, floatNumber(0.5), mustRunStribng("1 / 2"))
	assertValue(t, String("he he"), mustRunStribng(`"he\x20" + "he"`))
	assertValue(t, floatNumber(1.5), mustRunStribng(`0 && true || 1.5`))
	assertValue(t, Bool(true), mustRunStribng(`!(0 && true)`))
	assertValue(t, Bool(true), mustRunStribng(`1 == "1"`))
	assertValue(t, Bool(true), mustRunStribng(`"hehe" != ("1" == true)`))
	assertValue(t, Bool(true), mustRunStribng("1.1 >= 1"))
	assertValue(t, Bool(true), mustRunStribng(`"abc" > "aba"`))
}

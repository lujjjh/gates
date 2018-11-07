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
	assertValue(t, Int(34), mustRunString("4 + 5 * 6"))
	assertValue(t, Float(0.5), mustRunString("1 / 2"))
	assertValue(t, String("he he"), mustRunString(`"he\x20" + "he"`))
	assertValue(t, Float(1.5), mustRunString(`0 && true || 1.5`))
	assertValue(t, Bool(true), mustRunString(`!(0 && true)`))
	assertValue(t, Bool(true), mustRunString(`1 == "1"`))
	assertValue(t, Bool(true), mustRunString(`"hehe" != ("1" == true)`))
	assertValue(t, Bool(true), mustRunString("1.1 >= 1"))
	assertValue(t, Bool(true), mustRunString(`"abc" > "aba"`))
	assertValue(t, String("nullhehe"), mustRunString(`null + "hehe"`))

	assertValue(t, Int(42), mustRunStringWithGlobal(`a.b["c"]`, map[string]interface{}{
		"a": getterFunc(func(r *Runtime, v Value) Value {
			return Map(map[string]interface{}{
				"c": 42,
			})
		}),
	}))

	assertValue(t, Int(42), mustRunStringWithGlobal(`a[1*2]`, map[string]interface{}{
		"a": []interface{}{40, 41, 42},
	}))

	assertValue(t, Int(4), mustRunString(`("he" + "he").length`))
	assertValue(t, String("e"), mustRunString(`"hehe"[1]`))
	assertValue(t, Null, mustRunString(`"hehe"[-1]`))
	assertValue(t, Null, mustRunString(`"hehe"[4]`))
}

func BenchmarkRunProgram(b *testing.B) {
	program, err := Compile(`"hello"[0] + 1 && false`)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		r := New()
		for pb.Next() {
			r.Reset()
			r.RunProgram(program)
		}
	})
}

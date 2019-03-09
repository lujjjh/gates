package gates

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"
)

type panicErr struct {
	message string
}

func (p *panicErr) Error() string {
	if p.message != "" {
		return p.message
	}
	return "assertion failed"
}

func mustRunStringWithGlobal(s string, global map[string]Value) Value {
	r := New()
	for k, v := range global {
		r.Global().Set(k, v)
	}
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
	assertValue(t, String("hehe"), mustRunString(`null + "hehe"`))

	assertValue(t, Int(42), mustRunStringWithGlobal(`a.b["c"]`, map[string]Value{
		"a": ref(getterFunc(func(r *Runtime, v Value) Value {
			return Map(map[string]Value{
				"c": Int(42),
			})
		})),
	}))

	assertValue(t, Int(42), mustRunStringWithGlobal(`a[1*2]`, map[string]Value{
		"a": Array([]Value{Int(40), Int(41), Int(42)}),
	}))

	assertValue(t, Int(4), mustRunString(`("he" + "he").length`))
	assertValue(t, String("e"), mustRunString(`"hehe"[1]`))
	assertValue(t, Null, mustRunString(`"hehe"[-1]`))
	assertValue(t, Null, mustRunString(`"hehe"[4]`))

	assertValue(t, Float(3), mustRunStringWithGlobal(`add(1, 2)`, map[string]Value{
		"add": FunctionFunc(func(fc FunctionCall) Value {
			var result float64
			for _, arg := range fc.Args() {
				result += arg.ToFloat()
			}
			return Float(result)
		}),
	}))

	assertValue(t, Bool(true), mustRunString(`[] == []`))
	assertValue(t, Bool(true), mustRunString(`[1] == [1]`))
	assertValue(t, Bool(false), mustRunString(`[1] == [1, 2]`))
	assertValue(t, Bool(true), mustRunString(`{} == {}`))
	assertValue(t, Bool(true), mustRunString(`{ a: 1 } == { a: 1 }`))
	assertValue(t, Bool(false), mustRunString(`{ a: 1 } == { a: 1, b: 2 }`))

	assertValue(t, Int(42), mustRunString(`[0, 42][1]`))
	assertValue(t, String("bar"), mustRunString(`({foo: "bar"}).foo`))
	assertValue(t, String("bar"), mustRunString(`({"foo": "bar"}).foo`))
	assertValue(t, String("bar"), mustRunString(`({["foo"]: "bar", bar: "baz"}).foo`))

	assertValue(t, Int(42), mustRunString(`function (a, b) {
		return function (c) { return a + c; }(b + 1);
	}(1, 40)`))

	assertValue(t, Null, mustRunString(`function () {}()`))
	assertValue(t, Null, mustRunString(`function () { return; }()`))
	assertValue(t, Null, mustRunString(`function (a) { return a; }()`))

	assertValue(t, Int(89), mustRunString(`
		(function (x) {
			return function (f) {
				return function (n) {
					return f(x(x)(f))(n);
				};
			};
		})(function (x) {
			return function (f) {
				return function (n) {
					return f(x(x)(f))(n);
				};
			};
		})(function (f) {
			return function (n) {
				return (n == 0 || n == 1) && 1 || f(n - 2) + f(n - 1);
			};
		})(10)
	`))
}

func TestRunProgramStackOverflow(t *testing.T) {
	src := `
		(function (x) {
			return function (f) {
				return function () {
					return f(x(x)(f))();
				};
			};
		})(function (x) {
			return function (f) {
				return function () {
					return f(x(x)(f))();
				};
			};
		})(function (f) {
			return function () {
				return f();
			};
		})()
	`

	r := New()
	_, err := r.RunString(src)
	if err != ErrStackOverflow {
		t.Errorf("stack overflow expected")
	}
}

func TestRunFiles(t *testing.T) {
	fileInfo, err := ioutil.ReadDir("testdata/")
	if err != nil {
		t.Error(err)
		return
	}
	for _, f := range fileInfo {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !strings.HasSuffix(name, ".gates") {
			continue
		}

		s, err := ioutil.ReadFile("testdata/" + name)
		if err != nil {
			t.Error(err)
			continue
		}
		r := New()
		g := r.Global()
		_assert := FunctionFunc(func(fc FunctionCall) Value {
			args := fc.Args()
			argc := len(args)
			if argc < 1 {
				panic(&panicErr{message: "assert takes at least 1 arguments"})
			}
			if !args[0].ToBool() {
				if argc < 2 {
					panic(&panicErr{})
				}
				panic(&panicErr{message: args[1].ToString()})
			}
			return Bool(true)
		})
		g.Set("assert", _assert)
		g.Set("assert_eq", FunctionFunc(func(fc FunctionCall) Value {
			args := fc.Args()
			argc := len(args)
			if argc < 2 {
				panic(&panicErr{message: "assert_equal takes 2 arguments"})
			}
			message := String(args[0].ToString() + " expected, got " + args[1].ToString())
			return r.Call(_assert, Bool(args[0].Equals(args[1])), message)
		}))
		_, err = r.RunString(string(s))
		if err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkRunProgram(b *testing.B) {
	program, err := Compile(`"hello"[0] + 1 && false`)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		r := New()
		for pb.Next() {
			r.RunProgram(ctx, program)
		}
	})
}

package gates_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gates/gates"
)

var r = gates.New()

func assertEqual(a, b string, messages ...interface{}) {
	if len(messages) == 0 {
		messages = []interface{}{"it should be equaled"}
	}
	if a != b {
		panic(fmt.Sprint(messages...))
	}
}

func assertNotEqual(a, b string, messages ...interface{}) {
	if len(messages) == 0 {
		messages = []interface{}{"it should not be equaled"}
	}
	if a == b {
		panic(fmt.Sprint(messages...))
	}
}

func assertTrue(b bool, messages ...interface{}) {
	if len(messages) == 0 {
		messages = []interface{}{"it should be true"}
	}
	if !b {
		panic(fmt.Sprint(messages...))
	}
}

func assertFalse(b bool, messages ...interface{}) {
	if len(messages) == 0 {
		messages = []interface{}{"it should be false"}
	}
	if b {
		panic(fmt.Sprint(messages...))
	}
}

func assertNil(v interface{}, messages ...interface{}) {
	if len(messages) == 0 {
		messages = []interface{}{"it should be nil"}
	}
	if v != nil {
		panic(fmt.Sprint(messages...))
	}
}

func assertNotNil(v interface{}, messages ...interface{}) {
	if len(messages) == 0 {
		messages = []interface{}{"it should not be nil"}
	}
	if v == nil {
		panic(fmt.Sprint(messages...))
	}
}

func TestHasPrefix(t *testing.T) {
	v, err := r.RunString(`strings.has_prefix("foobar", "foo")`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.has_prefix("foobar")`)
	assertNil(err)
	assertFalse(v.ToBool())
}

func TestHasSuffix(t *testing.T) {
	v, err := r.RunString(`strings.has_suffix("foobar", "bar")`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.has_suffix("foobar", "foo")`)
	assertNil(err)
	assertFalse(v.ToBool())

	v, err = r.RunString(`strings.has_suffix("foobar")`)
	assertNil(err)
	assertFalse(v.ToBool())

	v, err = r.RunString(`strings.has_suffix()`)
	assertNil(err)
	assertFalse(v.ToBool())
}

func TestToLower(t *testing.T) {
	v, err := r.RunString(`strings.to_lower("Foo_Bar")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo_bar")

	v, err = r.RunString(`strings.to_lower("foo_bar")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo_bar")

	v, err = r.RunString(`strings.to_lower("1234")`)
	assertNil(err)
	assertEqual(v.ToString(), "1234")
}

func TestToUpper(t *testing.T) {
	v, err := r.RunString(`strings.to_upper("Foo_Bar")`)
	assertNil(err)
	assertEqual(v.ToString(), "FOO_BAR")

	v, err = r.RunString(`strings.to_upper("FOO_BAR")`)
	assertNil(err)
	assertEqual(v.ToString(), "FOO_BAR")

	v, err = r.RunString(`strings.to_upper("1234")`)
	assertNil(err)
	assertEqual(v.ToString(), "1234")
}

func TestTrim(t *testing.T) {
	v, err := r.RunString(`strings.trim("foo__foo", "foo")`)
	assertNil(err)
	assertEqual(v.ToString(), "__")

	v, err = r.RunString(`strings.trim("foo__foo", "__")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo__foo")

	v, err = r.RunString(`strings.trim("foo__", "__")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo")

	v, err = r.RunString(`strings.trim("foo__", "foo")`)
	assertNil(err)
	assertEqual(v.ToString(), "__")

	v, err = r.RunString(`strings.trim("  foo__bar  ")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo__bar")

	v, err = r.RunString(`strings.trim("foo  bar")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo  bar")

	v, err = r.RunString(`strings.trim("  foo  bar   ")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo  bar")
}

func TestTrimLeft(t *testing.T) {
	v, err := r.RunString(`strings.trim_left("foo__foo", "foo")`)
	assertNil(err)
	assertEqual(v.ToString(), "__foo")

	v, err = r.RunString(`strings.trim_left("foo__bar", "foo")`)
	assertNil(err)
	assertEqual(v.ToString(), "__bar")

	v, err = r.RunString(`strings.trim_left("foo__", "__")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo__")

	v, err = r.RunString(`strings.trim_left("foo__", "bar")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo__")
}

func TestTrimRight(t *testing.T) {
	v, err := r.RunString(`strings.trim_right("foo__foo", "foo")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo__")

	v, err = r.RunString(`strings.trim_right("foo__bar", "bar")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo__")

	v, err = r.RunString(`strings.trim_right("__bar", "__")`)
	assertNil(err)
	assertEqual(v.ToString(), "__bar")

	v, err = r.RunString(`strings.trim_right("foo__", "foo")`)
	assertNil(err)
	assertEqual(v.ToString(), "foo__")
}

func TestSplit(t *testing.T) {
	v, err := r.RunString(`strings.split("foo,bar", ",")`)
	assertNil(err)
	_ = v
	// TODO(cloverstd): test array elements
}

func TestJoin(t *testing.T) {
	v, err := r.RunString(`strings.join([1, "2", true, 1.1], "|")`)
	assertNil(err)
	assertEqual(v.ToString(), "1|2|true|1.1")
}

func TestMatch(t *testing.T) {
	v, err := r.RunString(`strings.match("(?P<first_name>\\w+) (?P<last_name>\\w+)", "Malcolm Reynolds").group(1) == "Malcolm"`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.match("(?P<first_name>\\w+) (?P<last_name>\\w+)", "Malcolm Reynolds").group("first_name") == "Malcolm"`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.match("(\\w+) (?P<last_name>\\w+)", "Malcolm Reynolds").group("first_name") == null`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.match("(\\w+) (?P<last_name>\\w+)", "Malcolm Reynolds").group(1) == "Malcolm"`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.match("(\\w+) (?P<last_name>\\w+)", "Malcolm Reynolds").group("last_name") == "Reynolds"`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.match("(\\w+) (?P<last_name>\\w+)", "Malcolm Reynolds").group(-1) == null`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.match("(\\w+) (?P<last_name>\\w+)", "Malcolm Reynolds").group(10000) == null`)
	assertNil(err)
	assertTrue(v.ToBool())

	v, err = r.RunString(`strings.match("(?i)test", "Test") != null`)
	assertNil(err)
	assertTrue(v.ToBool())
}

func TestFindALl(t *testing.T) {
	v, err := r.RunString(`
	(function() {
		let results = strings.find_all("(?i)(foo)", "foo\nfOo\nFOO");
		return results[0] == "foo" && results[1] == "fOo" && results[2] == "FOO";
	})()
	`)
	assertNil(err)
	assertTrue(v.ToBool())
}

func TestContains(t *testing.T) {
	v, err := r.RunString(`strings.contains("foobarfoo", "foo")`)
	assertNil(err)
	assertTrue(v == gates.True)

	v, err = r.RunString(`strings.contains("foobarfoo", "abc")`)
	assertNil(err)
	assertTrue(v == gates.False)

	v, err = r.RunString(`strings.contains_any("foobarfoobar", "barfoo")`)
	assertNil(err)
	assertTrue(v == gates.True)
}

func TestIndex(t *testing.T) {
	v, err := r.RunString(`strings.index("foobarfoo", "foo")`)
	assertNil(err)
	assertTrue(v.ToInt() == 0)

	v, err = r.RunString(`strings.index("foobarfoo", "bar")`)
	assertNil(err)
	assertTrue(v.ToInt() == 3)

	v, err = r.RunString(`strings.index("foobarfoo", "rba")`)
	assertNil(err)
	assertTrue(v.ToInt() == -1)

	v, err = r.RunString(`strings.index_any("barfoo", "aroo")`)
	assertNil(err)
	assertTrue(v.ToInt() == 1)

	v, err = r.RunString(`strings.index_any("barfoo", "ttoo")`)
	assertNil(err)
	assertTrue(v.ToInt() == 4)

	v, err = r.RunString(`strings.index_any("barfoo", "ttt")`)
	assertNil(err)
	assertTrue(v.ToInt() == -1)

	v, err = r.RunString(`strings.last_index("foobarfoo", "foo")`)
	assertNil(err)
	assertTrue(v.ToInt() == 6)

	v, err = r.RunString(`strings.last_index_any("foobarfoo", "foobb")`)
	assertNil(err)
	assertTrue(v.ToInt() == 8)

	v, err = r.RunString(`strings.last_index_any("foobarfoo", "tttt")`)
	assertNil(err)
	assertTrue(v.ToInt() == -1)
}

func TestRepeat(t *testing.T) {
	v, err := r.RunString(`strings.repeat("foo", 5)`)
	assertNil(err)
	assertEqual(strings.Repeat("foo", 5), v.ToString())

	v, err = r.RunString(`strings.repeat("foo", "a")`)
	assertNil(err)
	assertEqual(strings.Repeat("foo", 0), v.ToString())

	v, err = r.RunString(`strings.repeat("foo", -1) == null`)
	assertNil(err)
	assertTrue(v.ToBool())
}

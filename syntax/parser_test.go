package syntax

import (
	"testing"
)

var fset = NewFileSet()

func TestParseExpr(t *testing.T) {
	src := "a + b"
	x, err := ParseExpr(src)
	if err != nil {
		t.Errorf("ParseExpr(%q): %v", src, err)
	}
	// sanity check
	if _, ok := x.(*BinaryExpr); !ok {
		t.Errorf("ParseExpr(%q): got %T, want *BinaryExpr", src, x)
	}

	// an invalid expression
	src = "a + *"
	if _, err := ParseExpr(src); err == nil {
		t.Errorf("ParseExpr(%q): got no error", src)
	}

	// newline is permitted
	src = "(a + b) * -2.5\n"
	if _, err := ParseExpr(src); err != nil {
		t.Errorf("ParseExpr(%q): got error %s", src, err)
	}

	// various other stuff following a valid expression
	const validExpr = "a + b"
	const anything = "dh3*#D)#_"
	for _, c := range "!)]};," {
		src := validExpr + string(c) + anything
		if _, err := ParseExpr(src); err == nil {
			t.Errorf("ParseExpr(%q): got no error", src)
		}
	}

	// function call
	src = `a.b["c"](b(1+1))(2,"hello"[42])`
	if _, err := ParseExpr(src); err != nil {
		t.Errorf("ParseExpr(%q): got error %s", src, err)
	}
}

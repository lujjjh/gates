package syntax

import (
	"path/filepath"
	"testing"
)

type elt struct {
	tok Token
	lit string
}

var scanTokens = [...]elt{
	// Identifiers and basic type literals
	{IDENT, "foobar"},
	{IDENT, "$_42"},
	{IDENT, "中国"},
	{NUMBER, "0"},
	{NUMBER, "1"},
	{NUMBER, "123456789012345678890"},
	{NUMBER, "01234567"},
	{NUMBER, "0xcafebabe"},
	{NUMBER, "0."},
	{NUMBER, ".0"},
	{NUMBER, "3.14159265"},
	{NUMBER, "1e0"},
	{NUMBER, "1e+100"},
	{NUMBER, "1e-100"},
	{NUMBER, "2.71828e-1000"},
	{STRING, `"foobar"`},
	{STRING, `"foobar\n\0123\x0020"`},

	// Operators and delimiters
	{ADD, "+"},
	{SUB, "-"},
	{MUL, "*"},
	{QUO, "/"},
	{REM, "%"},

	{AND, "&"},
	{OR, "|"},
	{XOR, "^"},
	{SHL, "<<"},
	{SHR, ">>"},

	{LAND, "&&"},
	{LOR, "||"},

	{EQL, "=="},
	{LSS, "<"},
	{GTR, ">"},
	{NOT, "!"},

	{NEQ, "!="},
	{LEQ, "<="},
	{GEQ, ">="},

	{LPAREN, "("},
	{LBRACK, "["},
	{LBRACE, "{"},
	{COMMA, ","},
	{PERIOD, "."},
	{ELLIPSIS, "..."},

	{RPAREN, ")"},
	{RBRACK, "]"},
	{RBRACE, "}"},
	{SEMICOLON, ";"},
	{COLON, ":"},
}

const whitespace = "  \t  \n\n\n" // to separate tokens

var source = func() []byte {
	var src []byte
	for _, t := range scanTokens {
		src = append(src, t.lit...)
		src = append(src, whitespace...)
	}
	return src
}()

func newlineCount(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n++
		}
	}
	return n
}

func checkPos(t *testing.T, lit string, p Pos, expected Position) {
	pos := fset.Position(p)
	// Check cleaned filenames so that we don't have to worry about
	// different os.PathSeparator values.
	if pos.Filename != expected.Filename && filepath.Clean(pos.Filename) != filepath.Clean(expected.Filename) {
		t.Errorf("bad filename for %q: got %s, expected %s", lit, pos.Filename, expected.Filename)
	}
	if pos.Offset != expected.Offset {
		t.Errorf("bad position for %q: got %d, expected %d", lit, pos.Offset, expected.Offset)
	}
	if pos.Line != expected.Line {
		t.Errorf("bad line for %q: got %d, expected %d", lit, pos.Line, expected.Line)
	}
	if pos.Column != expected.Column {
		t.Errorf("bad column for %q: got %d, expected %d", lit, pos.Column, expected.Column)
	}
}

// Verify that calling Scan() provides the correct results.
func TestScan(t *testing.T) {
	whitespaceLinecount := newlineCount(whitespace)

	// error handler
	eh := func(_ Position, msg string) {
		t.Errorf("error handler called (msg = %s)", msg)
	}

	// verify scan
	var s Scanner
	s.Init(fset.AddFile("", fset.Base(), len(source)), source, eh)

	// set up expected position
	epos := Position{
		Filename: "",
		Offset:   0,
		Line:     1,
		Column:   1,
	}

	index := 0
	for {
		pos, tok, lit := s.Scan()

		// check position
		if tok == EOF {
			// correction for EOF
			epos.Line = newlineCount(string(source))
			epos.Column = 2
		}
		checkPos(t, lit, pos, epos)

		// check token
		e := elt{EOF, ""}
		if index < len(scanTokens) {
			e = scanTokens[index]
			index++
		}
		if tok != e.tok {
			t.Errorf("bad token for %q: got %s, expected %s", lit, tok, e.tok)
		}

		// check literal
		elit := ""
		switch e.tok {
		case IDENT:
			elit = e.lit
		case SEMICOLON:
			elit = ";"
		default:
			if e.tok.IsLiteral() {
				elit = e.lit
			}
		}
		if lit != elit {
			t.Errorf("bad literal for %q: got %q, expected %q", lit, lit, elit)
		}

		if tok == EOF {
			break
		}

		// update position
		epos.Offset += len(e.lit) + len(whitespace)
		epos.Line += newlineCount(e.lit) + whitespaceLinecount
	}

	if s.ErrorCount != 0 {
		t.Errorf("found %d errors", s.ErrorCount)
	}
}

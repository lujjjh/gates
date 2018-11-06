package syntax

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// An ErrorHandler may be provided to Scanner.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a
// position and an error message. The position points to the beginning of
// the offending
//
type ErrorHandler func(pos Position, msg string)

// Scanner scans tokens.
type Scanner struct {
	// immutable state
	file *File
	src  []byte
	err  ErrorHandler

	// scanning state
	ch         rune
	offset     int
	rdOffset   int
	lineOffset int

	// public state - ok to modify
	ErrorCount int // number of errors encountered
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError {
				s.error(s.offset, "illegal UTF-8 encoding")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		s.ch = -1 // eof
	}
}

// Init prepares the scanner s to tokenize the text src.
func (s *Scanner) Init(file *File, src []byte, err ErrorHandler) {
	// Explicitly initialize all fields since a scanner may be reused.
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}
	s.file = file
	s.src = src
	s.err = err

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.ErrorCount = 0

	s.next()
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offs)), msg)
	}
	s.ErrorCount++
}

func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.ch) {
		s.next()
	}
}

// Helper functions for scanning multi-byte tokens such as >> += >>= .
// Different routines recognize different length tok_i based on matches
// of ch_i. If a token ends in '=', the result is tok1 or tok3
// respectively. Otherwise, the result is tok0 if there was no other
// matching character, or tok2 if the matching character was ch2.

func (s *Scanner) switch2(tok0, tok1 Token) Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	return tok0
}

func (s *Scanner) switch3(tok0, tok1 Token, ch2 rune, tok2 Token) Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	if s.ch == ch2 {
		s.next()
		return tok2
	}
	return tok0
}

func (s *Scanner) switch4(tok0, tok1 Token, ch2 rune, tok2, tok3 Token) Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	if s.ch == ch2 {
		s.next()
		if s.ch == '=' {
			s.next()
			return tok3
		}
		return tok2
	}
	return tok0
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || s.ch == '$' || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

func (s *Scanner) scanMantissa(base int) {
	for digitVal(s.ch) < base {
		s.next()
	}
}

func (s *Scanner) scanNumber(seenDecimalPoint bool) string {
	// digitVal(s.ch) < 10
	offs := s.offset

	if seenDecimalPoint {
		offs--
		s.scanMantissa(10)
		goto exponent
	}

	if s.ch == '0' {
		// int or float
		offs := s.offset
		s.next()
		if s.ch == 'x' || s.ch == 'X' {
			// hexadecimal int
			s.next()
			s.scanMantissa(16)
			if s.offset-offs <= 2 {
				// only scanned "0x" or "0X"
				s.error(offs, "illegal hexadecimal number")
			}
		} else {
			// octal int or float
			seenDecimalDigit := false
			s.scanMantissa(8)
			if s.ch == '8' || s.ch == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				s.scanMantissa(10)
			}
			if s.ch == '.' || s.ch == 'e' || s.ch == 'E' {
				goto fraction
			}
			// octal int
			if seenDecimalDigit {
				s.error(offs, "illegal octal number")
			}
		}
		goto exit
	}

	// decimal int or float
	s.scanMantissa(10)

fraction:
	if s.ch == '.' {
		s.next()
		s.scanMantissa(10)
	}

exponent:
	if s.ch == 'e' || s.ch == 'E' {
		s.next()
		if s.ch == '-' || s.ch == '+' {
			s.next()
		}
		if digitVal(s.ch) < 10 {
			s.scanMantissa(10)
		} else {
			s.error(offs, "illegal floating-point exponent")
		}
	}

exit:
	return string(s.src[offs:s.offset])
}

// scanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax error, it stops at the offending
// character (without consuming it) and returns false. Otherwise
// it returns true.
func (s *Scanner) scanEscape(quote rune) bool {
	offs := s.offset

	var n int
	var base, max uint32
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.ch < 0 {
			msg = "escape sequence not terminated"
		}
		s.error(offs, msg)
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", s.ch)
			if s.ch < 0 {
				msg = "escape sequence not terminated"
			}
			s.error(s.offset, msg)
			return false
		}
		x = x*base + d
		s.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(offs, "escape sequence is invalid Unicode code point")
		return false
	}

	return true
}

func (s *Scanner) scanString(quote rune) string {
	// '"' opening already consumed
	offs := s.offset - 1

	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			s.error(offs, "string literal not terminated")
			break
		}
		s.next()
		if ch == quote {
			break
		}
		if ch == '\\' {
			s.scanEscape(quote)
		}
	}

	return string(s.src[offs:s.offset])
}

// Scan scans tokens.
func (s *Scanner) Scan() (pos Pos, tok Token, lit string) {
	s.skipWhitespace()

	// current token start
	pos = s.file.Pos(s.offset)

	switch ch := s.ch; {
	case isLetter(ch) || ch == '$':
		tok = IDENT
		lit = s.scanIdentifier()
		switch lit {
		case "true", "false":
			tok = BOOL
		case "null":
			tok = NULL
		}
	case '0' <= ch && ch <= '9':
		tok = NUMBER
		lit = s.scanNumber(false)
	default:
		s.next() // always make progress
		switch ch {
		case -1:
			tok = EOF
		case '"':
			tok = STRING
			lit = s.scanString(ch)
		case ':':
			tok = COLON
		case '.':
			if '0' <= s.ch && s.ch <= '9' {
				tok = NUMBER
				lit = s.scanNumber(true)
			} else {
				tok = PERIOD
			}
		case ',':
			tok = COMMA
		case ';':
			tok = SEMICOLON
			lit = ";"
		case '(':
			tok = LPAREN
		case ')':
			tok = RPAREN
		case '[':
			tok = LBRACK
		case ']':
			tok = RBRACK
		case '{':
			tok = LBRACE
		case '}':
			tok = RBRACE
		case '+':
			tok = ADD
		case '-':
			tok = SUB
		case '*':
			tok = MUL
		case '/':
			tok = QUO
		case '%':
			tok = REM
		case '^':
			tok = XOR
		case '<':
			tok = s.switch3(LSS, LEQ, '<', SHL)
		case '>':
			tok = s.switch3(GTR, GEQ, '>', SHR)
		case '=':
			if s.ch == '=' {
				s.next()
				tok = EQL
			}
		case '!':
			tok = s.switch2(NOT, NEQ)
		case '&':
			tok = s.switch3(AND, ILLEGAL, '&', LAND)
		case '|':
			tok = s.switch3(OR, ILLEGAL, '|', LOR)
		default:
			tok = ILLEGAL
			lit = string(ch)
		}
	}

	return
}

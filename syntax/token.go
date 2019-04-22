package syntax

import "strconv"

// Token indicates a token.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF

	literalBeg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT    // main
	NUMBER   // 123.45
	STRING   // "abc"
	BOOL     // true
	NULL     // null
	LET      // let
	FUNCTION // function
	IF       // if
	ELSE     // else
	FOR      // for
	RETURN   // return
	literalEnd

	operatorBeg
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	PIPE // |

	XOR // ^
	SHL // <<
	SHR // >>

	LAND // &&
	LOR  // ||

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ // !=
	LEQ // <=
	GEQ // >=

	LPAREN   // (
	LBRACK   // [
	LBRACE   // {
	COMMA    // ,
	PERIOD   // .
	ELLIPSIS // ...

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	operatorEnd

	othersBeg
	ARROW // =>
	othersEnd
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF: "EOF",

	IDENT:    "IDENT",
	NUMBER:   "NUMBER",
	STRING:   "STRING",
	BOOL:     "BOOL",
	NULL:     "NULL",
	LET:      "LET",
	FUNCTION: "FUNCTION",
	IF:       "IF",
	ELSE:     "ELSE",
	FOR:      "FOR",
	RETURN:   "RETURN",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	PIPE: "|",

	XOR: "^",
	SHL: "<<",
	SHR: ">>",

	LAND: "&&",
	LOR:  "||",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ: "!=",
	LEQ: "<=",
	GEQ: ">=",

	LPAREN:   "(",
	LBRACK:   "[",
	LBRACE:   "{",
	COMMA:    ",",
	PERIOD:   ".",
	ELLIPSIS: "...",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	ARROW: "=>",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 7
	HighestPrec = 8
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (op Token) Precedence() int {
	switch op {
	case PIPE:
		return 1
	case LOR:
		return 2
	case LAND:
		return 3
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 4
	case ADD, SUB, XOR:
		return 5
	case MUL, QUO, REM, SHL, SHR:
		return 6
	}
	return LowestPrec
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literalBeg < tok && tok < literalEnd }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operatorBeg < tok && tok < operatorEnd }

package gridregexp

import (
	"bytes"

	"github.com/pkg/errors"
)

type Op uint8

const (
	OpLiteral Op = 1 + iota // match Literal
	OpSelect                // match one of Option
	OpEmpty                 // match empty string (inside Option only)
)

type GridRegexp struct {
	Op Op // operator

	// op == OpLiteral, matching literal
	Literal string

	// op == OpSelect, matching one
	Option []*GridRegexp

	// following part
	Next *GridRegexp
}

func Parse(exp string) (*GridRegexp, error) {
	n := len(exp)
	if n < 3 || exp[0] != '^' || exp[n-1] != '$' {
		return nil, errors.New("short/invalid regexp")
	}

	parser := &parser{
		src: exp[1 : n-1],
	}
	return parser.parse(), nil
}

func (x *GridRegexp) String() string {
	var buf bytes.Buffer
	writeExpr(&buf, x)
	return "^" + buf.String() + "$"
}

func writeExpr(buf *bytes.Buffer, x *GridRegexp) {
	if x == nil {
		buf.WriteString("nil")
	}
	for ; x != nil; x = x.Next {
		switch x.Op {

		case OpLiteral:
			buf.WriteString(x.Literal)

		case OpEmpty:
			// pass

		case OpSelect:
			pre := byte('(')
			for _, sub := range x.Option {
				buf.WriteByte(pre)
				writeExpr(buf, sub)
				pre = '|'
			}
			buf.WriteByte(')')

		default:
			buf.WriteByte('?')

		}
	}
}

type parser struct {
	src string

	index int // into src
}

func (x *parser) peekChar() byte {
	if x.index < len(x.src) {
		return x.src[x.index]
	}
	return 0
}

func lit(ch byte) bool {
	switch ch {
	case 'N', 'S', 'E', 'W':
		return true
	}
	return false
}

func litOrOpen(ch byte) bool {
	return ch == '(' || lit(ch)
}

func (x *parser) nextToken() string {
	switch ch := x.peekChar(); ch {
	case '(', '|', ')':
		x.index++
		return string(rune(ch))

	case 'N', 'S', 'E', 'W':
		// literal
		start := x.index
		for x.index++; lit(x.peekChar()); x.index++ {
		}
		return x.src[start:x.index]

	default:
		x.unexpected(x.index)
	}

	panic("unreachable")
}

func (x *parser) unread() { x.index-- }

func (x *parser) unexpected(i int) {
	panic(errors.Errorf("unexpected rune %q at %d", x.peekChar(), x.index))
}

func (x *parser) parse() *GridRegexp {
	var first, last *GridRegexp
	for litOrOpen(x.peekChar()) {

		tok := x.nextToken()

		var cur *GridRegexp
		if tok == "(" {
			cur = x.parseSelect()
		} else {
			cur = &GridRegexp{
				Op:      OpLiteral,
				Literal: tok,
			}
		}

		if last == nil {
			first = cur
		} else {
			last.Next = cur
		}
		last = cur
	}
	return first
}

func (x *parser) parseSelect() *GridRegexp {
	exp := &GridRegexp{
		Op: OpSelect,
	}

	addempty := false
	wassep := true // last bit was separator

Loop:
	for {
		switch ch := x.peekChar(); ch {

		case '|', ')':
			x.index++
			if wassep {
				addempty = true
			}
			if ch == ')' {
				break Loop
			}
			wassep = true

		default:
			sub := x.parse()
			if sub == nil {
				panic(errors.New("parseSelect logic"))
			}
			exp.Option = append(exp.Option, sub)
			wassep = false
		}
	}

	if addempty {
		exp.Option = append(exp.Option, &GridRegexp{Op: OpEmpty})
	}

	return exp
}

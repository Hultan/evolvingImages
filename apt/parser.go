package apt

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	EOF rune = -1
)

type tokenType int

const (
	openParen tokenType = iota
	closeParen
	operator
	constant
)

type token struct {
	typ   tokenType
	value string
}

type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens chan token
}

type stateFunc func(*lexer) stateFunc

func stringToNode(s string) Node {
	switch s {
	case "Picture":
		return NewPicture()
	case "+":
		return NewPlus()
	case "-":
		return NewMinus()
	case "*":
		return NewMult()
	case "/":
		return NewDiv()
	case "Atan2":
		return NewAtan2()
	case "Atan":
		return NewAtan()
	case "Cos":
		return NewCos()
	case "Sin":
		return NewSin()
	case "SimplexNoise":
		return NewNoise()
	case "Square":
		return NewSquare()
	case "Log2":
		return NewLog2()
	case "Negate":
		return NewNegate()
	case "Ceil":
		return NewCeil()
	case "Floor":
		return NewFloor()
	case "Abs":
		return NewAbs()
	case "Clip":
		return NewClip()
	case "Wrap":
		return NewWrap()
	case "Lerp":
		return NewLerp()
	case "FBM":
		return NewFBM()
	case "Turbulence":
		return NewTurbulence()
	case "Swirl":
		return NewSwirl()
	case "x":
		return NewX()
	case "y":
		return NewY()
	default:
		panic(fmt.Sprintf("Unknown token type: %s", s))
	}
}

func parse(tokens chan token, parent Node) Node {
	for {
		tok, ok := <-tokens
		if !ok {
			panic("no more tokens")
		}
		switch tok.typ {
		case operator:
			n := stringToNode(tok.value)
			n.SetParent(parent)
			for i := range n.GetChildren() {
				n.GetChildren()[i] = parse(tokens, n)
			}
			return n
		case constant:
			n := NewConstant()
			n.SetParent(parent)
			num, err := strconv.ParseFloat(tok.value, 64)
			if err != nil {
				panic(fmt.Sprintf("Error while parsing constant : %s", err))
			}
			n.Value = num
			return n
		case openParen:
			continue
		case closeParen:
			continue
		}
	}
	return nil
}

func BeginLexing(input string) Node {
	l := &lexer{
		input:  input,
		tokens: make(chan token, 100), // Buffered so that the lexer can work independently of the parser
	}

	go l.run()
	return parse(l.tokens, nil)
}

func (l *lexer) run() {
	for state := determineToken; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func determineToken(l *lexer) stateFunc {
	for {
		switch r := l.next(); {
		case isWhiteSpace(r):
			l.ignore()
		case r == '(':
			l.emit(openParen)
		case r == ')':
			l.emit(closeParen)
		case isStartOfNumber(r):
			return lexNumber
		case r == EOF:
			return nil
		default:
			return lexOp
		}
	}
}

func lexOp(l *lexer) stateFunc {
	l.acceptRun("+-/*abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	l.emit(operator)
	return determineToken
}

func lexNumber(l *lexer) stateFunc {
	l.accept("-.")
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.input[l.start:l.pos] == "-" {
		l.emit(operator)
	} else {
		l.emit(constant)
	}

	return determineToken
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{
		t,
		l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func isStartOfNumber(r rune) bool {
	return (r >= '0' && r <= '9') || r == '-' || r == '.'
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) peek() rune {
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

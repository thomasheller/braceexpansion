package braceexpansion

import (
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType
	val string
}

type itemType int

const (
	itemError itemType = iota

	itemOpen
	itemClose
	itemSeparator
	itemText
	itemEOF
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	start int
	pos   int
	width int
	items chan item
	opts  ParseOpts
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

func (l *lexer) drain() {
	for range l.items {
	}
}

func lex(input string, opts ParseOpts) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
		opts:  opts,
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// state functions

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], l.opts.OpenBrace) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexOpen
		}
		if strings.HasPrefix(l.input[l.pos:], l.opts.CloseBrace) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexClose
		}
		if strings.HasPrefix(l.input[l.pos:], l.opts.Separator) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexSeparator
		}
		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexOpen(l *lexer) stateFn {
	l.pos += len(l.opts.OpenBrace)
	l.emit(itemOpen)
	return lexText
}

func lexClose(l *lexer) stateFn {
	l.pos += len(l.opts.CloseBrace)
	l.emit(itemClose)
	return lexText
}

func lexSeparator(l *lexer) stateFn {
	l.pos += len(l.opts.Separator)
	l.emit(itemSeparator)
	return lexText
}

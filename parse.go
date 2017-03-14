package braceexpansion

import (
	"fmt"
	"runtime"
)

type Tree struct {
	Root      *ListNode
	lex       *lexer
	token     item
	peekCount int
	opts      ParseOpts
}

type ParseOpts struct {
	OpenBrace             string
	CloseBrace            string
	Separator             string
	TreatRootAsList       bool
	TreatSingleAsOptional bool
}

func (t *Tree) recover(err *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.lex.drain()
			t.stopParse()
		}
		*err = e.(error)
	}
	return
}

func (t *Tree) startParse(l *lexer) {
	t.Root = nil
	t.lex = l
}

func (t *Tree) stopParse() {
	t.lex = nil
}

func (t *Tree) Parse(input string) (tree *Tree, err error) {
	defer t.recover(&err)
	opts := ParseOpts{OpenBrace: "{", CloseBrace: "}", Separator: ","}
	t.opts = opts
	t.startParse(lex(input, opts))
	t.parseRoot()
	return t, nil
}

func (t *Tree) ParseCustom(input string, opts ParseOpts) (tree *Tree, err error) {
	defer t.recover(&err)
	t.opts = opts
	t.startParse(lex(input, opts))
	if opts.TreatRootAsList {
		t.parseRootList()
	} else {
		t.parseRoot()
	}
	return t, nil
}

// parseRoot pretends root is regular text, not a list, for
// compatibility with traditional brace expansion.
func (t *Tree) parseRoot() {
	ln := t.newListNode()
	t.Root = &ln

	pn := &PhraseNode{}

	for t.peek().typ != itemEOF {
		if t.peek().typ == itemText || t.peek().typ == itemOpen {
			pn.append(t.exprOrText())
		} else if t.peek().typ == itemSeparator {
			t.next()
			pn.append(t.newTextNode(t.opts.Separator))
		} else {
			t.errorf("unexpected item type %v in root", t.peek().typ)
		}
	}

	if t.Root != nil {
		t.Root.append(*pn)
	}
}

// parseRootList is essentially the same as list, except it runs until EOF.
func (t *Tree) parseRootList() {
	ln := t.newListNode()
	t.Root = &ln

	if t.peek().typ == itemSeparator || t.peek().typ == itemEOF {
		t.Root.append(t.newEmptyPhraseNode())
	}

	for t.peek().typ != itemEOF {
		if t.peek().typ == itemText || t.peek().typ == itemOpen {
			t.Root.append(t.phrase())

		} else if t.peek().typ == itemSeparator {
			t.next()

			if t.peek().typ == itemSeparator || t.peek().typ == itemEOF {
				t.Root.append(t.newEmptyPhraseNode())
			}
		} else {
			t.errorf("unexpected item type %v in root", t.peek().typ)
		}
	}
}

func (t *Tree) list() ListNode {
	ln := t.newListNode()

	if t.peek().typ == itemSeparator {
		ln.append(t.newEmptyPhraseNode())
	}

	for t.peek().typ != itemClose {
		if t.peek().typ == itemText || t.peek().typ == itemOpen {
			ln.append(t.phrase())
		} else if t.peek().typ == itemSeparator {
			t.next()
			if t.peek().typ == itemSeparator || t.peek().typ == itemClose {
				ln.append(t.newEmptyPhraseNode())
			}
		} else {
			t.errorf("unexpected item type %v in list", t.peek().typ)
		}
	}

	t.next()

	return ln
}

func (t *Tree) phrase() PhraseNode {
	pn := PhraseNode{}

	for t.peek().typ == itemText || t.peek().typ == itemOpen {
		pn.append(t.exprOrText())
	}

	return pn
}

func (t *Tree) exprOrText() Node {
	switch t.peek().typ {
	case itemText:
		return t.text()
	case itemOpen:
		t.next()
		return t.list()
	default:
		t.errorf("unexpected item type %v expression or text", t.peek().typ)
	}

	panic("not reached")
}

func (t *Tree) text() TextNode {
	return t.newTextNode(t.next().val)
}

func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token = t.lex.nextItem()
	}
	return t.token
}

func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token
	}
	t.peekCount = 1
	t.token = t.lex.nextItem()
	return t.token
}

func New() *Tree {
	return &Tree{}
}

func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	panic(fmt.Errorf(format, args...))
}

func (t *Tree) error(err error) {
	t.errorf("%s", err)
}

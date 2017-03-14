package braceexpansion

import (
	"fmt"
	"testing"
)

var itemName = map[itemType]string{
	itemError:     "error",
	itemOpen:      "open",
	itemClose:     "close",
	itemSeparator: "separator",
	itemText:      "text",
	itemEOF:       "EOF",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

func (i item) String() string {
	return fmt.Sprintf("%v:\"%s\"", i.typ, i.val)
}

type lexTest struct {
	input string
	items []item
}

var lexTests = []lexTest{
	{"abc", []item{
		item{itemText, "abc"},
		item{itemEOF, ""},
	}},
	{"def", []item{
		item{itemText, "def"},
		item{itemEOF, ""},
	}},
	{"{{", []item{
		item{itemOpen, "{"},
		item{itemOpen, "{"},
		item{itemEOF, ""},
	}},
	{"{", []item{
		item{itemOpen, "{"},
		item{itemEOF, ""},
	}},
	{"}", []item{
		item{itemClose, "}"},
		item{itemEOF, ""},
	}},
	{",", []item{
		item{itemSeparator, ","},
		item{itemEOF, ""},
	}},
	{"a,", []item{
		item{itemText, "a"},
		item{itemSeparator, ","},
		item{itemEOF, ""},
	}},
	{",a", []item{
		item{itemSeparator, ","},
		item{itemText, "a"},
		item{itemEOF, ""},
	}},
	{"{,", []item{
		item{itemOpen, "{"},
		item{itemSeparator, ","},
		item{itemEOF, ""},
	}},
	{",,", []item{
		item{itemSeparator, ","},
		item{itemSeparator, ","},
		item{itemEOF, ""},
	}},
	{"a,b", []item{
		item{itemText, "a"},
		item{itemSeparator, ","},
		item{itemText, "b"},
		item{itemEOF, ""},
	}},
	{"{a,b}", []item{
		item{itemOpen, "{"},
		item{itemText, "a"},
		item{itemSeparator, ","},
		item{itemText, "b"},
		item{itemClose, "}"},
		item{itemEOF, ""},
	}},
	{"{a{1,2},b}", []item{
		item{itemOpen, "{"},
		item{itemText, "a"},
		item{itemOpen, "{"},
		item{itemText, "1"},
		item{itemSeparator, ","},
		item{itemText, "2"},
		item{itemClose, "}"},
		item{itemSeparator, ","},
		item{itemText, "b"},
		item{itemClose, "}"},
		item{itemEOF, ""},
	}},
	{"{a,b}x{1,2}", []item{
		item{itemOpen, "{"},
		item{itemText, "a"},
		item{itemSeparator, ","},
		item{itemText, "b"},
		item{itemClose, "}"},
		item{itemText, "x"},
		item{itemOpen, "{"},
		item{itemText, "1"},
		item{itemSeparator, ","},
		item{itemText, "2"},
		item{itemClose, "}"},
		item{itemEOF, ""},
	}},
}

func TestLex(t *testing.T) {
	for _, lt := range lexTests {
		t.Run(lt.input, func(t *testing.T) {
			items := collect(lt.input)
			if !equal(items, lt.items) {
				t.Errorf("want\n%+v\nhave\n%+v", lt.items, items)
			}
		})
	}
}

func collect(input string) (items []item) {
	l := lex(input, ParseOpts{OpenBrace: "{", CloseBrace: "}", Separator: ","})
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(a, b []item) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].typ != b[i].typ {
			return false
		}
		if a[i].val != b[i].val {
			return false
		}
	}
	return true
}

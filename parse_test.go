package braceexpansion

import (
	"fmt"
	"testing"
)

func (l ListNode) String() string {
	return fmt.Sprintf("List: %v", l.Phrases)
}

func (p PhraseNode) String() string {
	return fmt.Sprintf("Phrase: %v", p.Parts)
}

func (t TextNode) String() string {
	return fmt.Sprintf("\"%s\"", t.text)
}

type parseTest struct {
	input string
	ok    bool
	debug string
}

var parseTests = []parseTest{
	{"abc", true, `List: [Phrase: ["abc"]]`},
	{"a,b", true, `List: [Phrase: ["a" "," "b"]]`},
	{"a,,b", true, `List: [Phrase: ["a" "," "," "b"]]`},
	{"a,", true, `List: [Phrase: ["a" ","]]`},
	{",a", true, `List: [Phrase: ["," "a"]]`},
	{",a,", true, `List: [Phrase: ["," "a" ","]]`},
	{",", true, `List: [Phrase: [","]]`},
	{",,", true, `List: [Phrase: ["," ","]]`},
	{",,,", true, `List: [Phrase: ["," "," ","]]`},
	{"{,}", true, `List: [Phrase: [List: [Phrase: [""] Phrase: [""]]]]`},
	{"{,,}", true, `List: [Phrase: [List: [Phrase: [""] Phrase: [""] Phrase: [""]]]]`},
	{"{,,,}", true, `List: [Phrase: [List: [Phrase: [""] Phrase: [""] Phrase: [""] Phrase: [""]]]]`},
	{"abc,def", true, `List: [Phrase: ["abc" "," "def"]]`},
	{"{}", true, `List: [Phrase: [List: []]]`},
	{"a{}", true, `List: [Phrase: ["a" List: []]]`},
	{"a{,}", true, `List: [Phrase: ["a" List: [Phrase: [""] Phrase: [""]]]]`},
	{"{abc,}", true, `List: [Phrase: [List: [Phrase: ["abc"] Phrase: [""]]]]`},
	{"{,abc}", true, `List: [Phrase: [List: [Phrase: [""] Phrase: ["abc"]]]]`},
	{"{abc}", true, `List: [Phrase: [List: [Phrase: ["abc"]]]]`},
	{"{a{}}", true, `List: [Phrase: [List: [Phrase: ["a" List: []]]]]`},
	{"{abc,def}", true, `List: [Phrase: [List: [Phrase: ["abc"] Phrase: ["def"]]]]`},
	{"{a,b}{1,2}", true, `List: [Phrase: [List: [Phrase: ["a"] Phrase: ["b"]] List: [Phrase: ["1"] Phrase: ["2"]]]]`},
	{"{abc}def", true, `List: [Phrase: [List: [Phrase: ["abc"]] "def"]]`},
	{"}", false, ``},
	{"}}", false, ``},
	{"{{}", false, ``},
	{"{}}", false, ``},
	{"{,abc", false, ``},
	{"{abc,def", false, ``},
}

var parseTestsCustom = []parseTest{
	{"abc", true, `List: [Phrase: ["abc"]]`},
	{"a,b", true, `List: [Phrase: ["a"] Phrase: ["b"]]`},
	{"a,,b", true, `List: [Phrase: ["a"] Phrase: [""] Phrase: ["b"]]`},
	{"a,", true, `List: [Phrase: ["a"] Phrase: [""]]`},
	{",a", true, `List: [Phrase: [""] Phrase: ["a"]]`},
	{",a,", true, `List: [Phrase: [""] Phrase: ["a"] Phrase: [""]]`},
	{",", true, `List: [Phrase: [""] Phrase: [""]]`},
	{",,", true, `List: [Phrase: [""] Phrase: [""] Phrase: [""]]`},
	{",,,", true, `List: [Phrase: [""] Phrase: [""] Phrase: [""] Phrase: [""]]`},
	{"(,)", true, `List: [Phrase: [List: [Phrase: [""] Phrase: [""]]]]`},
	{"(,,)", true, `List: [Phrase: [List: [Phrase: [""] Phrase: [""] Phrase: [""]]]]`},
	{"(,,,)", true, `List: [Phrase: [List: [Phrase: [""] Phrase: [""] Phrase: [""] Phrase: [""]]]]`},
	{"abc,def", true, `List: [Phrase: ["abc"] Phrase: ["def"]]`},
	{"()", true, `List: [Phrase: [List: []]]`},
	{"a()", true, `List: [Phrase: ["a" List: []]]`},
	{"a(,)", true, `List: [Phrase: ["a" List: [Phrase: [""] Phrase: [""]]]]`},
	{"(abc,)", true, `List: [Phrase: [List: [Phrase: ["abc"] Phrase: [""]]]]`},
	{"(,abc)", true, `List: [Phrase: [List: [Phrase: [""] Phrase: ["abc"]]]]`},
	{"(abc)", true, `List: [Phrase: [List: [Phrase: ["abc"]]]]`},
	{"(a())", true, `List: [Phrase: [List: [Phrase: ["a" List: []]]]]`},
	{"(abc,def)", true, `List: [Phrase: [List: [Phrase: ["abc"] Phrase: ["def"]]]]`},
	{"(a,b)(1,2)", true, `List: [Phrase: [List: [Phrase: ["a"] Phrase: ["b"]] List: [Phrase: ["1"] Phrase: ["2"]]]]`},
	{"(abc)def", true, `List: [Phrase: [List: [Phrase: ["abc"]] "def"]]`},
	{")", false, ``},
	{"))", false, ``},
	{"(()", false, ``},
	{"())", false, ``},
	{"(,abc", false, ``},
	{"(abc,def", false, ``},
}

func TestParse(t *testing.T) {
	testParse(t, parseTests, parse)
}

func TestParseCustom(t *testing.T) {
	testParse(t, parseTestsCustom, parseCustom)
}

func testParse(t *testing.T, tests []parseTest, f parseFunc) {
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tree, err := f(test.input)

			switch {
			case err == nil && !test.ok:
				t.Error("Expected error, got none")
			case err != nil && test.ok:
				t.Errorf("Unexpected error: %v", err)
			case err == nil && fmt.Sprintf("%v", tree.Root) != test.debug:
				t.Errorf("Unexpected tree.\nWant:\n%v\nHave:\n%v", test.debug, fmt.Sprintf("%v", tree.Root))
			default:
				if err != nil {
					t.Logf("\"%s\" => error", test.input)
				} else {
					t.Logf("\"%s\" => %v", test.input, tree.Root)
				}
			}
		})
	}
}

type parseFunc func(input string) (*Tree, error)

func parse(input string) (*Tree, error) {
	return New().Parse(input)
}

func parseCustom(input string) (*Tree, error) {
	opts := ParseOpts{OpenBrace: "(", CloseBrace: ")", Separator: ",", TreatRootAsList: true, TreatSingleAsOptional: true}
	return New().ParseCustom(input, opts)
}

package braceexpansion

import (
	"fmt"
	"testing"

	"github.com/thomasheller/slicecmp"
)

type expandTest struct {
	input  string
	output []string
}

var expandTests = []expandTest{
	{"a", []string{"a"}},
	{"a,b", []string{"a,b"}},
	{"a,,b", []string{"a,,b"}},
	{"a,", []string{"a,"}},
	{",a", []string{",a"}},
	{",", []string{","}},
	{",,", []string{",,"}},
	{"abc,def", []string{"abc,def"}},
	{"{abc,def}", []string{"abc", "def"}},
	{"{}", []string{"{}"}},
	{"a{}", []string{"a{}"}},
	{"{a{}}", []string{"{a{}}"}}, // brace expansion default
	{"a{,}", []string{"a", "a"}},
	{"a{b,c}", []string{"ab", "ac"}},
	{"a{,b,c}", []string{"a", "ab", "ac"}},
	{"a{b,c,}", []string{"ab", "ac", "a"}},
	{"a{b,d{e,f}}", []string{"ab", "ade", "adf"}},
	{"a{b,{c,d}}", []string{"ab", "ac", "ad"}},
	{"{a,b}{1,2}", []string{"a1", "a2", "b1", "b2"}},
	{"{a,b}x{1,2}", []string{"ax1", "ax2", "bx1", "bx2"}},
	{"{abc}", []string{"{abc}"}},                 // brace expansion default
	{"{abc}def", []string{"{abc}def"}},           // brace expansion default
	{"{abc{def}}ghi", []string{"{abc{def}}ghi"}}, // brace expansion default
	{"{,}", []string{"", ""}},
	{"{,,}", []string{"", "", ""}},
	{"{,,,}", []string{"", "", "", ""}},
	{"{a,{{{b}}}}", []string{"a", "{{{b}}}"}},
	{"{a{1,2}b}", []string{"{a1b}", "{a2b}"}},
}

var expandTestsCustom = []expandTest{
	{"a", []string{"a"}},
	{"a,b", []string{"a", "b"}},
	{"a,,b", []string{"a", "", "b"}},
	{"a,", []string{"a", ""}},
	{",a", []string{"", "a"}},
	{",", []string{"", ""}},
	{",,", []string{"", "", ""}},
	{"abc,def", []string{"abc", "def"}},
	{"(abc,def)", []string{"abc", "def"}},
	{"()", []string{"()"}},
	{"a()", []string{"a()"}},
	{"(a())", []string{"a()", ""}}, // single-as-optional mode
	{"a(,)", []string{"a", "a"}},
	{"a(b,c)", []string{"ab", "ac"}},
	{"a(,b,c)", []string{"a", "ab", "ac"}},
	{"a(b,c,)", []string{"ab", "ac", "a"}},
	{"a(b,d(e,f))", []string{"ab", "ade", "adf"}},
	{"a(b,(c,d))", []string{"ab", "ac", "ad"}},
	{"(a,b)(1,2)", []string{"a1", "a2", "b1", "b2"}},
	{"(a,b)x(1,2)", []string{"ax1", "ax2", "bx1", "bx2"}},
	{"(abc)", []string{"abc", ""}},                            // single-as-optional mode
	{"(abc)def", []string{"abcdef", "def"}},                   // single-as-optional mode
	{"(abc(def))ghi", []string{"abcdefghi", "abcghi", "ghi"}}, // single-as-optional mode
	{"(,)", []string{"", ""}},
	{"(,,)", []string{"", "", ""}},
	{"(,,,)", []string{"", "", "", ""}},
	{"(a,(((b))))", []string{"a", "b", "", "", ""}}, // single-as-optional mode
	{"(a(1,2)b)", []string{"a1b", "a2b", ""}},       // single-as-optional mode
}

func TestExpand(t *testing.T) {
	testExpand(t, expandTests, parse)
}

func TestExpandCustom(t *testing.T) {
	testExpand(t, expandTestsCustom, parseCustom)
}

func testExpand(t *testing.T, tests []expandTest, f parseFunc) {
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tree, err := f(test.input)

			if err != nil {
				t.Errorf("Parse error: %v", err)
			}

			output := tree.Expand()

			if !slicecmp.Equal(test.output, output) {
				t.Errorf("Unexpected output:\n%s", slicecmp.Sprint([]string{"want", "have"}, test.output, output))
			}
		})
	}
}

func TestExpandTree(t *testing.T) {
	tree := &Tree{}
	ln := tree.newListNode()
	tree.Root = &ln

	pn := tree.newPhraseNode()

	ln1 := tree.newListNode()
	ln2 := tree.newListNode()

	ln1.append(tree.newPhraseNodeWithText("a"))
	ln1.append(tree.newPhraseNodeWithText("b"))
	ln2.append(tree.newPhraseNodeWithText("1"))
	ln2.append(tree.newPhraseNodeWithText("2"))

	pn.append(ln1)
	pn.append(ln2)

	ln.append(pn)

	output := tree.Expand()
	expected := []string{"a1", "a2", "b1", "b2"}

	if !slicecmp.Equal(expected, output) {
		t.Errorf("Unexpected output:\n%s", slicecmp.Sprint([]string{"want", "have"}, expected, output))
	}
}

func TestCartesian(t *testing.T) {
	cartesianTests := []struct {
		input  [][]string
		output []string
	}{
		{[][]string{}, []string{}},
		{[][]string{{"a", "b"}}, []string{"a", "b"}},
		{[][]string{{"a", "b"}, {"1", "2"}}, []string{"a1", "a2", "b1", "b2"}},
		{[][]string{{"a", "b"}, {"1", "2"}, {"x", "y"}}, []string{"a1x", "a1y", "a2x", "a2y", "b1x", "b1y", "b2x", "b2y"}},
		{[][]string{{"a", "b"}, {"x"}, {"1", "2"}}, []string{"ax1", "ax2", "bx1", "bx2"}},
	}

	for i, ct := range cartesianTests {
		t.Run(fmt.Sprintf("Cartesian%d", i), func(t *testing.T) {
			output := Cartesian(ct.input)
			if !slicecmp.Equal(ct.output, output) {
				t.Errorf("Unexpected output:\n%s", slicecmp.Sprint([]string{"want", "have"}, ct.output, output))
			}
		})
	}
}

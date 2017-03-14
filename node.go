package braceexpansion

type Node interface{}

type NodeType int

const (
	NodeList NodeType = iota
	NodePhrase
	NodeText
)

type ListNode struct {
	NodeType
	Phrases []PhraseNode
	Tree    *Tree
}

func (l *ListNode) append(n PhraseNode) {
	l.Phrases = append(l.Phrases, n)
}

type PhraseNode struct {
	NodeType
	Parts []Node // TextNode or ListNode
}

func (p *PhraseNode) append(n Node) { // TextNode or ListNode
	p.Parts = append(p.Parts, n)
}

func (t *Tree) newListNode() ListNode {
	return ListNode{Tree: t}
}

func (t *Tree) newPhraseNode() PhraseNode {
	return PhraseNode{}
}

func (t *Tree) newTextNode(val string) TextNode {
	return TextNode{text: val}
}

func (t *Tree) newEmptyPhraseNode() PhraseNode {
	return PhraseNode{Parts: []Node{t.newTextNode("")}}
}

func (t *Tree) newPhraseNodeWithText(val string) PhraseNode {
	return PhraseNode{Parts: []Node{t.newTextNode(val)}}
}

type TextNode struct {
	NodeType
	text string
}

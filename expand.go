package braceexpansion

func (t *Tree) Expand() []string {
	return t.Root.Expand(true)
}

func (l ListNode) Expand(root bool) []string {
	// empty brace expressions like "{}" are printed as regular text:
	if len(l.Phrases) == 0 {
		return []string{l.Tree.opts.OpenBrace + l.Tree.opts.CloseBrace}
	}

	// if ListNode has only one child, it may be treated as
	// optional depending on ParseOpts, but not at root level:
	if len(l.Phrases) == 1 {
		lines := l.Phrases[0].Expand()
		if root {
			return lines
		} else {
			if l.Tree.opts.TreatSingleAsOptional {
				lines = append(lines, "")
				return lines
			} else {
				result := []string{}
				for _, line := range lines {
					result = append(result, l.Tree.opts.OpenBrace+line+l.Tree.opts.CloseBrace)
				}
				return result
			}
		}
	}

	lines := []string{}

	for _, phrase := range l.Phrases {
		for _, line := range phrase.Expand() {
			lines = append(lines, line)
		}
	}

	return lines
}

func (p PhraseNode) Expand() []string {
	sets := [][]string{}
	for _, part := range p.Parts {
		set := p.expandPart(part)
		sets = append(sets, set)
	}

	result := Cartesian(sets)

	return result
}

func (p *PhraseNode) expandPart(part Node) []string {
	result := []string{}

	switch node := part.(type) {
	case TextNode:
		for _, line := range node.Expand() {
			result = append(result, line)
		}
	case ListNode:
		for _, line := range node.Expand(false) {
			result = append(result, line)
		}
	default:
		panic("unexpected node type")
	}

	return result
}

func (t TextNode) Expand() []string {
	return []string{t.text}
}

func Cartesian(sets [][]string) []string {
	if len(sets) == 0 {
		return []string{}
	}

	result := sets[0]

	for i := 1; i < len(sets); i++ {
		result = Cartesian2(result, sets[i])
	}

	return result
}

func Cartesian2(first, second []string) []string {
	result := []string{}

	for _, e1 := range first {
		for _, e2 := range second {
			result = append(result, e1+e2)
		}
	}

	return result
}

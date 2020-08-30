/*
	字典树，用于做脏词替换,unicode 支持
*/
package main

type trieNode struct {
	r        rune
	children []*trieNode
}

func buildTrie(words ...string) *trieNode {
	root := &trieNode{}

	for _, word := range words {
		scanNode := root

		for _, r := range word {
			scanNode = scanNode.append(r)
		}
	}

	return root
}

func (tn *trieNode) append(r rune) *trieNode {
	for _, child := range tn.children {
		if child == nil {
			continue
		}

		if child.r == r {
			return child
		}
	}

	tn.children = append(tn.children, &trieNode{
		r: r,
	})

	return tn.children[len(tn.children)-1]
}

func (tn *trieNode) replace(content string, rep rune) string {
	runes := []rune(content)

	scanNode, b := tn, -1
	for i := 0; i < len(runes); i++ {
		var findChild *trieNode
		for _, child := range scanNode.children {
			if child == nil || child.r != runes[i] {
				continue
			}

			findChild = child
			break
		}

		if findChild == nil {
			scanNode, b = tn, -1
			continue
		}

		scanNode = findChild

		if b < 0 {
			b = i
		}

		if len(findChild.children) == 0 {
			for r := b; r <= i; r++ {
				runes[r] = rep
			}
		}
	}

	return string(runes)
}

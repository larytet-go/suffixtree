package suffixtree

import (
	"fmt"
	"testing"
)

type myRune struct {
	r rune
}

func (r myRune) IsEqual(other Symbol) bool {
	return r.r == other.(myRune).r
}

func newWord(s string) []Symbol {
	symbols := []Symbol{}
	for _, c := range s {
		r := myRune{c}
		symbols = append(symbols, r)
	}
	return symbols
}

func TestSuffixTree(t *testing.T) {
	words := [][]Symbol{
		newWord("banana"),
		newWord("apple"),
		newWord("中文app"),
	}
	s := ([]Symbol)(words[0])
	tree := NewGeneralizedSuffixTree()
	for k, word := range words {
		tree.Put(word, k)
	}
	indexs := tree.Search(newWord("a"), -1)

	if len(indexs) != 3 {
		t.Error("indexs len should be 3,but ", len(indexs))
	}
	fmt.Println(indexs)
	for _, index := range indexs {
		fmt.Println(words[index])
	}

	indexs = tree.Search(newWord("文"), 0)

	if len(indexs) != 1 && indexs[0] != 2 {
		t.Error("indexs len should be 1 and indexs[0] must be 2,but ", len(indexs))
	}

	printnode("\t", tree.root)
}

func printnode(flag string, n *node) {
	for _, e := range n.edges {
		fmt.Printf("%s %v %v \n", flag, e.label, e.node.data)
		printnode(flag+"\t-", e.node)
	}
}

package suffixtree

import (
	"fmt"
	"math/rand"
	"testing"
)

type myRune struct {
	r rune
}

func (r myRune) IsEqual(other Symbol) bool {
	return r.r == other.(myRune).r
}

func (r myRune) IsLess(other Symbol) bool {
	return r.r < other.(myRune).r
}

func newWord(s string) []Symbol {
	symbols := []Symbol{}
	for _, c := range s {
		r := myRune{c}
		symbols = append(symbols, r)
	}
	return symbols
}

func toRunes(word []Symbol) []rune {
	runes := []rune{}
	for _, c := range word {
		runes = append(runes, c.(myRune).r)
	}
	return runes
}

func TestSuffixTree(t *testing.T) {
	words := [][]Symbol{
		newWord("banana"),
		newWord("apple"),
		newWord("中文app"),
	}
	tree := NewGeneralizedSuffixTree()
	for k, word := range words {
		tree.Put(word, k)
	}

	nodesCount := tree.NodesCount()
	expectedNodesCount := 16
	if nodesCount != expectedNodesCount {
		t.Errorf("Should be %d nodes instead of %d", expectedNodesCount, nodesCount)
	}

	edgesCount := tree.EdgesCount()
	expectedEdgesCount := 8
	if edgesCount != expectedEdgesCount {
		t.Errorf("Should be %d edges instead of %d", expectedEdgesCount, edgesCount)
	}

	indexs := tree.Search(newWord("a"), -1)

	if len(indexs) != 3 {
		t.Error("indexs len should be 3,but ", len(indexs))
	}
	fmt.Println(indexs)
	for _, index := range indexs {
		fmt.Println(string(toRunes(words[index])))
	}

	indexs = tree.Search(newWord("文"), 0)

	if len(indexs) != 1 && indexs[0] != 2 {
		t.Error("indexs len should be 1 and indexs[0] must be 2,but ", len(indexs))
	}

	tree.Put(words[0], 3)
	nodesCount = tree.NodesCount()
	if nodesCount != expectedNodesCount {
		t.Errorf("Should be %d nodes instead of %d", expectedNodesCount, nodesCount)
	}

	printnode("\t", tree.Root)
}

func TestSuffixTree1(t *testing.T) {
	words := [][]Symbol{
		newWord("banana"),
		newWord("apple"),
		newWord("中文app"),
	}
	tree := NewGeneralizedSuffixTree()
	k := 0
	for _, word := range words {
		for _, c := range word {
			t := []Symbol{c}
			tree.Put(t, k)
		}
		k++
	}

	printnode("\t", tree.Root)
}

func randomWord(size int) []Symbol {
	symbols := []Symbol{}
	for i := 0; i < size; i++ {
		c := 'a' + rand.Intn(0xFFFFFF)
		r := myRune{rune(c)}
		symbols = append(symbols, r)
	}
	return symbols
}

func TestSuffixTreeRandom(t *testing.T) {
	rand.Seed(0)
	tree := NewGeneralizedSuffixTree()
	for k := 0; k < 100; k++ {
		word := randomWord(8)
		tree.Put(word, k)
	}
	rand.Seed(0)
	for k := 0; k < 500; k++ {
		word := randomWord(8)
		found := tree.Search(word, 0)
		if len(found) == 0 {
			t.Errorf("Not found %v", word)
		}
	}
}

func printnode(flag string, n *Node) {
	for _, e := range n.Edges {
		fmt.Printf("%s %v %v \n", flag, string(toRunes(e.Label)), e.Node.Data)
		printnode(flag+"\t-", e.Node)
	}
}

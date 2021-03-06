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
	if size == 0 {
		size = 4 + rand.Intn(8)
	}

	symbols := []Symbol{}
	for i := 0; i < size; i++ {
		c := 'a' + rand.Intn(0xFFFFFF)
		r := myRune{rune(c)}
		symbols = append(symbols, r)
	}
	return symbols
}

func TestSuffixTreeRandom(t *testing.T) {
	count := 500
	tree := NewGeneralizedSuffixTree()
	rand.Seed(0)
	for k := 0; k < count; k++ {
		word := randomWord(0)
		tree.Put(word, k)
		found := tree.Search(word, 0)
		if len(found) != 1 {
			t.Errorf("Not found word=%v,found=%v", word, found)
		}
	}
	rand.Seed(0)
	for k := 0; k < count; k++ {
		word := randomWord(0)
		found := tree.Search(word, 0)
		if len(found) == 0 {
			t.Errorf("Not found %v", word)
		}
	}
	word := randomWord(0)
	found := tree.Search(word, 0)
	if len(found) != 0 {
		t.Errorf("Found %v", word)
	}
	edgesCount := tree.EdgesCount()
	if edgesCount != 3747 {
		t.Errorf("Edges count %d", edgesCount)
	}
	// Try again to check that the code can process the same word
	// more than once
	rand.Seed(0)
	for k := 0; k < count; k++ {
		word := randomWord(0)
		tree.Put(word, k)
	}
}

func printnode(flag string, n *Node) {
	for _, e := range n.Edges {
		fmt.Printf("%s %v %v \n", flag, string(toRunes(e.Label)), e.Node.Data)
		printnode(flag+"\t-", e.Node)
	}
}

func poulateTree(count int) *GeneralizedSuffixTree {
	tree := NewGeneralizedSuffixTree()
	rand.Seed(0)
	for k := 0; k < count; k++ {
		word := randomWord(0)
		tree.Put(word, k)
	}
	return tree
}

func BenchmarkSearch(b *testing.B) {
	count := 100
	tree := poulateTree(count)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rand.Seed(0)
		for k := 0; k < count; k++ {
			word := randomWord(0)
			found := tree.Search(word, 0)
			if len(found) == 0 {
				b.Errorf("Not found %v", word)
			}
		}
	}
}


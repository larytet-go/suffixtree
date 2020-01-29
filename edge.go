package suffixtree

// Symbol is an abstract type. An edge label in the suffix tree
// is an array of symbols
type Symbol interface {
	IsEqual(other Symbol) bool
	IsLess(other Symbol) bool
}

func isEqual(s1 []Symbol, s2 []Symbol) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, _s1 := range s1 {
		_s2 := s2[i]
		if !_s1.IsEqual(_s2) {
			return false
		}
	}
	return true
}

func indexSymbol(s []Symbol, c Symbol) int {
	for _, _s := range s {
		if _s.IsEqual(c) {
			return 0
		}
	}
	return -1
}

// Loosely follow https://golang.org/src/strings/strings.go
func indexOf(s, substr []Symbol) int {
	n := len(substr)
	switch {
	case n == 0:
		return 0
	case n == 1:
		return indexSymbol(s, substr[0])
	case n == len(s):
		if isEqual(s, substr) {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	}
	// TODO I can compare faster by skipping already performed
	// comparisons or Karp–Rabin algorithm
	for i := 0; i < len(s)-len(substr)+1; i++ {
		isEqual := true
		for j := 0; j < len(substr); j++ {
			if !s[j+i].IsEqual(substr[j]) {
				isEqual = false
				break
			}
		}
		if isEqual {
			return i
		}
	}
	return -1
}

type Edge struct {
	Label []Symbol
	*Node
}

func newEdge(label []Symbol, node *Node) *Edge {
	return &Edge{Label: label, Node: node}
}

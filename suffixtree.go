// Package suffixtree implements A Generalized Suffix Tree, based on the Ukkonen's paper "On-line construction of suffix trees"
package suffixtree

type generalizedSuffixTree struct {
	root       *node //The root of the suffix tree
	activeLeaf *node //The last leaf that was added during the update operation
	edgesCount int   //Number of edges in the tree
}

// Search search for the given word within the GST and returns at most the given number of matches.
// numElments <= 0 get all matches
func (t *generalizedSuffixTree) Search(word []Symbol, numElements int) []int {
	node := t.searchNode(word)
	if node == nil {
		return nil
	}
	return node.getData(numElements)
}

func (t *generalizedSuffixTree) GetEdgesCount() int {
	return t.edgesCount
}

// searchNode returns the tree node (if present) that corresponds to the given string.
func (t *generalizedSuffixTree) searchNode(word []Symbol) *node {
	/*
	 * Verifies if exists a path from the root to a node such that the concatenation
	 * of all the labels on the path is a superstring of the given word.
	 * If such a path is found, the last node on it is returned.
	 */
	var currentNode = t.root
	var currentEdge *edge

	for i, symbol := range word {
		currentEdge = currentNode.getEdge(symbol)
		if currentEdge == nil {
			// there is no edge starting with this symbol
			return nil
		} else {
			label := currentEdge.label
			lenToMatch := len(word) - i
			if lenToMatch > len(label) {
				lenToMatch = len(label)
			}
			if !isEqual(word[i:i+lenToMatch], label[:lenToMatch]) {
				// the label on the edge does not correspond to the one in the string to search
				return nil
			}

			if len(label) >= len(word)-i {
				return currentEdge.node
			} else {
				// advance to next node
				currentNode = currentEdge.node
				i += lenToMatch
			}
		}
	}

	return nil
}

// Put adds the specified index to the GST under the given key.
func (t *generalizedSuffixTree) Put(symbols []Symbol, index int) {
	// reset activeLeaf
	t.activeLeaf = t.root
	s := t.root

	// proceed with tree construction (closely related to procedure in
	// Ukkonen's paper)
	var text []Symbol
	// iterate over the string, one rune at a time
	for k, r := range symbols {
		// line 6
		text = append(text, r)
		// line 7: update the tree with the new transitions due to this new rune
		s, text = t.update(s, text, symbols[k:], index)
		// line 8: make sure the active pair is canonical
		s, text = t.canonize(s, text)
	}

	// add leaf suffix link, is necessary
	if t.activeLeaf.suffix == nil && t.activeLeaf != t.root && t.activeLeaf != s {
		t.activeLeaf.suffix = s
	}
}

/*
 * update updates the tree starting from inputNode and by adding stringPart.
 *
 * Returns a reference (*node,[]rune) pair for the string that has been added so far.
 * This means:
 * - the Node will be the Node that can be reached by the longest path string (S1)
 *   that can be obtained by concatenating consecutive edges in the tree and
 *   that is a substring of the string added so far to the tree.
 * - the String will be the remainder that must be added to S1 to get the string
 *   added so far.
 *
 * @param inputNode the node to start from
 * @param stringPart the string to add to the tree
 * @param rest the rest of the string
 * @param value the value to add to the index
 */
func (t *generalizedSuffixTree) update(inputNode *node, stringPart []Symbol, rest []Symbol, value int) (s *node, symbols []Symbol) {
	s = inputNode
	symbols = stringPart
	newSymbol := stringPart[len(stringPart)-1]

	// line 1
	oldroot := t.root

	// line 1b
	endpoint, r := t.testAndSplit(s, stringPart[:len(stringPart)-1], newSymbol, rest, value)

	var leaf *node
	// line 2
	for !endpoint {
		// line 3
		tempEdge := r.getEdge(newSymbol)
		if tempEdge != nil {
			// such a node is already present. This is one of the main differences from Ukkonen's case:
			// the tree can contain deeper nodes at this stage because different strings were added by previous iterations.
			leaf = tempEdge.node
		} else {
			// must build a new leaf
			leaf = newNode()
			leaf.addRef(value)
			newedge := newEdge(rest, leaf)
			r.addEdge(newSymbol, newedge)
			t.edgesCount++
		}

		// update suffix link for newly created leaf
		if t.activeLeaf != t.root {
			t.activeLeaf.suffix = leaf
		}
		t.activeLeaf = leaf

		// line 4
		if oldroot != t.root {
			oldroot.suffix = r
		}

		// line 5
		oldroot = r

		// line 6
		if s.suffix == nil { // root node
			// this is a special case to handle what is referred to as node _|_ on the paper
			symbols = symbols[1:]
		} else {
			n, b := t.canonize(s.suffix, safeCutLastChar(symbols))
			s = n
			// use intern to ensure that runes is a reference from the string pool
			symbols = append(b, symbols[len(symbols)-1])
		}

		// line 7
		endpoint, r = t.testAndSplit(s, safeCutLastChar(symbols), newSymbol, rest, value)
	}

	// line 8
	if oldroot != t.root {
		oldroot.suffix = r
	}

	return
}

/*
 * canonize return a (*node, []rune) (n, remainder) pair such that n is a farthest descendant of
 * s (the input node) that can be reached by following a path of edges denoting
 * a prefix of inputstr and remainder will be string that must be
 * appended to the concatenation of labels from s to n to get inpustr.
 */
func (t *generalizedSuffixTree) canonize(s *node, symbols []Symbol) (*node, []Symbol) {

	currentNode := s
	if len(symbols) > 0 {
		g := s.getEdge(symbols[0])
		// descend the tree as long as a proper label is found
		for g != nil && indexOf(symbols, g.label) == 0 {
			symbols = symbols[len(g.label):]
			currentNode = g.node
			if len(symbols) > 0 {
				g = currentNode.getEdge(symbols[0])
			}
		}
	}
	return currentNode, symbols
}

/*
 * testAndSplit tests whether the string stringPart + r is contained in the subtree that has inputs as root.
 * If that's not the case, and there exists a path of edges e1, e2, ... such that
 *     e1.label + e2.label + ... + $end = stringPart
 * and there is an edge g such that
 *     g.label = stringPart + rest
 *
 * Then g will be split in two different edges, one having $end as label, and the other one
 * having rest as label.
 *
 * @param inputs the starting node
 * @param stringPart the string to search
 * @param r the following character
 * @param remainder the remainder of the string to add to the index
 * @param value the value to add to the index
 * @return a pair containing
 *                  true/false depending on whether (stringPart + t) is contained in the subtree starting in inputs
 *                  the last node that can be reached by following the path denoted by stringPart starting from inputs
 *
 */
func (t *generalizedSuffixTree) testAndSplit(inputs *node, stringPart []Symbol, r Symbol, remainder []Symbol, value int) (bool, *node) {
	// descend the tree as far as possible
	s, str := t.canonize(inputs, stringPart)

	if len(str) > 0 {
		g := s.getEdge(str[0])

		// must see whether "str" is substring of the label of an edge
		if len(g.label) > len(str) && g.label[len(str)] == r {
			return true, s
		} else {
			// need to split the edge
			newlabel := g.label[len(str):]

			// build a new node
			w := newNode()
			// build a new edge
			newedge := newEdge(str, w)
			s.addEdge(str[0], newedge)

			// link s -> r
			g.label = newlabel
			w.addEdge(newlabel[0], g)

			t.edgesCount += 2

			return false, w
		}
	} else {
		e := s.getEdge(r)
		if e == nil {
			// if there is no t-transtion from s
			return false, s
		} else {
			if isEqual(remainder, e.label) {
				// update payload of destination node
				e.node.addRef(value)
				return true, s
			} else if indexOf(remainder, e.label) == 0 {
				return true, s
			} else if indexOf(e.label, remainder) == 0 {
				// need to split as above
				newNode := newNode()
				newNode.addRef(value)
				newEdge := newEdge(remainder, newNode)
				s.addEdge(r, newEdge)

				e.label = e.label[len(remainder):]
				newNode.addEdge(e.label[0], e)
				t.edgesCount += 2
				return false, s
			} else {
				// they are different words. No prefix. but they may still share some common substr
				return true, s
			}
		}
	}

}

func safeCutLastChar(symbols []Symbol) []Symbol {
	if len(symbols) == 0 {
		return nil
	}
	return symbols[:len(symbols)-1]
}

func NewGeneralizedSuffixTree() *generalizedSuffixTree {
	t := &generalizedSuffixTree{}
	t.root = newNode()
	t.activeLeaf = t.root
	return t
}

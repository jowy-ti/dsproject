package pkg

// Node represents a single node/edge in the Patricia Trie.
type Node struct {
	prefix string
	isWord bool
	edges  []*Node
}

// Trie represents the root of the Patricia Trie.
type Trie struct {
	root *Node
}

// NewTrie initializes a new empty Patricia Trie.
func NewTrie() *Trie {
	return &Trie{
		root: &Node{},
	}
}

// Insert adds a new word to the trie.
func (t *Trie) Insert(word string) {
	if word == "" {
		return
	}
	t.root.insert(word)
}

func (n *Node) insert(word string) {
	// Check against existing edges
	for _, edge := range n.edges {
		commonLen := getCommonPrefixLength(edge.prefix, word)

		// If there is a shared prefix, we found the path to go down
		if commonLen > 0 {
			if commonLen == len(edge.prefix) {
				if commonLen == len(word) {
					// The word matches the edge exactly; just mark it as a valid word.
					edge.isWord = true
				} else {
					// The word is longer than the edge. Recursively insert the remainder.
					edge.insert(word[commonLen:])
				}
			} else {
				// The word and the edge share a prefix, but diverge. We must split the edge.

				// 1. Create a new node for the remainder of the existing edge
				remainderEdge := &Node{
					prefix: edge.prefix[commonLen:],
					isWord: edge.isWord,
					edges:  edge.edges,
				}

				// 2. Modify the current edge to stop at the shared prefix
				edge.prefix = edge.prefix[:commonLen]
				edge.edges = []*Node{remainderEdge}
				edge.isWord = false // Temporarily false, updated below

				// 3. Handle the remainder of the new word
				if commonLen == len(word) {
					// The new word ends exactly at the split point
					edge.isWord = true
				} else {
					// There is more to the new word, add it as a sibling to the remainder edge
					newEdge := &Node{
						prefix: word[commonLen:],
						isWord: true,
					}
					edge.edges = append(edge.edges, newEdge)
				}
			}
			return
		}
	}

	// No common prefix found with any existing edge. Create a completely new branch.
	newNode := &Node{
		prefix: word,
		isWord: true,
	}
	n.edges = append(n.edges, newNode)
}

// Search checks if an exact word exists in the trie.
func (t *Trie) Search(word string) bool {
	if word == "" {
		return false
	}
	return t.root.search(word)
}

func (n *Node) search(word string) bool {
	for _, edge := range n.edges {
		commonLen := getCommonPrefixLength(edge.prefix, word)

		if commonLen > 0 {
			if commonLen == len(edge.prefix) {
				if commonLen == len(word) {
					// Exact match found
					return edge.isWord
				}
				// The word is longer, keep searching down this edge
				return edge.search(word[commonLen:])
			}
			// Word matches part of the edge, but edge is longer, so word doesn't exist
			return false
		}
	}
	return false
}

// StartsWith checks if any word in the trie starts with the given prefix.
func (t *Trie) StartsWith(prefix string) bool {
	if prefix == "" {
		return true
	}
	return t.root.startsWith(prefix)
}

func (n *Node) startsWith(prefix string) bool {
	for _, edge := range n.edges {
		commonLen := getCommonPrefixLength(edge.prefix, prefix)

		if commonLen > 0 {
			if commonLen == len(prefix) {
				// We've matched the whole search prefix
				return true
			}
			if commonLen == len(edge.prefix) {
				// Matched the whole edge, keep checking the remainder
				return edge.startsWith(prefix[commonLen:])
			}
			// The edge diverges before the prefix is fully matched
			return false
		}
	}
	return false
}

// Helper to find how many characters match at the start of two strings.
func getCommonPrefixLength(s1, s2 string) int {
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}
	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}
	return minLen
}

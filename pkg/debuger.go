package pkg

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func PrintTrie(node Node) {
	printNode(node, "", true)
}

func printNode(
	node Node,
	prefix string,
	isLast bool,
) {

	if node == nil {
		fmt.Printf("%s<nil>\n", prefix)
		return
	}

	// Tree formatting
	connector := "├── "
	nextPrefix := prefix + "│   "

	if isLast {
		connector = "└── "
		nextPrefix = prefix + "    "
	}

	switch n := node.(type) {

	// --------------------------------------------------------
	// LEAF
	// --------------------------------------------------------

	case *Leaf:

		fmt.Printf(
			"%s%sLeaf\n",
			prefix,
			connector,
		)

		fmt.Printf(
			"%s    Path  : %v\n",
			nextPrefix,
			n.Path_unique,
		)

		fmt.Printf(
			"%s    Value : %s\n",
			nextPrefix,
			string(n.Value),
		)

		fmt.Printf(
			"%s    Hash  : %s\n",
			nextPrefix,
			shortHash(n.Hash()),
		)

	// --------------------------------------------------------
	// EXTENSION
	// --------------------------------------------------------

	case *Extension:

		fmt.Printf(
			"%s%sExtension\n",
			prefix,
			connector,
		)

		fmt.Printf(
			"%s    Shared Path : %v\n",
			nextPrefix,
			n.Path_shared,
		)

		fmt.Printf(
			"%s    Hash        : %s\n",
			nextPrefix,
			shortHash(n.Hash()),
		)

		printNode(
			n.next_branch,
			nextPrefix,
			true,
		)

	// --------------------------------------------------------
	// BRANCH
	// --------------------------------------------------------

	case *Branch:

		fmt.Printf(
			"%s%sBranch\n",
			prefix,
			connector,
		)

		fmt.Printf(
			"%s    Hash : %s\n",
			nextPrefix,
			shortHash(n.Hash()),
		)

		// Count children
		childIndexes := make([]int, 0)

		for i, child := range n.childs {
			if child != nil {
				childIndexes = append(childIndexes, i)
			}
		}

		for i, index := range childIndexes {

			child := n.childs[index]

			lastChild := i == len(childIndexes)-1

			label := fmt.Sprintf(
				"%s[%X]",
				nextPrefix,
				index,
			)

			printNodeWithLabel(
				child,
				label,
				lastChild,
			)
		}

	default:

		fmt.Printf(
			"%s%sUnknown node type %T\n",
			prefix,
			connector,
			node,
		)
	}
}

func printNodeWithLabel(
	node Node,
	prefix string,
	isLast bool,
) {

	indent := strings.Repeat(" ", 4)

	switch n := node.(type) {

	case *Leaf:

		fmt.Printf(
			"%s ── Leaf\n",
			prefix,
		)

		fmt.Printf(
			"%s%sPath  : %v\n",
			prefix,
			indent,
			n.Path_unique,
		)

		fmt.Printf(
			"%s%sValue : %s\n",
			prefix,
			indent,
			string(n.Value),
		)

		fmt.Printf(
			"%s%sHash  : %s\n",
			prefix,
			indent,
			shortHash(n.Hash()),
		)

	case *Extension:

		fmt.Printf(
			"%s ── Extension\n",
			prefix,
		)

		fmt.Printf(
			"%s%sShared Path : %v\n",
			prefix,
			indent,
			n.Path_shared,
		)

		fmt.Printf(
			"%s%sHash        : %s\n",
			prefix,
			indent,
			shortHash(n.Hash()),
		)

		printNode(
			n.next_branch,
			prefix+indent,
			true,
		)

	case *Branch:

		fmt.Printf(
			"%s ── Branch\n",
			prefix,
		)

		fmt.Printf(
			"%s%sHash : %s\n",
			prefix,
			indent,
			shortHash(n.Hash()),
		)

		childIndexes := make([]int, 0)

		for i, child := range n.childs {
			if child != nil {
				childIndexes = append(childIndexes, i)
			}
		}

		for i, index := range childIndexes {

			child := n.childs[index]

			lastChild := i == len(childIndexes)-1

			label := fmt.Sprintf(
				"%s%s[%X]",
				prefix,
				indent,
				index,
			)

			printNodeWithLabel(
				child,
				label,
				lastChild,
			)
		}
	}
}

func shortHash(hash [32]byte) string {
	return hex.EncodeToString(hash[:8])
}

func (t *Trie) insertPath(path []byte, value string) {
	val := []byte(value)

	if t.node == nil {
		t.node = newLeaf(path, val)
	} else {
		t.node = update_trie(t.node, path, val)
	}
}

func PrintTriedb(node Node) {
	printNode(node, "", true)
}

func printNodedb(
	node Node,
	prefix string,
	isLast bool,
) {

	if node == nil {
		fmt.Printf("%s<nil>\n", prefix)
		return
	}

	// Tree formatting
	connector := "├── "
	nextPrefix := prefix + "│   "

	if isLast {
		connector = "└── "
		nextPrefix = prefix + "    "
	}

	switch n := node.(type) {

	// --------------------------------------------------------
	// LEAF
	// --------------------------------------------------------

	case *Leaf:

		fmt.Printf(
			"%s%sLeaf\n",
			prefix,
			connector,
		)

		fmt.Printf(
			"%s    Path  : %v\n",
			nextPrefix,
			n.Path_unique,
		)

		fmt.Printf(
			"%s    Value : %s\n",
			nextPrefix,
			string(n.Value),
		)

		fmt.Printf(
			"%s    Hash  : %s\n",
			nextPrefix,
			shortHash(n.Hash()),
		)

	// --------------------------------------------------------
	// EXTENSION
	// --------------------------------------------------------

	case *Extension:

		fmt.Printf(
			"%s%sExtension\n",
			prefix,
			connector,
		)

		fmt.Printf(
			"%s    Shared Path : %v\n",
			nextPrefix,
			n.Path_shared,
		)

		fmt.Printf(
			"%s    Hash        : %s\n",
			nextPrefix,
			shortHash(n.Hash()),
		)

		printNode(
			n.next_branch,
			nextPrefix,
			true,
		)

	// --------------------------------------------------------
	// BRANCH
	// --------------------------------------------------------

	case *Branch:

		fmt.Printf(
			"%s%sBranch\n",
			prefix,
			connector,
		)

		fmt.Printf(
			"%s    Hash : %s\n",
			nextPrefix,
			shortHash(n.Hash()),
		)

		// Count children
		childIndexes := make([]int, 0)

		for i, child := range n.childs {
			if child != nil {
				childIndexes = append(childIndexes, i)
			}
		}

		for i, index := range childIndexes {

			child := n.childs[index]

			lastChild := i == len(childIndexes)-1

			label := fmt.Sprintf(
				"%s[%X]",
				nextPrefix,
				index,
			)

			printNodeWithLabel(
				child,
				label,
				lastChild,
			)
		}

	default:

		fmt.Printf(
			"%s%sUnknown node type %T\n",
			prefix,
			connector,
			node,
		)
	}
}

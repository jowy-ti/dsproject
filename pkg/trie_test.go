package pkg

import (
	"bytes"
	"testing"
)

func TestBytesToNibbles(t *testing.T) {
	var input [32]byte
	input[0] = 0xAB

	nibbles := bytesToNibbles(input)

	if nibbles[0] != 0x0A {
		t.Fatalf("expected high nibble 0x0A, got %x", nibbles[0])
	}

	if nibbles[1] != 0x0B {
		t.Fatalf("expected low nibble 0x0B, got %x", nibbles[1])
	}
}

func TestPathMismatch(t *testing.T) {
	a := []byte{1, 2, 3, 4}
	b := []byte{1, 2, 9, 4}

	index, newNibble, trieNibble := path_mismatch(a, b)

	if index != 2 {
		t.Fatalf("expected mismatch at 2, got %d", index)
	}

	if newNibble != 3 {
		t.Fatalf("expected new nibble 3, got %d", newNibble)
	}

	if trieNibble != 9 {
		t.Fatalf("expected trie nibble 9, got %d", trieNibble)
	}
}

func TestLeafHashDeterministic(t *testing.T) {
	leaf1 := newLeaf([]byte{1, 2}, []byte("hello"))
	leaf2 := newLeaf([]byte{1, 2}, []byte("hello"))

	if leaf1.Hash() != leaf2.Hash() {
		t.Fatal("equal leaves should have same hash")
	}
}

func TestLeafSerializeDeserialize(t *testing.T) {
	leaf := newLeaf([]byte{1, 2, 3}, []byte("hello"))

	serialized := leaf.Serialize()

	node := deserializeNode(serialized)

	decodedLeaf, ok := node.(*Leaf)
	if !ok {
		t.Fatal("deserialized node is not Leaf")
	}

	if !bytes.Equal(decodedLeaf.Value, leaf.Value) {
		t.Fatal("values differ")
	}

	if !bytes.Equal(decodedLeaf.Path_unique, leaf.Path_unique) {
		t.Fatal("paths differ")
	}
}

func TestBranchInsert(t *testing.T) {
	branch := newBranch()

	leaf := newLeaf([]byte{1, 2}, []byte("hello"))

	branch.insert(5, leaf)

	if branch.childs[5] == nil {
		t.Fatal("child not inserted")
	}

	if branch.Childs_hash[5] != leaf.Hash() {
		t.Fatal("child hash mismatch")
	}
}

func TestTrieInsertSingle(t *testing.T) {
	var trie Trie

	trie.insert("hello")

	if trie.node == nil {
		t.Fatal("root should not be nil")
	}

	_, ok := trie.node.(*Leaf)
	if !ok {
		t.Fatal("first insertion should create leaf")
	}
}

func TestTrieInsertCreatesBranch(t *testing.T) {
	var trie Trie

	trie.insert("hello")
	trie.insert("world")

	_, ok := trie.node.(*Branch)

	if !ok {
		_, extOk := trie.node.(*Extension)

		if !extOk {
			t.Fatal("expected root to become branch or extension")
		}
	}
}

func TestHashChangesAfterInsert(t *testing.T) {
	var trie Trie

	trie.insert("hello")

	hash1 := trie.node.Hash()

	trie.insert("world")

	hash2 := trie.node.Hash()

	if hash1 == hash2 {
		t.Fatal("root hash should change after insertion")
	}
}

func TestTrieDeterministic(t *testing.T) {
	var trie1 Trie
	var trie2 Trie

	values := []string{
		"hello",
		"world",
		"foo",
		"bar",
	}

	for _, v := range values {
		trie1.insert(v)
		trie2.insert(v)
	}

	if trie1.node.Hash() != trie2.node.Hash() {
		t.Fatal("tries should have identical root hashes")
	}
}

func verifyStoredNodesRecursive(
	t *testing.T,
	node Node,
	store map[[32]byte][]byte,
) {

	t.Helper()

	if node == nil {
		return
	}

	hash := node.Hash()

	serialized, ok := store[hash]

	if !ok {
		t.Fatalf(
			"missing node in store: %x",
			hash,
		)
	}

	deserialized := deserializeNode(serialized)

	// --------------------------------------------------------
	// HASH MUST MATCH
	// --------------------------------------------------------

	if deserialized.Hash() != hash {
		t.Fatalf(
			"deserialized hash mismatch for %T",
			node,
		)
	}

	// --------------------------------------------------------
	// TYPE MUST MATCH
	// --------------------------------------------------------

	switch original := node.(type) {

	case *Leaf:

		decodedLeaf, ok := deserialized.(*Leaf)

		if !ok {
			t.Fatalf(
				"expected Leaf after deserialize, got %T",
				deserialized,
			)
		}

		if !bytes.Equal(
			original.Path_unique,
			decodedLeaf.Path_unique,
		) {
			t.Fatal("leaf path mismatch")
		}

		if !bytes.Equal(
			original.Value,
			decodedLeaf.Value,
		) {
			t.Fatal("leaf value mismatch")
		}

	case *Branch:

		decodedBranch, ok := deserialized.(*Branch)

		if !ok {
			t.Fatalf(
				"expected Branch after deserialize, got %T",
				deserialized,
			)
		}

		// verify child hashes
		for i := 0; i < len(original.childs); i++ {

			originalChild := original.childs[i]
			decodedChildHash := decodedBranch.Childs_hash[i]

			if originalChild == nil {

				var empty [32]byte

				if decodedChildHash != empty {
					t.Fatalf(
						"expected empty child hash at %d",
						i,
					)
				}

				continue
			}

			if decodedChildHash != originalChild.Hash() {
				t.Fatalf(
					"child hash mismatch at index %d",
					i,
				)
			}

			verifyStoredNodesRecursive(
				t,
				originalChild,
				store,
			)
		}

	case *Extension:

		decodedExt, ok := deserialized.(*Extension)

		if !ok {
			t.Fatalf(
				"expected Extension after deserialize, got %T",
				deserialized,
			)
		}

		if !bytes.Equal(
			original.Path_shared,
			decodedExt.Path_shared,
		) {
			t.Fatal("extension path mismatch")
		}

		if decodedExt.Next_hash != original.next_branch.Hash() {
			t.Fatal("extension next hash mismatch")
		}

		verifyStoredNodesRecursive(
			t,
			original.next_branch,
			store,
		)
	}
}

func TestExtractKeyValues(t *testing.T) {

	var trie Trie

	/*
		Structure we want:

		[1 2 3 4] -> A
		[1 2 3 5] -> B
		[1 9 9 9] -> C

		This should generate multiple nodes:
		- leaves
		- branches
		- extension(s)
	*/

	trie.insertPath([]byte{1, 2, 3, 4}, "A")
	trie.insertPath([]byte{1, 2, 3, 5}, "B")
	trie.insertPath([]byte{1, 9, 9, 9}, "C")

	// --------------------------------------------------------
	// EXTRACT ALL NODES
	// --------------------------------------------------------

	nodes := make(map[[32]byte][]byte)

	extractKeyValues(trie.node, nodes)

	// --------------------------------------------------------
	// VERIFY ROOT EXISTS
	// --------------------------------------------------------

	rootHash := trie.node.Hash()

	rootSerialized, ok := nodes[rootHash]

	if !ok {
		t.Fatal("root node not found in extracted map")
	}

	// --------------------------------------------------------
	// VERIFY ROOT DESERIALIZES CORRECTLY
	// --------------------------------------------------------

	rootNode := deserializeNode(rootSerialized)

	if rootNode.Hash() != rootHash {
		t.Fatal("deserialized root hash mismatch")
	}

	// --------------------------------------------------------
	// RECURSIVELY VERIFY ALL STORED NODES
	// --------------------------------------------------------

	verifyStoredNodesRecursive(
		t,
		trie.node,
		nodes,
	)

	// --------------------------------------------------------
	// ENSURE MAP IS NOT EMPTY
	// --------------------------------------------------------

	if len(nodes) == 0 {
		t.Fatal("expected extracted nodes")
	}

	// --------------------------------------------------------
	// ENSURE ALL HASHES MATCH CONTENT
	// --------------------------------------------------------

	for hash, serialized := range nodes {

		node := deserializeNode(serialized)

		if node.Hash() != hash {
			t.Fatalf(
				"hash mismatch for node %T",
				node,
			)
		}
	}
}

func TestDebugTrie(t *testing.T) {

	var trie Trie

	// trie.insertPath([]byte{1, 2, 3, 4}, "A")
	// trie.insertPath([]byte{1, 2, 3, 5}, "B")
	// trie.insertPath([]byte{1, 7, 9, 4}, "C")
	// trie.insertPath([]byte{5, 5, 5, 6}, "D")

	trie.insertPath([]byte{1, 2, 3, 4}, "A")
	trie.insertPath([]byte{1, 2, 3, 5}, "B")
	trie.insertPath([]byte{1, 2, 8, 9}, "C")
	trie.insertPath([]byte{1, 7, 1, 1}, "D")
	trie.insertPath([]byte{1, 7, 1, 2}, "E")
	trie.insertPath([]byte{1, 7, 9, 9}, "F")
	trie.insertPath([]byte{5, 5, 5, 5}, "G")
	trie.insertPath([]byte{5, 5, 6, 6}, "H")
	trie.insertPath([]byte{9, 1, 1, 1}, "I")

	PrintTrie(trie.node)
}

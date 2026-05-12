package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Trie struct {
	node Node
}

type Node interface {
	Hash() [32]byte
	Serialize() []byte
}

func init() {
	gob.Register(&Extension{})
	gob.Register(&Branch{})
	gob.Register(&Leaf{})
}

func (t *Trie) insert(value string) {
	val := []byte(value)
	path := bytesToNibbles(sha256.Sum256(val))

	if t.node == nil {
		t.node = newLeaf(path[:], val)
	} else {
		t.node = update_trie(t.node, path[:], val)
	}
}

func update_trie(node Node, path []byte, value []byte) Node {
	switch n := node.(type) {
	case *Leaf:
		index, new_leaf_nible_mismatch, old_leaf_nible_mismatch := path_mismatch(path, n.Path_unique)

		if index == -1 {
			n.Value = value
			return n
		}

		branch := newBranch()

		// Add new leaf to the branch
		new_leaf := newLeaf(path[index+1:], value) // Create leaf with updated path
		branch.insert(new_leaf_nible_mismatch, new_leaf)

		// Add old leaf to the branch
		n.Path_unique = n.Path_unique[index+1:] // Update path
		branch.insert(old_leaf_nible_mismatch, n)

		if index == 0 {
			return branch
		}

		return newExtension(path[:index], branch.Hash(), branch)

	case *Branch:
		next_nible := path[0]
		leaf_path_left := path[1:]

		if n.childs[next_nible] == nil {
			n.childs[next_nible] = newLeaf(leaf_path_left, value)
		} else {
			n.childs[next_nible] = update_trie(n.childs[next_nible], leaf_path_left, value)
		}

		n.Childs_hash[next_nible] = n.childs[next_nible].Hash()
		return n

	case *Extension:
		index, new_leaf_nible_mismatch, extension_nible_mismatch := path_mismatch(path, n.Path_shared)
		len_path_shared := len(n.Path_shared)

		if index == -1 {
			n.next_branch = update_trie(n.next_branch, path[len_path_shared:], value)
			return n
		}

		branch := newBranch()

		// Add new leaf to the branch
		new_leaf := newLeaf(path[index+1:], value) // Update path and create leaf
		branch.insert(new_leaf_nible_mismatch, new_leaf)

		var first_part_path []byte = n.Path_shared[:index]

		// mismatch of a nible before the last one
		if index < len_path_shared-1 {
			n.Path_shared = n.Path_shared[index+1:] // Update extension shared_path
			branch.insert(extension_nible_mismatch, n)
		} else {
			branch.insert(extension_nible_mismatch, n.next_branch)
		}

		// mismatch of first nible (shared_path)
		if index == 0 {
			return branch
		}

		return newExtension(first_part_path, branch.Hash(), branch)
	}

	return node
}

// Helper

// returns the index where the mismatch is. Invariant (len(trie_path) <= len(new_path))
func path_mismatch(new_path []byte, trie_path []byte) (index int, new_path_mismatch_nible int, trie_path_mismatch_nible int) {
	length := len(trie_path)

	for i := 0; length > i; i++ {
		if new_path[i] != trie_path[i] {
			return i, int(new_path[i]), int(trie_path[i])
		}
	}

	return -1, -1, -1
}

// transform []byte to nibles doubling its length
func bytesToNibbles(data [32]byte) [64]byte {
	var nibbles [64]byte
	for i, b := range data {
		nibbles[i*2] = b >> 4     // High nibble
		nibbles[i*2+1] = b & 0x0F // Low nibble
	}
	return nibbles
}

// extractKeyValues retrieves a map with the hash and serialized node.
func extractKeyValues(node Node, m map[[32]byte][]byte) map[[32]byte][]byte {
	switch n := node.(type) {
	case *Leaf:
		m[n.Hash()] = n.Serialize()

	case *Branch:
		for i := 0; len(n.childs) > i; i++ {
			if n.childs[i] != nil {
				extractKeyValues(n.childs[i], m)
			}
		}
		m[n.Hash()] = n.Serialize()

	case *Extension:
		extractKeyValues(n.next_branch, m)
		m[n.Hash()] = n.Serialize()
	}

	return m
}

func deserializeNode(serializedNode []byte) Node {
	var node Node

	source := bytes.NewReader(serializedNode)
	decoder := gob.NewDecoder(source)
	err := decoder.Decode(&node)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return node
}

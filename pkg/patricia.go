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

type Extension struct {
	Path_shared []byte
	Next_hash   []byte
	Next_branch *Branch
}

type Branch struct {
	Childs_hash [16][32]byte
	Childs      [16]Node
}

type Leaf struct {
	Path_unique []byte
	Value       []byte
}

func newTrie() *Trie {
	return &Trie{}
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
		index := path_mismatch(path, n.Path_unique)

		if index == -1 {
			return n
		}

		branch := newBranch()

		// Add new leaf
		new_leaf_index := int(path[index])
		new_leaf_unique_path := path[index+1:] // Update path
		new_leaf := newLeaf(new_leaf_unique_path, value)
		branch.Childs[new_leaf_index] = new_leaf
		branch.Childs_hash[new_leaf_index] = new_leaf.Hash() // Add hash

		// Add old leaf
		n.Path_unique = n.Path_unique[index+1:] // Update path
		old_leaf_index := int(n.Path_unique[index])
		branch.Childs[old_leaf_index] = n
		branch.Childs_hash[old_leaf_index] = n.Hash() // Add hash

		if index == 0 {
			return branch
		}

		branch_hash := branch.Hash()
		return newExtension(path[:index], branch_hash[:], branch)

	case *Branch:
		index := path[0]
		leaf_unique_path := path[1:]

		if n.Childs[index] != nil {
			n.Childs[index] = update_trie(n.Childs[index], leaf_unique_path, value)
		} else {
			n.Childs[index] = newLeaf(leaf_unique_path, value)
		}

		n.Childs_hash[index] = n.Childs[index].Hash()
		return n

	case *Extension:
		index := path_mismatch(path, n.Path_shared)

	}
}

// Helper

// returns the index where the mismatch is
func path_mismatch(path1 []byte, path2 []byte) int {
	if len(path1) != len(path2) {
		log.Panicf("error: comparing paths with different lengths")
	}

	length := len(path1)

	for i := 0; length > i; i++ {
		if path1[i] != path2[i] {
			return i
		}
	}

	return -1
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

// Extension

func newExtension(shared_path []byte, next_hash []byte, next_branch *Branch) *Extension {
	return &Extension{
		Path_shared: shared_path,
		Next_hash:   next_hash,
		Next_branch: next_branch,
	}
}

func (e *Extension) Hash() [32]byte {
	headers := bytes.Join(
		[][]byte{
			e.Next_hash,
			e.Path_shared,
		},
		[]byte{},
	)
	return sha256.Sum256(headers)
}

func (e *Extension) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(e)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

// Branch

func newBranch() *Branch {
	return &Branch{}
}

func (b *Branch) Hash() [32]byte {
	allPieces := make([][]byte, 0, 17)
	for _, h := range b.Childs_hash {
		allPieces = append(allPieces, h[:])
	}
	allPieces = append(allPieces)
	data := bytes.Join(allPieces, []byte{})

	return sha256.Sum256(data)
}

func (b *Branch) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

// Leaf

func newLeaf(path_unique []byte, value []byte) *Leaf {
	return &Leaf{
		Path_unique: path_unique,
		Value:       value,
	}
}

func (l *Leaf) Hash() [32]byte {
	headers := bytes.Join(
		[][]byte{
			l.Path_unique,
			l.Value,
		},
		[]byte{},
	)
	return sha256.Sum256(headers)
}

func (l *Leaf) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(l)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

package pkg

import (
	"bytes"
	"crypto/sha256"
)

type Trie struct {
	Values [][]byte
	node   *Node
}

type Node interface {
	Hash() []byte
	Serialize() []byte
}

type Branch struct {
	Childs_hash [16][]byte
	Value       []byte
	Childs      [16]*Node
}

type Leaf struct {
	Path_rest []byte
	Value     []byte
}

type Extension struct {
	Shared_path []byte
	Next_hash   []byte
	Next_branch *Branch
}

func newTrie() *Trie {
	return &Trie{}
}

func (t *Trie) insert_value(value string) {
	t.Values = append(t.Values, []byte(value))
}

func (t *Trie) build() {
	for _, value := range t.Values {

	}
}

// Interface functions

func (b *Branch) Hash() [32]byte {
	allPieces := make([][]byte, 0, 17)
	for _, h := range b.Childs_hash {
		allPieces = append(allPieces, h)
	}
	allPieces = append(allPieces, b.Value)
	data := bytes.Join(allPieces, []byte{})

	return sha256.Sum256(data)
}

func (b *Branch) Serialize() []byte {

}

func (l *Leaf) Hash() [32]byte {
	headers := bytes.Join(
		[][]byte{
			l.Path_rest,
			l.Value,
		},
		[]byte{},
	)
	return sha256.Sum256(headers)
}

func (l *Leaf) Serialize() []byte {

}

func (e *Extension) Hash() [32]byte {
	headers := bytes.Join(
		[][]byte{
			e.Next_hash,
			e.Shared_path,
		},
		[]byte{},
	)
	return sha256.Sum256(headers)
}

func (e *Extension) Serialize() []byte {

}

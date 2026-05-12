package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Branch struct {
	Childs_hash [16][32]byte
	childs      [16]Node
}

func newBranch() *Branch {
	return &Branch{}
}

func (b *Branch) insert(index int, node Node) {
	b.childs[index] = node
	b.Childs_hash[index] = node.Hash()
}

func (b *Branch) Hash() [32]byte {
	allPieces := make([][]byte, 0, 17)
	for _, h := range b.Childs_hash {
		allPieces = append(allPieces, h[:])
	}
	data := bytes.Join(allPieces, []byte{})

	return sha256.Sum256(data)
}

func (b *Branch) Serialize() []byte {
	var result bytes.Buffer

	var node Node = b
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(&node)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

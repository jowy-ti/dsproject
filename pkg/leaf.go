package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Leaf struct {
	Path_unique []byte
	Value       []byte
}

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

	var node Node = l
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(&node)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

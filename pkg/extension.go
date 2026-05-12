package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Extension struct {
	Path_shared []byte
	Next_hash   [32]byte
	next_branch Node // always a branch
}

func newExtension(shared_path []byte, next_hash [32]byte, next_branch *Branch) *Extension {
	return &Extension{
		Path_shared: shared_path,
		Next_hash:   next_hash,
		next_branch: next_branch,
	}
}

func (e *Extension) Hash() [32]byte {
	headers := bytes.Join(
		[][]byte{
			e.Next_hash[:],
			e.Path_shared,
		},
		[]byte{},
	)
	return sha256.Sum256(headers)
}

func (e *Extension) Serialize() []byte {
	var result bytes.Buffer

	var node Node = e
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(&node)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

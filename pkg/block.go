package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"math"
	"time"
)

type Block struct {
	Hash          [32]byte
	PrevBlockHash [32]byte
	RootHash      [32]byte
	Nonce         uint64
	Difficulty    uint64
	Timestamp     int64
}

// newBlock creates a new Block using the provided data and the previous block hash
func newBlock(prevBlockHash [32]byte, rootHash [32]byte, difficulty uint64) *Block {
	block := &Block{
		PrevBlockHash: prevBlockHash,
		RootHash:      rootHash,
		Nonce:         math.MaxUint64,
		Difficulty:    difficulty,
		Timestamp:     time.Now().Unix(),
	}
	block.Hash = block.computeHash()
	return block
}

// newGenesisBlock creates the first block in the chain
func newGenesisBlock() *Block {
	return newBlock([32]byte{}, [32]byte{}, 0)
}

func (b *Block) setHash(hash [32]byte) {
	b.Hash = hash
}

func (b *Block) getHash() [32]byte {
	return b.Hash
}

func (b *Block) getPrevBlockHash() [32]byte {
	return b.PrevBlockHash
}

// serialize encodes the Block into a byte slice using GOB
func (b *Block) serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

// deserializeBlock converts bytes back into a Block struct
func deserializeBlock(serializedBlock []byte) *Block {
	var block Block

	source := bytes.NewReader(serializedBlock)
	decoder := gob.NewDecoder(source)
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return &block
}

// Iterates nonce
func (b *Block) nextNonce() {
	b.Nonce += 1
}

// setHash computes and assigns the SHA-256 hash of the block
func (b *Block) computeHash() [32]byte {
	headers := bytes.Join(
		[][]byte{
			b.PrevBlockHash[:],
			b.RootHash[:],
			intToHex(uint64(b.Timestamp)),
			intToHex(b.Difficulty),
			intToHex(b.Nonce),
		},
		[]byte{},
	)
	return sha256.Sum256(headers)
}

// Helper function to convert int64 to hex bytes
func intToHex(num uint64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, num)
	return buff
}

// validateTrie validates integrity through hashes. Should be passed the root hash of the Trie
func (b *Block) validateTrie(hash [32]byte, boltDB *boltStorage) bool {
	serializedNode := boltDB.dbGetTrieValue(hash)
	node := deserializeNode(serializedNode)
	// fmt.Printf("%x", hash)
	// println()

	if serializedNode == nil {
		return false
	}

	switch n := node.(type) {
	case *Branch:
		for i := 0; len(n.Childs_hash) > i; i++ {
			if n.Childs_hash[i] != [32]byte{} {
				if !b.validateTrie(n.Childs_hash[i], boltDB) {
					return false
				}
			}
		}

	case *Extension:
		if !b.validateTrie(n.Next_hash, boltDB) {
			return false
		}
	}
	return hash == node.Hash()
}

// extractKeyValues retrieves a map with the hash and serialized node.
func getHashesValuesStored(hash [32]byte, boltDB *boltStorage, m map[[32]byte]string) map[[32]byte]string {
	serializedNode := boltDB.dbGetTrieValue(hash)
	node := deserializeNode(serializedNode)

	if serializedNode == nil {
		return nil
	}

	switch n := node.(type) {
	case *Leaf:
		m[n.Hash()] = string(n.Value)

	case *Branch:
		for i := 0; len(n.Childs_hash) > i; i++ {
			if n.Childs_hash[i] != [32]byte{} {
				getHashesValuesStored(n.Childs_hash[i], boltDB, m)
			}
		}

	case *Extension:
		getHashesValuesStored(n.Next_hash, boltDB, m)
	}
	return m
}

func (b *Block) getTrieInfo(boltDB *boltStorage) map[[32]byte]string {
	return getHashesValuesStored(b.RootHash, boltDB, make(map[[32]byte]string))
}

// verifyValue verifies if the value is inside the Trie
func verifyValue(boltDB *boltStorage, hash [32]byte, value string) bool {
	serializedNode := boltDB.dbGetTrieValue(hash)
	node := deserializeNode(serializedNode)

	if serializedNode == nil {
		return false
	}

	switch n := node.(type) {
	case *Leaf:
		return value == string(n.Value)

	case *Branch:
		for i := 0; len(n.Childs_hash) > i; i++ {
			if n.Childs_hash[i] != [32]byte{} {
				if verifyValue(boltDB, n.Childs_hash[i], value) {
					return true
				}
			}
		}

	case *Extension:
		if verifyValue(boltDB, n.Next_hash, value) {
			return true
		}
	}
	return false
}

func (b *Block) verifyValueInTrie(boltDB *boltStorage, value string) bool {
	return verifyValue(boltDB, b.RootHash, value)
}

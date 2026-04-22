package blockchain

import (
	"bytes"
	"crypto/sha256"
	global "dsproject/internal"
	"encoding/binary"
	"encoding/json"
	"log"
	"math"
	"time"
)

type Block struct {
	Data          []byte
	Hash          [32]byte
	PrevBlockHash [32]byte
	Nonce         uint64
	Difficulty    uint64
	Timestamp     int64
}

// newBlock creates a new Block using the provided data and the previous block hash
func newBlock(data string, prevBlockHash [32]byte, difficulty uint64) *Block {
	block := &Block{
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Nonce:         math.MaxUint64,
		Difficulty:    difficulty,
		Timestamp:     time.Now().Unix(),
	}
	block.Hash = block.computeHash()
	return block
}

// newGenesisBlock creates the first block in the chain
func newGenesisBlock() *Block {
	return newBlock(global.GenesisBlockName, [32]byte{}, 0)
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

// serialize encodes the Block into a byte slice using JSON
func (b *Block) serialize() []byte {
	var result bytes.Buffer

	encoder := json.NewEncoder(&result)
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
	decoder := json.NewDecoder(source)
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
			b.Data,
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

// Old compute hash
// func (pow *ProofOfWork) computeHash(nonce uint64, difficulty uint64) [32]byte {
// 	timestamp := []byte(strconv.FormatInt(pow.block.timestamp, 10)) // converted to string to produce the same results with different CPU architectures
// 	nonceB := []byte(strconv.FormatUint(nonce, 10))
// 	difficultyB := []byte(strconv.FormatUint(difficulty, 10))
// 	headers := bytes.Join([][]byte{pow.block.prevBlockHash, pow.block.data, timestamp, nonceB, difficultyB}, []byte{})

// 	return sha256.Sum256(headers)
// }

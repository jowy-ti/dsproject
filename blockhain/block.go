package blockchain

import (
	"bytes"
	global "dsproject/internal"
	"encoding/json"
	"log"
	"time"
)

type Block struct {
	hash          []byte
	prevBlockHash []byte
	data          []byte
	nonce         uint64
	difficulty    uint64
	timestamp     int64
}

// newBlock creates a new Block using the provided data and the previous block hash
func newBlock(data string, prevBlockHash []byte, difficulty uint64) *Block {
	block := &Block{
		hash:          []byte{},
		prevBlockHash: prevBlockHash,
		data:          []byte(data),
		nonce:         0,
		difficulty:    difficulty,
		timestamp:     time.Now().Unix(),
	}

	return block
}

func (b *Block) setInfoAfterMining(nonce uint64, hash []byte) {
	b.hash = hash
	b.nonce = nonce
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

// newGenesisBlock creates the first block in the chain
func newGenesisBlock() *Block {
	return newBlock(global.GenesisBlockName, []byte{}, 0)
}

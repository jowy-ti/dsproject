package blockchain

import (
	"bytes"
	global "dsproject/internal"
	"encoding/json"
	"log"
	"time"
)

type Block struct {
	timestamp     int64
	data          []byte // This will eventually hold transactions
	prevBlockHash []byte
	hash          []byte
	nonce         uint
	difficulty    uint
}

// newBlock creates a new Block using the provided data and the previous block hash
func newBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		timestamp:     time.Now().Unix(),
		data:          []byte(data),
		prevBlockHash: prevBlockHash,
		hash:          []byte{},
	}

	return block
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
	return newBlock(global.GenesisBlockName, []byte{})
}

package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte // This will eventually hold transactions
	PrevBlockHash []byte
	Hash          []byte
}

const (
	genesisBlockName string = "Genesis Block"
)

func newBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}
	block.setHash()
	return block
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10)) // converted to string to produce the same results with different CPU architectures
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func (b *Block) serialize() []byte {
	var result bytes.Buffer

	encoder := json.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Bytes()
}

// DeserializeBlock converts bytes back into a Block struct
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

// NewGenesisBlock creates the first block in the chain
func newGenesisBlock() *Block {
	return newBlock(genesisBlockName, []byte{})
}

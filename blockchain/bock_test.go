package blockchain

import (
	"bytes"
	"testing"
)

func TestNewBlock(t *testing.T) {
	data := "Test Block Data"
	prevHash := []byte("prev_hash_example")
	block := newBlock(data, prevHash)

	// Verify Data was stored correctly
	if !bytes.Equal(block.Data, []byte(data)) {
		t.Errorf("Expected data %s, got %s", data, string(block.Data))
	}

	// Verify Hash was actually generated
	if len(block.Hash) == 0 {
		t.Error("Hash should not be empty")
	}

	// Verify the Link
	if !bytes.Equal(block.PrevBlockHash, prevHash) {
		t.Error("PrevBlockHash does not match the provided hash")
	}
}

func TestSerialization(t *testing.T) {
	block := newBlock("Serialize Me", []byte{1, 2, 3})

	// Convert to bytes
	serialized := block.serialize()

	// Convert back to struct
	deserialized := deserializeBlock(serialized)

	// Compare fields
	if !bytes.Equal(block.Hash, deserialized.Hash) {
		t.Error("Hashes do not match after deserialization")
	}

	if string(deserialized.Data) != "Serialize Me" {
		t.Error("Data was corrupted during serialization")
	}

	if block.Timestamp != deserialized.Timestamp {
		t.Error("Timestamp was lost during serialization")
	}

	// var prettyJSON bytes.Buffer
	// json.Indent(&prettyJSON, serialized, "", "  ")
	// fmt.Printf("\n--- Block JSON View ---\n%s\n", prettyJSON.String())
}

func TestGenesis(t *testing.T) {
	genesis := newGenesisBlock()

	if string(genesis.Data) != "Genesis Block" {
		t.Error("Genesis block data is incorrect")
	}

	if len(genesis.PrevBlockHash) != 0 {
		t.Error("Genesis block should not have a previous hash")
	}
}

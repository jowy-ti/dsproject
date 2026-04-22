package blockchain

import (
	global "dsproject/internal"
	"os"
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	if bc.tip == [32]byte{} {
		t.Error("Blockchain tip should not be empty after initialization")
	}

	// Verify genesis block exists by iterating
	it := bc.newIterator()
	if !it.validHash() {
		t.Error("Iterator should be valid for a new blockchain (Genesis block)")
	}
}

func TestAddBlock(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	oldTip := bc.tip

	bc.AddBlock("Send 1 BTC to Alice")

	if bc.tip == oldTip {
		t.Error("Blockchain tip should have updated after AddBlock")
	}
}

func TestSearchBlock(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	data := "Searchable Data"
	bc.AddBlock(data)

	targetHash := bc.tip

	found := bc.SearchBlock(targetHash)
	if found == nil {
		t.Fatal("Should have found the block by its hash")
	}

	if string(found.Data) != data {
		t.Errorf("Expected data %s, got %s", data, string(found.Data))
	}

	// Test searching for non-existent hash
	var fakeHash [32]byte
	fakeHash[0] = 0xFF
	if bc.SearchBlock(fakeHash) != nil {
		t.Error("Should not find a block with a fake hash")
	}
}

func TestIterator(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	bc.AddBlock("Block 1")
	bc.AddBlock("Block 2")

	count := 0
	it := bc.newIterator()

	for it.validHash() {
		block := it.getBlockAndAdvance()
		count++
		if block == nil {
			t.Fatal("Iterator returned nil block")
		}
	}

	// Genesis + Block 1 + Block 2 = 3 blocks
	if count != 3 {
		t.Errorf("Expected 3 blocks in chain, got %d", count)
	}
}

func TestValidateChain(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	bc.AddBlock("First valid block")
	bc.AddBlock("Second valid block")

	if !bc.ValidateChain() {
		t.Error("Chain should be valid")
	}
}

// newTestBlockchain creates a temporary blockchain instance for testing
func newTestBlockchain() (*Blockchain, func()) {
	bc := NewBlockchain(global.TestDBPath)

	cleanup := func() {
		bc.boltDB.db.Close()
		os.Remove(global.TestDBPath)
	}

	return bc, cleanup
}

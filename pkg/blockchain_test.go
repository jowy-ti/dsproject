package pkg

import (
	global "dsproject/pkg/config"
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

	bc.AddBlock()

	if bc.tip == oldTip {
		t.Error("Blockchain tip should have updated after AddBlock")
	}
}

func TestIterator(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	bc.AddBlock()
	bc.AddBlock()

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

	bc.AddBlock()
	bc.AddBlock()

	if !bc.ValidateChain() {
		t.Error("Chain should be valid")
	}
}

// newTestBlockchain creates a temporary blockchain instance for testing
func newTestBlockchain() (*Blockchain, func()) {
	bc := NewBlockchain(global.TestDBPath)

	cleanup := func() {
		bc.BoltDB.db.Close()
		os.Remove(global.TestDBPath)
	}

	return bc, cleanup
}

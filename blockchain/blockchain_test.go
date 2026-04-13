package blockchain

import (
	"os"
	"testing"
)

// TestBlockchainGenesis verifies that the blockchain initializes with a genesis block
func TestBlockchainGenesis(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	if bc.tip == nil {
		t.Fatal("expected genesis block, got nil tip")
	}

	block := bc.searchBlock(bc.tip)

	if block == nil {
		t.Fatal("genesis block not found")
	}
}

// TestAddBlock checks that adding a block updates the chain tip
func TestAddBlock(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	prevTip := bc.tip

	newBlock := newBlock("test block", prevTip)
	bc.addBlock(newBlock)

	if string(bc.tip) != string(newBlock.Hash) {
		t.Fatal("tip was not updated")
	}
}

// TestForEachBlock ensures all blocks are visited during iteration
func TestForEachBlock(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	block1 := newBlock("block1", bc.tip)
	bc.addBlock(block1)

	block2 := newBlock("block2", bc.tip)
	bc.addBlock(block2)

	count := 0

	bc.forEachBlock(func(b *Block) bool {
		count++
		return false
	})

	if count != 3 {
		t.Fatalf("expected 3 blocks, got %d", count)
	}
}

// TestForEachBlockBreak verifies that iteration stops when the callback returns true
func TestForEachBlockBreak(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	block1 := newBlock("block1", bc.tip)
	bc.addBlock(block1)

	count := 0

	bc.forEachBlock(func(b *Block) bool {
		count++
		return true // parar inmediatamente
	})

	if count != 1 {
		t.Fatalf("expected 1 iteration, got %d", count)
	}
}

// TestIteratorTraversal checks that the iterator traverses all blocks including genesis
func TestIteratorTraversal(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	// añadir bloques
	block1 := newBlock("block1", bc.tip)
	bc.addBlock(block1)

	block2 := newBlock("block2", bc.tip)
	bc.addBlock(block2)

	it := bc.newIterator()

	count := 0

	for it.validHash() {
		block := it.getCurrentBlock()
		count++
		it.nextBlock(block)
	}

	if count != 3 {
		t.Fatalf("expected 2 blocks + genesis block, got %d", count)
	}
}

// TestSearchBlockFound verifies that an existing block can be found by hash
func TestSearchBlockFound(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	block := newBlock("block1", bc.tip)
	bc.addBlock(block)

	found := bc.searchBlock(block.Hash)

	if found == nil {
		t.Fatal("block not found")
	}

	if string(found.Hash) != string(block.Hash) {
		t.Fatal("wrong block returned")
	}
}

// TestSearchBlockNotFound verifies that searching for a non-existent block returns nil
func TestSearchBlockNotFound(t *testing.T) {
	bc, cleanup := newTestBlockchain()
	defer cleanup()

	hash := []byte("nonexistent")

	found := bc.searchBlock(hash)

	if found != nil {
		t.Fatal("expected nil, got block")
	}
}

// newTestBlockchain creates a temporary blockchain instance for testing
func newTestBlockchain() (*Blockchain, func()) {
	bc := newBlockchain(testDBPath)

	cleanup := func() {
		bc.boltDB.db.Close()
		os.Remove(testDBPath)
	}

	return bc, cleanup
}

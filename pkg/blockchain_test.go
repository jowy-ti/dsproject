package pkg

import (
	"os"
	"testing"
)

// Helper function to setup and teardown a clean blockchain environment
func setupTestBlockchain(t *testing.T) (*Blockchain, func()) {
	dbPath := "test_blockchain.db"
	bc := NewBlockchain(dbPath)

	cleanup := func() {
		// Clean up the DB file if your boltStorage implementation creates one
		_ = os.Remove(dbPath)
	}

	return bc, cleanup
}

func TestNewBlockchain(t *testing.T) {
	bc, cleanup := setupTestBlockchain(t)
	defer cleanup()

	if bc == nil {
		t.Fatal("Expected blockchain to be initialized, got nil")
	}

	// Verify that the genesis block was created and the tip is updated
	if bc.tip == [32]byte{} {
		t.Error("Expected blockchain tip to be populated with Genesis block hash, got empty byte array")
	}

	// Verify database connection is attached
	if bc.BoltDB == nil {
		t.Error("Expected BoltDB reference to be set")
	}
}

func TestGetChainLengthAndForEachBlock(t *testing.T) {
	bc, cleanup := setupTestBlockchain(t)
	defer cleanup()

	// Initially, only the Genesis block should exist
	initialLength := bc.GetChainLength()
	if initialLength != 1 {
		t.Errorf("Expected initial chain length to be 1 (Genesis block), got %d", initialLength)
	}

	// Add some data and create a new block
	bc.InsertTrieValue("tx_data_1")
	bc.AddBlock()

	newLength := bc.GetChainLength()
	if newLength != 2 {
		t.Errorf("Expected chain length to be 2 after adding a block, got %d", newLength)
	}
}

func TestVerifyValueInBlock(t *testing.T) {
	bc, cleanup := setupTestBlockchain(t)
	defer cleanup()

	bc.InsertTrieValue("secret_value")
	bc.AddBlock() // This becomes pos 0 (the tip)

	// Since your iterator starts at the tip, position 0 is the block we just added.
	if !bc.VerifyValueInBlock(0, "secret_value") {
		t.Error("Expected to find 'secret_value' in block at position 0")
	}

	// Test checking a value that doesn't exist
	if bc.VerifyValueInBlock(0, "non_existent_value") {
		t.Error("Did not expect to find 'non_existent_value' in block at position 0")
	}

	// Test checking an out-of-bounds position
	if bc.VerifyValueInBlock(5, "secret_value") {
		t.Error("Did not expect to find a value at an invalid block position")
	}
}

func TestGetDataFromBlock(t *testing.T) {
	bc, cleanup := setupTestBlockchain(t)
	defer cleanup()

	bc.InsertTrieValue("important_transaction")
	bc.AddBlock()

	data, block := bc.GetDataFromBlock(0)
	if block == nil {
		t.Fatal("Expected to retrieve a block at position 0, got nil")
	}

	// Confirm that the returned block matches the current blockchain tip
	if block.getHash() != bc.tip {
		t.Errorf("Expected block hash %x to match chain tip %x", block.getHash(), bc.tip)
	}

	// Verify data map presence
	if data == nil {
		t.Error("Expected data map to be returned from block")
	}
}

func TestValidateChain(t *testing.T) {
	bc, cleanup := setupTestBlockchain(t)
	defer cleanup()

	// Add a block to make a multi-block chain
	bc.InsertTrieValue("block_data")
	bc.AddBlock()

	// The default stub structures should return true for standard validation
	if !bc.ValidateChain() {
		t.Error("Expected blockchain to validate successfully")
	}
}

func TestBlockchainIterator(t *testing.T) {
	bc, cleanup := setupTestBlockchain(t)
	defer cleanup()

	bc.InsertTrieValue("block_1")
	bc.AddBlock()
	bc.InsertTrieValue("block_2")
	bc.AddBlock()

	it := bc.newIterator()
	count := 0

	for it.validHash() {
		currentHashBeforeAdvance := it.currentHash
		block := it.getBlockAndAdvance()
		count++

		if block == nil {
			t.Fatalf("Iterator returned nil block at iteration %d", count)
		}

		// Ensure iterator advanced the hash to the previous block
		if it.currentHash == currentHashBeforeAdvance && it.validHash() {
			t.Errorf("Iterator failed to step backward; hash remained %x", currentHashBeforeAdvance)
		}
	}

	if count != 3 { // 2 added blocks + 1 Genesis block
		t.Errorf("Expected iterator to visit 3 blocks, visited %d", count)
	}
}

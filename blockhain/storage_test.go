package blockchain

import (
	"bytes"
	global "dsproject/internal"
	"os"
	"testing"
)

func TestNewBoltStorage(t *testing.T) {
	storage, cleanup := newTestDB()
	defer cleanup()

	if storage == nil || storage.db == nil {
		t.Fatal("Storage or DB connection is nil")
	}

	if !storage.dbExistsBucket() {
		t.Error("Bucket should have been created during initialization")
	}
}

func TestAddAndGetBlock(t *testing.T) {
	storage, cleanup := newTestDB()
	defer cleanup()

	// 1. Create mock data using [32]byte
	var mockHash [32]byte
	copy(mockHash[:], "block-hash-001") // Note: copy handles slice to array
	mockData := []byte("serialized-block-bytes")

	// 2. Test AddBlock
	storage.dbAddBlock(mockHash, mockData)

	// 3. Test Retrieval
	retrieved := storage.dbGetEncodedBlock(mockHash)
	if !bytes.Equal(retrieved, mockData) {
		t.Errorf("Expected data %s, got %s", string(mockData), string(retrieved))
	}
}

func TestLastHashPointer(t *testing.T) {
	storage, cleanup := newTestDB()
	defer cleanup()

	// Initially, last hash should be empty (32 zeros)
	if storage.dbGetLastHash() != ([32]byte{}) {
		t.Error("Initial last hash should be empty")
	}

	// Add a block
	var newHash [32]byte
	copy(newHash[:], "latest-block-hash")
	storage.dbAddBlock(newHash, []byte("some data"))

	// Verify the 'l' key updated
	lastHash := storage.dbGetLastHash()
	if lastHash != newHash {
		t.Errorf("Last hash pointer was not updated. Expected %x, got %x", newHash, lastHash)
	}
}

func TestMissingBlock(t *testing.T) {
	storage, cleanup := newTestDB()
	defer cleanup()

	var fakeHash [32]byte
	copy(fakeHash[:], "i-do-not-exist")

	retrieved := storage.dbGetEncodedBlock(fakeHash)
	if retrieved != nil {
		t.Error("Retrieving a non-existent block should return nil")
	}
}

// newTestDB used as a helper function to initialize, close and remove db
func newTestDB() (*boltStorage, func()) {
	storage := newBoltStorage(global.TestDBPath)

	cleanup := func() {
		storage.db.Close()
		os.Remove(global.TestDBPath)
	}

	return storage, cleanup
}

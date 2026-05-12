package pkg

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewBlock(t *testing.T) {
	data := "Test Block Data"
	prevHash := [32]byte{1, 2, 3}
	difficulty := uint64(12)

	b := newBlock(data, prevHash, [32]byte{}, difficulty)

	if string(b.Data) != data {
		t.Errorf("Expected data %s, got %s", data, b.Data)
	}
	if b.PrevBlockHash != prevHash {
		t.Errorf("Previous hash mismatch")
	}
	if b.Difficulty != difficulty {
		t.Errorf("Difficulty mismatch")
	}
	// Verify hash was actually computed and isn't just zeros
	if b.Hash == [32]byte{} {
		t.Error("Block hash should not be empty")
	}
}

func TestNewGenesisBlock(t *testing.T) {
	b := newGenesisBlock()

	if b.PrevBlockHash != [32]byte{} {
		t.Error("Genesis block should have an empty previous hash")
	}
	if b.Difficulty != 0 {
		t.Error("Genesis block should have 0 difficulty by default")
	}
}

func TestComputeHash(t *testing.T) {
	b := newBlock("Hash Test", [32]byte{}, [32]byte{}, 1)

	hash1 := b.computeHash()
	hash2 := b.computeHash()

	if hash1 != hash2 {
		t.Error("Hashes of the same data should be identical")
	}

	// Change the nonce and verify hash changes
	b.nextNonce()
	hash3 := b.computeHash()

	if hash1 == hash3 {
		t.Error("Hash should change when nonce changes")
	}
}

func TestSerialize(t *testing.T) {
	// NOTE: This test will fail unless you capitalize the Block fields!
	b := newBlock("Serialize Test", [32]byte{0xAA}, [32]byte{}, 1)

	serialized := b.serialize()
	if len(serialized) == 0 {
		t.Fatal("Serialized data should not be empty")
	}

	deserialized := deserializeBlock(serialized)

	// We use reflect.DeepEqual to compare the whole struct including slices
	if !reflect.DeepEqual(b, deserialized) {
		t.Errorf("Deserialized block does not match original")
	}
}

func TestIntToHex(t *testing.T) {
	val := uint64(1)
	expected := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	res := intToHex(val)

	if !bytes.Equal(res, expected) {
		t.Errorf("intToHex failed. Got %v, want %v", res, expected)
	}
}

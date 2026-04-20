package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork sets up a PoW for a specific block
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	// We shift 1 left by (256 - Difficulty) bits.
	// The higher the difficulty, the smaller the target number.
	target.Lsh(target, uint(256-config.Difficulty))

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		// Check if hashInt is less than the target
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

// Helper function to convert int64 to hex bytes
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(config.Difficulty)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

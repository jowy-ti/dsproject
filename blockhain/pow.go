package blockchain

import (
	"bytes"
	"crypto/sha256"
	global "dsproject/internal"
	"encoding/binary"
	"log"
	"math"
	"math/big"
	"strconv"
)

type ProofOfWork struct {
	block      *Block
	target     *big.Int
	difficulty uint
}

// NewProofOfWork sets up a PoW for a specific block
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	// We shift 1 left by (256 - Difficulty) bits.
	// The higher the difficulty, the smaller the target number.
	target.Lsh(target, uint(256-global.Difficulty))

	return &ProofOfWork{
		block:      b,
		target:     target,
		difficulty: global.Difficulty,
	}
}

func (pow *ProofOfWork) mine() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < math.MaxInt64 {
		hash = computeHash(nonce, global.Difficulty)
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

// setHash computes and assigns the SHA-256 hash of the block
func (pow *ProofOfWork) computeHash(nonce uint, difficulty uint) []byte {
	timestamp := []byte(strconv.FormatInt(pow.block.timestamp, 10)) // converted to string to produce the same results with different CPU architectures
	nonceB := []byte(strconv.FormatInt(int64(nonce), 10))
	difficultyB := []byte(strconv.FormatInt(int64(difficulty), 10))
	headers := bytes.Join([][]byte{pow.block.prevBlockHash, pow.block.data, timestamp, nonceB, difficultyB}, []byte{})
	hash := sha256.Sum256(headers)

	return hash[:]
}

// Helper function to convert int64 to hex bytes
func intToHex(num int64) []byte {
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

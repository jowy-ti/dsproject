package blockchain

import (
	"bytes"
	"crypto/sha256"
	global "dsproject/internal"
	"encoding/binary"
	"math"
	"math/big"
)

type ProofOfWork struct {
	block      *Block
	target     *big.Int
	difficulty uint64
}

// NewProofOfWork sets up a PoW for a specific block
func newProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-global.Difficulty))

	return &ProofOfWork{
		block:      b,
		target:     target,
		difficulty: global.Difficulty,
	}
}

func (pow *ProofOfWork) mine() (uint64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce uint64 = 0

	for nonce < math.MaxUint64 {
		hash = pow.computeHash(nonce)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

// setHash computes and assigns the SHA-256 hash of the block
func (pow *ProofOfWork) computeHash(nonce uint64) [32]byte {
	headers := pow.prepareData(nonce)
	return sha256.Sum256(headers)
}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.prevBlockHash,
			pow.block.data,
			IntToHex(uint64(pow.block.timestamp)),
			IntToHex(pow.difficulty),
			IntToHex(nonce),
		},
		[]byte{},
	)
	return data
}

// Helper function to convert int64 to hex bytes
func IntToHex(num uint64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, num)
	return buff
}

// Old compute hash
// func (pow *ProofOfWork) computeHash(nonce uint64, difficulty uint64) [32]byte {
// 	timestamp := []byte(strconv.FormatInt(pow.block.timestamp, 10)) // converted to string to produce the same results with different CPU architectures
// 	nonceB := []byte(strconv.FormatUint(nonce, 10))
// 	difficultyB := []byte(strconv.FormatUint(difficulty, 10))
// 	headers := bytes.Join([][]byte{pow.block.prevBlockHash, pow.block.data, timestamp, nonceB, difficultyB}, []byte{})

// 	return sha256.Sum256(headers)
// }

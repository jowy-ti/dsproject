package blockchain

import (
	global "dsproject/internal"
	"math/big"
)

type ProofOfWork struct {
	target     *big.Int
	difficulty uint64
}

// NewProofOfWork sets up a PoW for a specific block
func newProofOfWork() *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(global.HashLength-global.Difficulty))

	return &ProofOfWork{
		target:     target,
		difficulty: global.Difficulty,
	}
}

func (pow *ProofOfWork) mine(block *Block) {
	var hashInt big.Int
	var hash [32]byte = block.getHash()

	for hashInt.Cmp(pow.target) != -1 {
		block.nextNonce()
		hash = block.computeHash()
		hashInt.SetBytes(hash[:])
	}

	block.setHash(hash)
}

// fmt.Errorf("nextNonce. valid nonce out of limit uint64")

package dsproject

import "dsproject/blockchain"

const (
	dbPath string = "blockchain.db"
)

func main() {
	bc := blockchain.NewBlockchain(dbPath)

}

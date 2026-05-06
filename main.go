package main

import (
	blockchain "dsproject/pkg"
)

const (
	dbPath string = "blockchain.db"
)

func main() {
	bc := blockchain.NewBlockchain(dbPath)
	bc.AddBlock("first block")
	bc.AddBlock("second block")
	bc.AddBlock("third block")
	if !bc.ValidateChain() {
		println("inconsistent chain")
	}
	println("chain validated")
	println(bc.GetDataFromBlock(2))
	println(bc.GetDataFromBlock(1))
	println(bc.GetDataFromBlock(0))

	// value := "Holiiiii"
	// by := []byte(value)
	// print(fmt.Sprintf("string: %s, byte: % x", string(by), by))
}

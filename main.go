package main

import (
	"dsproject/pkg"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	dbPath string = "blockchain.db"
)

var blockPos int
var word string

var rootCmd = &cobra.Command{
	Use:   "bc",
	Short: "Blockchain CLI",
}

func main() {
	bc := pkg.NewBlockchain(dbPath)

	var addblockCmd = &cobra.Command{
		Use:   "addblock",
		Short: "Validate the chain creates the trie and add a new block",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			if !bc.ValidateChain() {
				fmt.Println("❌ Error: Inconsistent chain!")
			} else {
				fmt.Println("✅ chain validated.")
				for _, value := range args {
					bc.InsertTrieValue(value)
					fmt.Printf("-> Inserted: %s\n", value)
				}
				bc.AddBlock()
				fmt.Println("✅ block added.")
			}

		},
	}

	var verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Check if a value exists in a specific block",
		Run: func(cmd *cobra.Command, args []string) {

			chainLength := bc.GetChainLength()

			if blockPos >= chainLength || blockPos < 0 {
				fmt.Printf("🚫 There are no blocks in that position, try with a position between 0 and %d\n", chainLength-1)
				return
			}

			if bc.VerifyValueInBlock(blockPos, word) {
				fmt.Printf("✨ The value '%s' exists in block %d\n", word, blockPos)
			} else {
				fmt.Printf("🚫 The value '%s' was NOT found in block %d\n", word, blockPos)
			}
		},
	}

	verifyCmd.Flags().IntVarP(&blockPos, "block", "b", 0, "Block number to search")
	verifyCmd.Flags().StringVarP(&word, "word", "w", "", "Word to verify")
	verifyCmd.MarkFlagRequired("word")

	var debugCmd = &cobra.Command{
		Use:   "debugTrie",
		Short: "Print the trie of a specific block",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			blockPos, err := strconv.Atoi(args[0])

			if err != nil {
				print("Please introduce a valid number\n")
			}

			chainLength := bc.GetChainLength()

			if blockPos == chainLength-1 {
				print("The Genesis Block is empty\n")
				return
			}

			if blockPos >= chainLength || blockPos < 0 {
				fmt.Printf("🚫 There are no blocks in that position, try with a position between 0 and %d\n", chainLength-1)
				return
			}

			bc.PrintBlockTrie(blockPos)
		},
	}

	var dataCmd = &cobra.Command{
		Use:   "blockdata",
		Short: "Get the the content of a specific block",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			blockPos, err := strconv.Atoi(args[0])

			if err != nil {
				print("Please introduce a valid number\n")
			}

			chainLength := bc.GetChainLength()

			if blockPos >= chainLength || blockPos < 0 {
				fmt.Printf("🚫 There are no blocks in that position, try with a position between 0 and %d\n", chainLength-1)
				return
			}

			data, block := bc.GetDataFromBlock(blockPos)

			print("Block info:\n\n")
			fmt.Printf("Block hash: %x\n", block.Hash)
			fmt.Printf("Previous block hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Root hash of the trie: %x\n", block.RootHash)
			fmt.Printf("Difficulty: %d\n", block.Difficulty)
			fmt.Printf("Nonce: %d\n", block.Nonce)
			fmt.Printf("Timestamp: %d\n\n\n", block.Timestamp)

			print("Trie data:\n")
			for hash, value := range data {
				fmt.Printf("Hash: %x, Value: %s\n", hash, value)
			}
		},
	}

	var chainCmd = &cobra.Command{
		Use:   "chainpic",
		Short: "Print the current state of the chain",
		Run: func(cmd *cobra.Command, args []string) {

			chainLength := bc.GetChainLength()

			for i := range chainLength {
				fmt.Printf("%d ", i)
			}
			print("<- Genesis Block\n")
		},
	}

	rootCmd.AddCommand(addblockCmd, verifyCmd, debugCmd, dataCmd, chainCmd)

	// 3. Execute the Root Command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

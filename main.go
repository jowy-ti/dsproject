package main

import (
	"dsproject/pkg"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	dbPath string = "blockchain.db"
)

var blockNum int
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

			if bc.VerifyValueInBlock(blockNum, word) {
				fmt.Printf("✨ The value '%s' exists in block %d\n", word, blockNum)
			} else {
				fmt.Printf("🚫 The value '%s' was NOT found in block %d\n", word, blockNum)
			}
		},
	}
	verifyCmd.Flags().IntVarP(&blockNum, "block", "b", 0, "Block number to search")
	verifyCmd.Flags().StringVarP(&word, "word", "w", "", "Word to verify")
	verifyCmd.MarkFlagRequired("word")

	// var debugCmd = &cobra.Command{
	// 	Use:   "debug",
	// 	Short: "Print the trie of a specific block",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		PrintTrie(node Node)
	// 	},
	// }

	rootCmd.AddCommand(addblockCmd, verifyCmd)

	// 3. Execute the Root Command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/akhilesharora/go-merkle/pkg/client"
	"github.com/akhilesharora/go-merkle/pkg/merkle"
	"github.com/akhilesharora/go-merkle/pkg/server"
)

func main() {
	// Initialize the server with a set of files with dummy data
	files := []merkle.File{
		{Data: "File1 Contents"},
		{Data: "File2 Contents"},
		{Data: "File3 Contents"},
	}
	srv := server.NewServer(files)

	// Create a client
	cl := client.NewClient(srv)

	// Client computes and stores the Merkle root hash
	tree := merkle.BuildMerkleTree(files)
	rootHash := tree.Root.Hash
	err := os.WriteFile("rootHash.txt", rootHash[:], 0644)
	if err != nil {
		log.Fatalf("Failed to write root hash to file: %v", err)
	}

	// Client requests a file and its Merkle proof from the server
	fileIndex := 2 // Example: Requesting the second file
	file, proofBytes, directions, err := cl.RequestFile(fileIndex)
	if err != nil {
		log.Fatalf("Failed to request file from server: %v", err)
	}

	// Client verifies the file using the Merkle proof
	if !cl.VerifyFile(file, proofBytes, directions) {
		log.Fatal("Failed to verify file using the Merkle proof")
	} else {
		fmt.Printf("File %d verified successfully\n", fileIndex)
	}

	// Compare the reconstructed root hash with the stored root hash
	storedRootHashBytes, err := os.ReadFile("rootHash.txt")
	if err != nil {
		log.Fatalf("Failed to read stored root hash: %v", err)
	}

	if !compareHashes(rootHash, storedRootHashBytes) {
		log.Fatal("Reconstructed root hash does not match the stored root hash")
	} else {
		fmt.Println("Reconstructed root hash matches the stored root hash. File integrity verified.")
	}
}

// compareHashes compares two hash values
func compareHashes(hash1 [32]byte, hash2Bytes []byte) bool {
	if len(hash2Bytes) != len(hash1) {
		return false
	}
	var hash2 [32]byte
	copy(hash2[:], hash2Bytes)
	return bytes.Equal(hash1[:], hash2[:])
}

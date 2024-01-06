package server

import (
	"fmt"
	"github.com/akhilesharora/go-merkle/pkg/merkle"
)

type Server struct {
	merkleTree *merkle.MerkleTree
	fileData   []string // Slice to map indices to data
}

func NewServer(files []merkle.File) *Server {
	tree := merkle.BuildMerkleTree(files)

	var fileData []string
	for _, file := range files {
		fileData = append(fileData, file.Data)
	}

	return &Server{
		merkleTree: tree,
		fileData:   fileData,
	}
}

func (s *Server) GetFileData(fileIndex int) (string, error) {
	if fileIndex < 0 || fileIndex >= len(s.fileData) {
		return "", fmt.Errorf("file index out of range")
	}
	return s.fileData[fileIndex], nil
}

func (s *Server) GetMerkleRootHash() [32]byte {
	if s.merkleTree.Root == nil {
		return [32]byte{}
	}
	return s.merkleTree.Root.Hash
}

// GenerateMerkleProof generates a Merkle proof for the given file index.
// The proof is a slice of hashes with a sequence of directions
func (s *Server) GenerateMerkleProof(fileIndex int) ([][32]byte, []bool, error) {
	if fileIndex < 0 || fileIndex >= len(s.merkleTree.Leaves) {
		return nil, nil, fmt.Errorf("file index out of range")
	}

	var proof [][32]byte
	var directions []bool // false for left, true for right
	currentNode := s.merkleTree.Leaves[fileIndex]
	for currentNode != s.merkleTree.Root {
		sibling := currentNode.GetSibling()
		proof = append(proof, sibling.Hash)
		directions = append(directions, currentNode == currentNode.Parent.Left)
		currentNode = currentNode.Parent
	}
	return proof, directions, nil
}

package server

import (
	"fmt"

	"github.com/akhilesharora/go-merkle/pkg/merkle"
)

type ServerInterface interface {
	UploadFile(filename string, data []byte) uint
	GetFileData(fileIndex int) ([]byte, error)
	GetFileCount() int
	GetMerkleRootHash() [32]byte
	GenerateMerkleProof(fileIndex int) ([][32]byte, []bool, error)
}

type Server struct {
	MerkleTree *merkle.MerkleTree
	Files      [][]byte
}

func (s *Server) GetFileData(fileIndex int) ([]byte, error) {
	if fileIndex < 0 || fileIndex >= len(s.Files) {
		return nil, fmt.Errorf("file index out of range")
	}
	return s.Files[fileIndex], nil
}

func (s *Server) GetFileCount() int {
	return len(s.Files)
}

func NewServer() *Server {
	return &Server{
		MerkleTree: &merkle.MerkleTree{},
		Files:      [][]byte{},
	}
}

func (s *Server) UploadFile(filename string, data []byte) uint {
	s.Files = append(s.Files, data)
	s.updateMerkleTree()
	return uint(len(s.Files) - 1)
}

func (s *Server) updateMerkleTree() {
	var files []merkle.File
	for _, fileData := range s.Files {
		files = append(files, merkle.File{Data: string(fileData)})
	}
	s.MerkleTree = merkle.BuildMerkleTree(files)
}

func (s *Server) GetMerkleRootHash() [32]byte {
	if s.MerkleTree.Root == nil {
		return [32]byte{}
	}
	return s.MerkleTree.Root.Hash
}

func (s *Server) GenerateMerkleProof(fileIndex int) ([][32]byte, []bool, error) {
	if fileIndex < 0 || fileIndex >= len(s.MerkleTree.Leaves) {
		return nil, nil, fmt.Errorf("file index out of range")
	}
	return s.MerkleTree.GenerateProof(fileIndex)
}

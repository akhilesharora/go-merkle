package client

import (
	"bytes"

	"github.com/akhilesharora/go-merkle/pkg/merkle"
	"github.com/akhilesharora/go-merkle/pkg/server"
)

type Client struct {
	server *server.Server
}

func NewClient(server *server.Server) *Client {
	return &Client{server: server}
}

func (c *Client) RequestFile(fileIndex int) (merkle.File, []byte, []bool, error) {
	fileData, err := c.server.GetFileData(fileIndex)
	if err != nil {
		return merkle.File{}, nil, nil, err
	}

	proof, directions, err := c.server.GenerateMerkleProof(fileIndex)
	if err != nil {
		return merkle.File{}, nil, nil, err
	}

	// Convert proof to a byte slice
	var proofBytes []byte
	for _, hash := range proof {
		proofBytes = append(proofBytes, hash[:]...)
	}

	return merkle.File{Data: fileData}, proofBytes, directions, nil
}

func (c *Client) VerifyFile(file merkle.File, proofBytes []byte, directions []bool) bool {
	fileHash := merkle.Hash(file.Data)

	// Splitting the proof bytes back into hash slices
	var proof [][32]byte
	for i := 0; i < len(proofBytes); i += 32 {
		var hash [32]byte
		copy(hash[:], proofBytes[i:i+32])
		proof = append(proof, hash)
	}

	// Reconstruct the Merkle root using the proof and directions
	reconstructedRoot := reconstructMerkleRoot(fileHash, proof, directions)

	expectedRootHash := c.server.GetMerkleRootHash()
	return bytes.Equal(reconstructedRoot[:], expectedRootHash[:])
}

func reconstructMerkleRoot(fileHash [32]byte, proof [][32]byte, directions []bool) [32]byte {
	currentHash := fileHash
	for i, hash := range proof {
		if directions[i] {
			// If the current hash is the left child
			currentHash = merkle.HashPair(currentHash, hash)
		} else {
			// If the current hash is the right child
			currentHash = merkle.HashPair(hash, currentHash)
		}
	}
	return currentHash
}

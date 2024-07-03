package merkle

import (
	"crypto/sha256"
	"fmt"
)

// File represents a file with its data
type File struct {
	Data string
}

// Node represents a node in the Merkle tree
type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	Hash   [32]byte
}

// MerkleTree represents the entire Merkle tree
type MerkleTree struct {
	Root   *Node
	Leaves []*Node
}

// CreateHash computes the SHA-256 hash of input data
func CreateHash(data []byte) [32]byte {
	return sha256.Sum256(data)
}

// BuildMerkleTree constructs a Merkle tree from the given files
func BuildMerkleTree(files []File) *MerkleTree {
	if len(files) == 0 {
		return &MerkleTree{}
	}

	var leaves []*Node
	for _, file := range files {
		leaf := &Node{Hash: CreateHash([]byte(file.Data))}
		leaves = append(leaves, leaf)
	}

	root := buildTree(leaves)
	return &MerkleTree{Root: root, Leaves: leaves}
}

// buildTree builds the tree from the leaves up to the root
func buildTree(nodes []*Node) *Node {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var newLevel []*Node
	for i := 0; i < len(nodes); i += 2 {
		var left, right *Node = nodes[i], nil
		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			right = &Node{Hash: nodes[i].Hash} // Duplicate last node if number of nodes is odd
		}

		parent := &Node{
			Left:  left,
			Right: right,
			Hash:  CreateHash(append(left.Hash[:], right.Hash[:]...)),
		}
		left.Parent, right.Parent = parent, parent
		newLevel = append(newLevel, parent)
	}

	return buildTree(newLevel)
}

// GetSibling returns the sibling of the node. If the node is a left child, return the right child and vice versa.
func (n *Node) GetSibling() *Node {
	if n.Parent == nil {
		return nil // Root node has no sibling
	}
	if n == n.Parent.Left {
		return n.Parent.Right
	}
	return n.Parent.Left
}

// HashPair hashes two concatenated hashes
func HashPair(left, right []byte) [32]byte {
	combined := append(left, right...)
	return CreateHash(combined)
}

// GenerateProof generates a Merkle proof for the given file index.
func (t *MerkleTree) GenerateProof(index int) ([][32]byte, []bool, error) {
	if index < 0 || index >= len(t.Leaves) {
		return nil, nil, fmt.Errorf("index out of range")
	}

	var proof [][32]byte
	var directions []bool
	currentNode := t.Leaves[index]

	for currentNode.Parent != nil {
		sibling := currentNode.GetSibling()
		proof = append(proof, sibling.Hash)
		directions = append(directions, currentNode == currentNode.Parent.Left)
		currentNode = currentNode.Parent
	}

	// We shall reverse the proof and directions as we've built them from leaf to root
	for i, j := 0, len(proof)-1; i < j; i, j = i+1, j-1 {
		proof[i], proof[j] = proof[j], proof[i]
		directions[i], directions[j] = directions[j], directions[i]
	}

	return proof, directions, nil
}

package merkle

import (
	"crypto/sha256"
	"fmt"
)

type File struct {
	Data string
}

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	Hash   [32]byte
}

type MerkleTree struct {
	Root   *Node
	Leaves []*Node
}

// Hash computes the SHA-256 hash of input data
func Hash(data string) [32]byte {
	return sha256.Sum256([]byte(data))
}

// BuildMerkleTree constructs a Merkle tree from the given files
func BuildMerkleTree(files []File) *MerkleTree {
	if len(files) == 0 {
		return &MerkleTree{}
	}

	var leaves []*Node
	for _, file := range files {
		leaf := &Node{Hash: Hash(file.Data)}
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
			Hash:  Hash(string(left.Hash[:]) + string(right.Hash[:])),
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

func HashPair(hash1, hash2 [32]byte) [32]byte {
	combined := append(hash1[:], hash2[:]...)
	return Hash(string(combined))
}

// Update updates the data at the specified index and recalculates the necessary hashes.
func (t *MerkleTree) Update(index int, newData string) error {
	if index < 0 || index >= len(t.Leaves) {
		return fmt.Errorf("index out of range")
	}

	// Update the leaf node's data
	leafNode := t.Leaves[index]
	leafNode.Hash = Hash(newData)

	// Recalculate hashes up the tree
	currentNode := leafNode
	for currentNode.Parent != nil {
		parent := currentNode.Parent
		parent.Hash = Hash(string(parent.Left.Hash[:]) + string(parent.Right.Hash[:]))
		currentNode = parent
	}

	return nil
}

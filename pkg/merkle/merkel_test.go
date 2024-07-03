package merkle

import (
	"bytes"
	"testing"
)

func TestCreateHash(t *testing.T) {
	data := []byte("test data")
	hash := CreateHash(data)
	if hash == [32]byte{} {
		t.Error("CreateHash returned an empty hash")
	}
}

func TestBuildMerkleTree(t *testing.T) {
	files := []File{
		{Data: "file1"},
		{Data: "file2"},
		{Data: "file3"},
	}

	tree := BuildMerkleTree(files)

	if tree.Root == nil {
		t.Error("BuildMerkleTree returned a tree with nil root")
	}

	if len(tree.Leaves) != len(files) {
		t.Errorf("Expected %d leaves, got %d", len(files), len(tree.Leaves))
	}
}

func TestBuildMerkleTreeEmptyInput(t *testing.T) {
	tree := BuildMerkleTree([]File{})

	if tree.Root != nil {
		t.Error("BuildMerkleTree with empty input should return a tree with nil root")
	}

	if len(tree.Leaves) != 0 {
		t.Error("BuildMerkleTree with empty input should return a tree with no leaves")
	}
}

func TestGetSibling(t *testing.T) {
	files := []File{
		{Data: "file1"},
		{Data: "file2"},
	}

	tree := BuildMerkleTree(files)

	if tree.Leaves[0].GetSibling() != tree.Leaves[1] {
		t.Error("GetSibling failed for left leaf")
	}

	if tree.Leaves[1].GetSibling() != tree.Leaves[0] {
		t.Error("GetSibling failed for right leaf")
	}

	if tree.Root.GetSibling() != nil {
		t.Error("Root node should not have a sibling")
	}
}

func TestHashPair(t *testing.T) {
	left := []byte("left")
	right := []byte("right")
	hash := HashPair(left, right)

	if hash == [32]byte{} {
		t.Error("HashPair returned an empty hash")
	}

	// Test commutativity
	hashReverse := HashPair(right, left)
	if bytes.Equal(hash[:], hashReverse[:]) {
		t.Error("HashPair should not be commutative")
	}
}

func TestGenerateProof(t *testing.T) {
	files := []File{
		{Data: "file1"},
		{Data: "file2"},
		{Data: "file3"},
		{Data: "file4"},
	}

	tree := BuildMerkleTree(files)

	proof, directions, err := tree.GenerateProof(1)

	if err != nil {
		t.Errorf("GenerateProof returned an error: %v", err)
	}

	if len(proof) != 2 { // log2(4) = 2
		t.Errorf("Expected proof length 2, got %d", len(proof))
	}

	if len(directions) != len(proof) {
		t.Errorf("Directions length (%d) doesn't match proof length (%d)", len(directions), len(proof))
	}
}

func TestGenerateProofOutOfRange(t *testing.T) {
	files := []File{
		{Data: "file1"},
		{Data: "file2"},
	}

	tree := BuildMerkleTree(files)

	_, _, err := tree.GenerateProof(-1)
	if err == nil {
		t.Error("GenerateProof should return an error for negative index")
	}

	_, _, err = tree.GenerateProof(2)
	if err == nil {
		t.Error("GenerateProof should return an error for out of range index")
	}
}

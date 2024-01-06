package merkle

import (
	"reflect"
	"testing"
)

func TestBuildMerkleTree(t *testing.T) {
	files := []File{{Data: "File1"}, {Data: "File2"}}

	t.Run("ValidTreeConstruction", func(t *testing.T) {
		tree := BuildMerkleTree(files)
		if tree.Root == nil {
			t.Error("Merkle tree root should not be nil")
		}
		if len(tree.Leaves) != 2 {
			t.Errorf("Merkle tree should have 2 leaves, got %d", len(tree.Leaves))
		}
	})

	t.Run("SingleFileTree", func(t *testing.T) {
		singleFileTree := BuildMerkleTree([]File{{Data: "SingleFile"}})
		if len(singleFileTree.Leaves) != 1 {
			t.Errorf("Merkle tree should have 1 leaf for single file, got %d", len(singleFileTree.Leaves))
		}
	})

	t.Run("TreeWithOddNumberOfFiles", func(t *testing.T) {
		oddFiles := []File{{Data: "File1"}, {Data: "File2"}, {Data: "File3"}}
		tree := BuildMerkleTree(oddFiles)
		if len(tree.Leaves) != 3 {
			t.Errorf("Expected 3 leaves, got %d", len(tree.Leaves))
		}
	})

	t.Run("TreeStructureValidation", func(t *testing.T) {
		tree := BuildMerkleTree(files)
		for _, leaf := range tree.Leaves {
			current := leaf
			for current.Parent != nil {
				if current.Parent.Left != current && current.Parent.Right != current {
					t.Error("Tree structure invalid: parent-child relationship mismatch")
				}
				current = current.Parent
			}
			if current != tree.Root {
				t.Error("Tree structure invalid: didn't reach root")
			}
		}
	})

	t.Run("EmptyTree", func(t *testing.T) {
		emptyTree := BuildMerkleTree([]File{})
		if emptyTree.Root != nil {
			t.Error("Root of an empty Merkle tree should be nil")
		}
	})

	t.Run("LeafHashes", func(t *testing.T) {
		tree := BuildMerkleTree(files)
		for i, file := range files {
			if tree.Leaves[i].Hash != Hash(file.Data) {
				t.Errorf("Leaf hash does not match hash of file data for file %d", i)
			}
		}
	})

	t.Run("ParentHashes", func(t *testing.T) {
		tree := BuildMerkleTree(files)
		for _, leaf := range tree.Leaves {
			current := leaf
			for current.Parent != nil {
				expectedHash := Hash(string(current.Parent.Left.Hash[:]) + string(current.Parent.Right.Hash[:]))
				if current.Parent.Hash != expectedHash {
					t.Error("Parent hash does not match combined hash of its children")
				}
				current = current.Parent
			}
		}
	})

	t.Run("UpdateTree", func(t *testing.T) {
		tree := BuildMerkleTree(files)
		newData := "UpdatedFile"
		err := tree.Update(0, newData)
		if err != nil {
			t.Errorf("Update failed: %v", err)
		}
		if tree.Leaves[0].Hash != Hash(newData) {
			t.Error("Updated leaf hash does not match new data hash")
		}

		// Check if root hash has changed after update
		originalRootHash := tree.Root.Hash
		tree = BuildMerkleTree(files) // Rebuild tree with original files
		if reflect.DeepEqual(originalRootHash, tree.Root.Hash) {
			t.Error("Root hash did not change after leaf update")
		}
	})
}

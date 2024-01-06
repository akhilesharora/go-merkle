package server

import (
	"sync"
	"testing"

	"github.com/akhilesharora/go-merkle/pkg/merkle"
)

func TestServer(t *testing.T) {
	files := []merkle.File{{Data: "File1"}, {Data: "File2"}}
	srv := NewServer(files)

	t.Run("GetFileDataValid", func(t *testing.T) {
		data, err := srv.GetFileData(0)
		if err != nil || data != "File1" {
			t.Errorf("Expected File1, got %s, error: %v", data, err)
		}
	})

	t.Run("GetFileDataInvalidIndex", func(t *testing.T) {
		_, err := srv.GetFileData(-1)
		if err == nil {
			t.Error("Expected error for invalid file index, got none")
		}
	})

	t.Run("GenerateMerkleProofValid", func(t *testing.T) {
		proof, directions, err := srv.GenerateMerkleProof(0)
		if err != nil {
			t.Errorf("Failed to generate Merkle proof: %v", err)
		}
		if len(proof) == 0 || len(directions) == 0 {
			t.Errorf("Merkle proof or directions should not be empty")
		}
	})

	t.Run("GenerateMerkleProofInvalidIndex", func(t *testing.T) {
		_, _, err := srv.GenerateMerkleProof(-1)
		if err == nil {
			t.Error("Expected error for invalid file index in Merkle proof generation, got none")
		}
	})

	t.Run("GenerateMerkleProofForInvalidIndex", func(t *testing.T) {
		_, _, err := srv.GenerateMerkleProof(-1)
		if err == nil {
			t.Error("Expected error for invalid file index in Merkle proof generation, got none")
		}
	})

	t.Run("ServerUnderLoad", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, err := srv.GetFileData(i % len(files))
				if err != nil {
					t.Errorf("Error under load: %v", err)
				}
			}(i)
		}
		wg.Wait()
	})
}

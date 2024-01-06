package client

import (
	"sync"
	"testing"

	"github.com/akhilesharora/go-merkle/pkg/merkle"
	"github.com/akhilesharora/go-merkle/pkg/server"
)

func TestClient(t *testing.T) {
	files := []merkle.File{{Data: "File1"}, {Data: "File2"}}
	srv := server.NewServer(files)
	cl := NewClient(srv)

	t.Run("RequestFileValid", func(t *testing.T) {
		file, proofBytes, directions, err := cl.RequestFile(0)
		if err != nil || file.Data != "File1" {
			t.Errorf("Expected File1, got %s, error: %v", file.Data, err)
		}

		if !cl.VerifyFile(file, proofBytes, directions) {
			t.Errorf("Failed to verify file")
		}
	})

	t.Run("RequestFileInvalidIndex", func(t *testing.T) {
		_, _, _, err := cl.RequestFile(-1)
		if err == nil {
			t.Error("Expected error for invalid file index, got none")
		}
	})

	t.Run("RequestFileConcurrently", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, _, _, err := cl.RequestFile(i % len(files))
				if err != nil {
					t.Errorf("Concurrent request failed: %v", err)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("RequestFileFromEmptyServer", func(t *testing.T) {
		emptySrv := server.NewServer([]merkle.File{})
		emptyClient := NewClient(emptySrv)
		_, _, _, err := emptyClient.RequestFile(0)
		if err == nil {
			t.Error("Expected error when requesting file from empty server, got none")
		}
	})
}

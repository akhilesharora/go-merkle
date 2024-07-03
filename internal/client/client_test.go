package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/akhilesharora/go-merkle/pkg/merkle"
)

// Mock server handler for upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Mock server handler for download
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileData := []byte("mock file data")
	w.Write(fileData)
}

// Test for UploadFiles method
func TestUploadFiles(t *testing.T) {
	// Create a temporary file to upload
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte("temporary file content")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(uploadHandler))
	defer server.Close()

	client := NewClient(server.URL)
	rootHash, err := client.UploadFiles([]string{tmpfile.Name()})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rootHash == "" {
		t.Fatalf("expected a root hash, got empty string")
	}
}

func TestDownloadAndVerifyFile(t *testing.T) {
	mockFileData := []byte("mock file data")
	// Create a simple Merkle tree with one element for testing
	merkleTree := merkle.BuildMerkleTree([]merkle.File{{Data: string(mockFileData)}})
	rootHash := merkleTree.Root.Hash

	rootHashPath := "root_hash.txt"
	err := ioutil.WriteFile(rootHashPath, rootHash[:], 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(rootHashPath)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/download/") {
			w.Write(mockFileData)
		} else if strings.Contains(r.URL.Path, "/proof/") {
			proof, _, _ := merkleTree.GenerateProof(0)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"proof":      proof,
				"directions": proof,
			})
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	fileData, err := client.DownloadAndVerifyFile(0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Equal(fileData, mockFileData) {
		t.Fatalf("expected %s, got %s", mockFileData, fileData)
	}
}

// Test for VerifyProof method
func TestVerifyProof(t *testing.T) {
	client := NewClient("")

	fileHash := merkle.CreateHash([]byte("mock file data"))
	proofHash := merkle.CreateHash([]byte("mock proof"))
	proof := [][]byte{proofHash[:]} // Convert [32]byte to []byte
	directions := []bool{true}
	rootHash := merkle.CreateHash(append(fileHash[:], proofHash[:]...))

	if !client.VerifyProof([32]byte(fileHash[:]), proof, directions, rootHash[:]) {
		t.Fatalf("expected proof verification to succeed")
	}
}

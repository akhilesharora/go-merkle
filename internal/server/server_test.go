package server

import (
	"bytes"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.MerkleTree == nil {
		t.Error("NewServer: MerkleTree is nil")
	}
	if len(server.Files) != 0 {
		t.Errorf("NewServer: Expected 0 files, got %d", len(server.Files))
	}
}

func TestUploadFile(t *testing.T) {
	server := NewServer()
	testData := []byte("test data")
	server.UploadFile("test.txt", testData)

	if len(server.Files) != 1 {
		t.Errorf("UploadFile: Expected 1 file, got %d", len(server.Files))
	}
	if !bytes.Equal(server.Files[0], testData) {
		t.Error("UploadFile: Stored file data doesn't match uploaded data")
	}
	if server.MerkleTree.Root == nil {
		t.Error("UploadFile: MerkleTree root is nil after upload")
	}
}

func TestGetFileData(t *testing.T) {
	server := NewServer()
	testData := []byte("test data")
	server.UploadFile("test.txt", testData)

	data, err := server.GetFileData(0)
	if err != nil {
		t.Errorf("GetFileData: Unexpected error: %v", err)
	}
	if !bytes.Equal(data, testData) {
		t.Error("GetFileData: Retrieved data doesn't match uploaded data")
	}

	_, err = server.GetFileData(-1)
	if err == nil {
		t.Error("GetFileData: Expected error for negative index, got nil")
	}

	_, err = server.GetFileData(1)
	if err == nil {
		t.Error("GetFileData: Expected error for out of range index, got nil")
	}
}

func TestGetFileCount(t *testing.T) {
	server := NewServer()
	if server.GetFileCount() != 0 {
		t.Errorf("GetFileCount: Expected 0, got %d", server.GetFileCount())
	}

	server.UploadFile("test1.txt", []byte("test1"))
	server.UploadFile("test2.txt", []byte("test2"))
	if server.GetFileCount() != 2 {
		t.Errorf("GetFileCount: Expected 2, got %d", server.GetFileCount())
	}
}

func TestGetMerkleRootHash(t *testing.T) {
	server := NewServer()
	emptyHash := server.GetMerkleRootHash()
	if emptyHash != [32]byte{} {
		t.Error("GetMerkleRootHash: Expected empty hash for empty tree")
	}

	server.UploadFile("test.txt", []byte("test data"))
	hash := server.GetMerkleRootHash()
	if hash == [32]byte{} {
		t.Error("GetMerkleRootHash: Expected non-empty hash after upload")
	}
}

func TestGenerateMerkleProof(t *testing.T) {
	server := NewServer()
	server.UploadFile("test1.txt", []byte("test1"))
	server.UploadFile("test2.txt", []byte("test2"))

	proof, directions, err := server.GenerateMerkleProof(0)
	if err != nil {
		t.Errorf("GenerateMerkleProof: Unexpected error: %v", err)
	}
	if len(proof) == 0 {
		t.Error("GenerateMerkleProof: Expected non-empty proof")
	}
	if len(directions) != len(proof) {
		t.Error("GenerateMerkleProof: Mismatch between proof and directions length")
	}

	_, _, err = server.GenerateMerkleProof(-1)
	if err == nil {
		t.Error("GenerateMerkleProof: Expected error for negative index, got nil")
	}

	_, _, err = server.GenerateMerkleProof(2)
	if err == nil {
		t.Error("GenerateMerkleProof: Expected error for out of range index, got nil")
	}
}

func TestUpdateMerkleTree(t *testing.T) {
	server := NewServer()
	server.UploadFile("test1.txt", []byte("test1"))
	firstHash := server.GetMerkleRootHash()

	server.UploadFile("test2.txt", []byte("test2"))
	secondHash := server.GetMerkleRootHash()

	if firstHash == secondHash {
		t.Error("updateMerkleTree: Merkle root hash should change after adding a new file")
	}
}

func TestServerInterface(t *testing.T) {
	var _ ServerInterface = (*Server)(nil)
}

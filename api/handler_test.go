package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

// MockServer is a mock implementation of the server.ServerInterface
type MockServer struct {
	mock.Mock
}

func (m *MockServer) UploadFile(filename string, data []byte) uint {
	args := m.Called(filename, data)
	return args.Get(0).(uint)
}

func (m *MockServer) GetFileData(index int) ([]byte, error) {
	args := m.Called(index)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockServer) GetFileCount() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockServer) GetMerkleRootHash() [32]byte {
	args := m.Called()
	return args.Get(0).([32]byte)
}

func (m *MockServer) GenerateMerkleProof(index int) ([][32]byte, []bool, error) {
	args := m.Called(index)
	return args.Get(0).([][32]byte), args.Get(1).([]bool), args.Error(2)
}

func TestUploadHandler(t *testing.T) {
	mockServer := new(MockServer)
	handler := &Handlers{Server: mockServer}

	// Create a new file upload request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "testfile.txt")
	_, _ = part.Write([]byte("file content"))
	writer.Close()

	req, err := http.NewRequest("POST", "/upload?filename=testfile.txt", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	// Set up the mock expectation
	mockServer.On("UploadFile", "testfile.txt", []byte("file content")).Return(uint(0))

	handler.UploadHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	mockServer.AssertExpectations(t)
}

// TestDownloadHandler tests the DownloadHandler function
func TestDownloadHandler(t *testing.T) {
	mockServer := new(MockServer)
	handler := &Handlers{Server: mockServer}

	// Mock the server behavior
	mockServer.On("GetFileCount").Return(1)
	mockServer.On("GetFileData", 0).Return([]byte("file content"), nil)

	req, err := http.NewRequest("GET", "/download/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/download/{index}", handler.DownloadHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "file content"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	mockServer.AssertExpectations(t)
}

// TestProofHandler tests the ProofHandler function
func TestProofHandler(t *testing.T) {
	mockServer := new(MockServer)
	handler := &Handlers{Server: mockServer}

	// Mock the server behavior
	mockServer.On("GetFileCount").Return(1)
	mockProof := [][32]byte{{1, 2, 3}}
	mockServer.On("GenerateMerkleProof", 0).Return(mockProof, []bool{true}, nil)

	req, err := http.NewRequest("GET", "/proof/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/proof/{index}", handler.ProofHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	proof, ok := response["proof"].([]interface{})
	if !ok || len(proof) != 1 {
		t.Errorf("Unexpected proof in response: %v", response["proof"])
	}

	directions, ok := response["directions"].([]interface{})
	if !ok || len(directions) != 1 || directions[0] != true {
		t.Errorf("Unexpected directions in response: %v", response["directions"])
	}

	mockServer.AssertExpectations(t)
}

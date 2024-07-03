package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/akhilesharora/go-merkle/pkg/merkle"
)

type ClientInterface interface {
	UploadFiles(files []string) (string, error)
	DownloadAndVerifyFile(fileIndex int) ([]byte, error)
	VerifyProof(fileHash [32]byte, proof [][]byte, directions []bool, rootHash []byte) bool
}

type Client struct {
	serverURL string
}

func NewClient(serverURL string) *Client {
	return &Client{serverURL: serverURL}
}

func (c *Client) UploadFiles(files []string) (string, error) {
	var merkleFiles []merkle.File
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		merkleFiles = append(merkleFiles, merkle.File{Data: string(data)})

		if err := c.uploadFile(file, data); err != nil {
			return "", err
		}
	}

	merkleTree := merkle.BuildMerkleTree(merkleFiles)
	rootHash := merkleTree.Root.Hash

	if err := saveRootHash(rootHash); err != nil {
		return "", err
	}

	return hex.EncodeToString(rootHash[:]), nil
}

func (c *Client) uploadFile(file string, data []byte) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		return err
	}
	part.Write(data)
	writer.Close()

	req, err := http.NewRequest("POST", c.serverURL+"/upload?filename="+filepath.Base(file), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	log.Printf("Uploaded file: %s", file)
	return nil
}

func saveRootHash(rootHash [32]byte) error {
	err := os.WriteFile("root_hash.txt", rootHash[:], 0644)
	if err != nil {
		return err
	}
	log.Printf("Stored root hash: %s", hex.EncodeToString(rootHash[:]))
	return nil
}

func (c *Client) downloadFile(fileIndex int) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s/download/%d", c.serverURL, fileIndex))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Downloaded file index: %d", fileIndex)
	return fileData, nil
}

func (c *Client) getMerkleProof(fileIndex int) (*proofResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/proof/%d", c.serverURL, fileIndex))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var proofResponse proofResponse
	err = json.NewDecoder(resp.Body).Decode(&proofResponse)
	if err != nil {
		return nil, err
	}

	return &proofResponse, nil
}

type proofResponse struct {
	Proof      [][]byte `json:"proof"`
	Directions []bool   `json:"directions"`
}

func (c *Client) VerifyProof(fileHash [32]byte, proof [][]byte, directions []bool, rootHash []byte) bool {
	currentHash := fileHash
	log.Printf("Initial fileHash: %x", currentHash)

	for i, proofElement := range proof {
		log.Printf("Proof element %d: %x", i, proofElement)
		log.Printf("Direction %d: %v", i, directions[i])

		if directions[i] {
			currentHash = merkle.HashPair(currentHash[:], proofElement)
		} else {
			currentHash = merkle.HashPair(proofElement, currentHash[:])
		}
		log.Printf("Current hash after step %d: %x", i, currentHash)
	}

	log.Printf("Final calculated hash: %x", currentHash)
	log.Printf("Expected root hash: %x", rootHash)

	return bytes.Equal(currentHash[:], rootHash)
}

func (c *Client) DownloadAndVerifyFile(fileIndex int) ([]byte, error) {
	fileData, err := c.downloadFile(fileIndex)
	if err != nil {
		return nil, fmt.Errorf("download error: %w", err)
	}

	proofResponse, err := c.getMerkleProof(fileIndex)
	if err != nil {
		return nil, fmt.Errorf("proof retrieval error: %w", err)
	}

	fileHash := merkle.CreateHash(fileData)
	rootHash, err := os.ReadFile("root_hash.txt")
	if err != nil {
		return nil, fmt.Errorf("root hash read error: %w", err)
	}

	log.Printf("File hash: %x", fileHash)
	log.Printf("Root hash from file: %x", rootHash)

	if !c.VerifyProof(fileHash, proofResponse.Proof, proofResponse.Directions, rootHash) {
		return nil, fmt.Errorf("file verification failed: fileHash=%x, rootHash=%x", fileHash, rootHash)
	}

	log.Printf("Verified file index: %d", fileIndex)
	return fileData, nil
}

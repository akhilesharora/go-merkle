package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/akhilesharora/go-merkle/internal/server"
	"github.com/gorilla/mux"
)

type Handlers struct {
	Server server.ServerInterface
}

func (h *Handlers) UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := r.URL.Query().Get("filename")
	fileIndex := h.Server.UploadFile(filename, data)

	response := map[string]interface{}{
		"message":   "File uploaded successfully",
		"fileIndex": fileIndex,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, err := strconv.Atoi(vars["index"])
	if err != nil || index < 0 || index >= h.Server.GetFileCount() {
		http.Error(w, "Invalid file index", http.StatusBadRequest)
		return
	}

	fileData, err := h.Server.GetFileData(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(fileData)
}

func (h *Handlers) ProofHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, err := strconv.Atoi(vars["index"])
	if err != nil || index < 0 || index >= h.Server.GetFileCount() {
		http.Error(w, "Invalid file index", http.StatusBadRequest)
		return
	}

	proof, directions, err := h.Server.GenerateMerkleProof(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"proof":      proof,
		"directions": directions,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) ServeUI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

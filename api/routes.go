package api

import (
	"github.com/akhilesharora/go-merkle/internal/server"
	"github.com/gorilla/mux"
)

func SetupRoutes(s server.ServerInterface) *mux.Router {
	r := mux.NewRouter()
	h := &Handlers{Server: s}

	r.HandleFunc("/upload", CORSMiddleware(h.UploadHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/download/{index}", CORSMiddleware(h.DownloadHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/proof/{index}", CORSMiddleware(h.ProofHandler)).Methods("GET", "OPTIONS")
	r.PathPrefix("/").HandlerFunc(h.ServeUI)

	return r
}

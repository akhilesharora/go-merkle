package api

import (
	"github.com/akhilesharora/go-merkle/internal/server"
	"github.com/gorilla/mux"
)

func SetupRoutes(s server.ServerInterface) *mux.Router {
	r := mux.NewRouter()
	h := &Handlers{Server: s}

	r.HandleFunc("/upload", h.UploadHandler).Methods("POST")
	r.HandleFunc("/download/{index}", h.DownloadHandler).Methods("GET")
	r.HandleFunc("/proof/{index}", h.ProofHandler).Methods("GET")
	r.PathPrefix("/").HandlerFunc(h.ServeUI)

	return r
}

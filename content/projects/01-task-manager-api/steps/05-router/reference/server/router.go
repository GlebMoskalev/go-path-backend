package server

import (
	"encoding/json"
	"net/http"

	"taskmanager/handler"
)

func NewRouter(h *handler.TaskHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /tasks", h.Create)
	mux.HandleFunc("GET /tasks", h.List)
	mux.HandleFunc("GET /tasks/{id}", h.GetByID)
	mux.HandleFunc("PUT /tasks/{id}", h.Update)
	mux.HandleFunc("DELETE /tasks/{id}", h.Delete)

	return mux
}

func Run(addr string, h http.Handler) error {
	return http.ListenAndServe(addr, h)
}

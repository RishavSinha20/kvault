package api

import (
	"encoding/json"
	"kvault/store"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store *store.Store
}
type PutRequest struct {
	Value string `json:"value"`
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{
		store: s,
	}
}
func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	var req PutRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.store.Put(key, req.Value)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	val, ok := h.store.Get(key)

	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"value": val,
	})
}
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	ok := h.store.Delete(key)

	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

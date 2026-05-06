package api

import (
	"encoding/json"
	"kvault/store"
	"net/http"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{
		store: s,
	}
}

func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	h.store.Put(key, value)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

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
	key := r.URL.Query().Get("key")

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

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"kvault/store"
)

type Handler struct {
	store *store.Store
}

type PutRequest struct {
	Value string `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ValueResponse struct {
	Value string `json:"value"`
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{
		store: s,
		store: s,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: message,
	})
}

func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		writeError(w, http.StatusBadRequest, "missing key")
		return
	}

	var req PutRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.store.Put(key, req.Value)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to persist data")
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{
		Message: "key stored successfully",
	})
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		writeError(w, http.StatusBadRequest, "missing key")
		writeError(w, http.StatusBadRequest, "missing key")
		return
	}

	val, ok := h.store.Get(key)

	if !ok {
		writeError(w, http.StatusNotFound, "key not found")
		writeError(w, http.StatusNotFound, "key not found")
		return
	}

	writeJSON(w, http.StatusOK, ValueResponse{
		Value: val,
	writeJSON(w, http.StatusOK, ValueResponse{
		Value: val,
	})
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		writeError(w, http.StatusBadRequest, "missing key")
		writeError(w, http.StatusBadRequest, "missing key")
		return
	}

	ok, err := h.store.Delete(key)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to persist delete")
		return
	}

	if !ok {
		writeError(w, http.StatusNotFound, "key not found")
		writeError(w, http.StatusNotFound, "key not found")
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{
		Message: "key deleted successfully",
	})
}
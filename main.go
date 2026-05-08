package main

import (
	"log"
	"net/http"

	"kvault/api"
	"kvault/store"

	"github.com/gorilla/mux"
)

func main() {
	s := store.NewStore()

	h := api.NewHandler(s)

	r := mux.NewRouter()
	r.HandleFunc("store/{key}", h.PutHandler).Methods("PUT")
	r.HandleFunc("/store/{key}", h.GetHandler).Methods("GET")
	r.HandleFunc("/store/{key}", h.DeleteHandler).Methods("DELETE")

	log.Println("kvault server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

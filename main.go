package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"kvault/api"
	"kvault/store"
)

func main() {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		log.Fatal(err)
	}

	wal, err := store.NewWAL("data/wal.log")
	if err != nil {
		log.Fatal(err)
	}

	s, err := store.NewStore(wal)
	if err != nil {
		log.Fatal(err)
	}

	h := api.NewHandler(s)

	r := mux.NewRouter()

	r.HandleFunc("/store/{key}", h.PutHandler).Methods("PUT")
	r.HandleFunc("/store/{key}", h.GetHandler).Methods("GET")
	r.HandleFunc("/store/{key}", h.DeleteHandler).Methods("DELETE")

	log.Println("kvault server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}

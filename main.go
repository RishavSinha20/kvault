package main

import (
	"log"
	"net/http"

	"kvault/api"
	"kvault/store"
)

func main() {
	s := store.NewStore()

	h := api.NewHandler(s)

	http.HandleFunc("/put", h.PutHandler)
	http.HandleFunc("/get", h.GetHandler)
	http.HandleFunc("/delete", h.DeleteHandler)

	log.Println("kvault server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

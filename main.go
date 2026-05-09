package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"kvault/api"
	"kvault/replication"
	"kvault/store"
)

func main() {
	port := flag.String("port", "8080", "server port")
	role := flag.String("role", "leader", "leader or follower")
	flag.Parse()

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

	followers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	replicator := replication.NewReplicator(followers)

	isLeader := *role == "leader"

	h := api.NewHandler(s, replicator, isLeader)

	r := mux.NewRouter()

	r.HandleFunc("/store/{key}", h.PutHandler).Methods("PUT")
	r.HandleFunc("/store/{key}", h.GetHandler).Methods("GET")
	r.HandleFunc("/store/{key}", h.DeleteHandler).Methods("DELETE")

	r.HandleFunc("/internal/replicate", h.ReplicationHandler).Methods("POST")

	log.Printf("kvault node running on :%s as %s", *port, *role)

	log.Fatal(http.ListenAndServe(":"+*port, r))
}

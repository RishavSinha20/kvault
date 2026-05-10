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

	nodeAddress := "http://localhost:" + *port

	peers := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
	}

	var filteredPeers []string

	for _, peer := range peers {
		if peer != nodeAddress {
			filteredPeers = append(filteredPeers, peer)
		}
	}

	replicator := replication.NewReplicator(filteredPeers)

	initialLeader := *role == "leader"

	election := replication.NewElectionManager(
		nodeAddress,
		filteredPeers,
		initialLeader,
	)

	election.StartHeartbeatSender()
	election.StartElectionMonitor()

	h := api.NewHandler(s, replicator, election)

	r := mux.NewRouter()

	r.HandleFunc("/store/{key}", h.PutHandler).Methods("PUT")
	r.HandleFunc("/store/{key}", h.GetHandler).Methods("GET")
	r.HandleFunc("/store/{key}", h.DeleteHandler).Methods("DELETE")

	r.HandleFunc("/internal/replicate", h.ReplicationHandler).Methods("POST")

	r.HandleFunc("/internal/heartbeat", h.HeartbeatHandler).Methods("POST")

	log.Printf("kvault node running on :%s as %s", *port, *role)

	log.Fatal(http.ListenAndServe(":"+*port, r))
}
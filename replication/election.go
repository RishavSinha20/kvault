package replication

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Heartbeat struct {
	Leader string `json:"leader"`
}

type ElectionManager struct {
	mu sync.RWMutex

	nodeAddress string
	peers       []string

	isLeader bool

	lastHeartbeat time.Time
}

func NewElectionManager(
	nodeAddress string,
	peers []string,
	initialLeader bool,
) *ElectionManager {
	return &ElectionManager{
		nodeAddress:   nodeAddress,
		peers:         peers,
		isLeader:      initialLeader,
		lastHeartbeat: time.Now(),
	}
}
func (e *ElectionManager) IsLeader() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.isLeader
}
func (e *ElectionManager) BecomeLeader() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.isLeader = true

	log.Printf("Node %s became leader", e.nodeAddress)
}

func (e *ElectionManager) BecomeFollower() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.isLeader = false
}

func (e *ElectionManager) ReceiveHeartbeat() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.lastHeartbeat = time.Now()
}

func (e *ElectionManager) StartHeartbeatSender() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)

		for range ticker.C {
			if !e.IsLeader() {
				continue
			}

			heartbeat := Heartbeat{
				Leader: e.nodeAddress,
			}

			data, err := json.Marshal(heartbeat)
			if err != nil {
				continue
			}

			for _, peer := range e.peers {
				endpoint := peer + "/internal/heartbeat"

				resp, err := http.Post(
					endpoint,
					"application/json",
					bytes.NewBuffer(data),
				)

				if err != nil {
					continue
				}

				resp.Body.Close()
			}
		}
	}()
}

func (e *ElectionManager) StartElectionMonitor() {
	go func() {
		ticker := time.NewTicker(3 * time.Second)

		for range ticker.C {
			if e.IsLeader() {
				continue
			}

			e.mu.RLock()
			elapsed := time.Since(e.lastHeartbeat)
			e.mu.RUnlock()

			if elapsed > 5*time.Second {
				log.Println("Heartbeat timeout detected")
				e.BecomeLeader()
			}
		}
	}()
}

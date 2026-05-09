package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ReplicationRequest struct {
	Operation string `json:"operation"`
	Key       string `json:"key"`
	Value     string `json:"value,omitempty"`
}

type Replicator struct {
	followers []string
}

func NewReplicator(followers []string) *Replicator {
	return &Replicator{
		followers: followers,
	}
}

func (r *Replicator) Replicate(req ReplicationRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	for _, follower := range r.followers {
		endpoint := fmt.Sprintf("%s/internal/replicate", follower)

		resp, err := http.Post(
			endpoint,
			"application/json",
			bytes.NewBuffer(data),
		)

		if err != nil {
			return err
		}

		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("replication failed for follower: %s", follower)
		}
	}

	return nil
}

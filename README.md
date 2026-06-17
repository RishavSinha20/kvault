# KVault

> A fault-tolerant distributed key-value store built in Go featuring WAL-based persistence, crash recovery, leader-follower replication, and heartbeat-based leader election.

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Learning%20Project-blue)

---

## Overview

KVault is a distributed key-value store built from first principles to explore the core concepts behind modern databases and distributed systems.

The project evolved incrementally from a simple in-memory store into a fault-tolerant distributed system capable of:

- Storing data in memory with concurrent access support
- Persisting data using a Write-Ahead Log (WAL)
- Recovering state after crashes and restarts
- Replicating writes across multiple nodes
- Detecting node failures using heartbeats
- Electing new leaders automatically during failover

The goal of KVault is educational: to understand how systems such as Redis, etcd, Consul, and ZooKeeper work internally by implementing the fundamental building blocks from scratch.

---

# Architecture

## High-Level Cluster Architecture

```text
                    Client
                       │
                       ▼
                ┌────────────┐
                │   Leader   │
                └────────────┘
                       │
         ┌─────────────┴─────────────┐
         ▼                           ▼
 ┌──────────────┐           ┌──────────────┐
 │ Follower #1  │           │ Follower #2  │
 └──────────────┘           └──────────────┘
```

---

## Write Path

```text
Client Request
      │
      ▼
  Leader Node
      │
      ▼
 Write-Ahead Log
      │
      ▼
 In-Memory Store
      │
      ▼
 Replicate to Followers
```

---

## Recovery Path

```text
Server Restart
      │
      ▼
 Read WAL
      │
      ▼
 Replay Operations
      │
      ▼
 Rebuild In-Memory State
```

---

## Leader Election Flow

```text
Leader Sends Heartbeats
          │
          ▼
 Followers Track Heartbeats
          │
          ▼
 Heartbeat Timeout
          │
          ▼
 New Leader Promotion
```

---

# Features

## 1. Core Key-Value Store

- Thread-safe in-memory storage
- Concurrent read/write operations
- RWMutex-based synchronization
- RESTful HTTP API
- Structured JSON responses
- Clean modular architecture

---

## 2. Persistence Layer

- Write-Ahead Logging (WAL)
- Durable writes
- Crash recovery
- Automatic state reconstruction on restart

---

## 3. Replication Layer

- Leader-follower architecture
- Write replication
- Distributed state synchronization
- Read support on follower nodes

---

## 4. Leader Election

- Heartbeat-based failure detection
- Automatic leader promotion
- Basic fault tolerance
- Cluster failover support

---

# Tech Stack

| Category | Technology |
|-----------|------------|
| Language | Go |
| HTTP Router | Gorilla Mux |
| Storage | In-Memory Hash Map |
| Persistence | Write-Ahead Log |
| Replication | Leader-Follower |
| Concurrency | sync.RWMutex |
| Communication | HTTP |
| Deployment | Docker (Optional) |

---

# Project Structure

```text
kvault/
│
├── api/
│   └── handler.go
│
├── replication/
│   ├── replicator.go
│   └── election.go
│
├── store/
│   ├── store.go
│   └── wal.go
│
├── data/
│   └── wal.log
│
├── main.go
├── go.mod
└── go.sum
```

---

# Design Decisions

## Why Use RWMutex?

The store is read-heavy.

Using `sync.RWMutex` allows:

- Multiple readers simultaneously
- Single writer exclusivity
- Better performance than a standard mutex for read-heavy workloads

---

## Why Write-Ahead Logging?

Before modifying memory:

```text
Write Operation
      │
      ▼
 Persist to WAL
      │
      ▼
 Update Memory
```

This guarantees durability.

If the process crashes after writing to the WAL, the operation can be recovered during restart.

---

## Why Leader-Follower Replication?

Leader-follower replication provides:

- Simpler consistency model
- Centralized write coordination
- Easier failure handling

All writes go through the leader and are propagated to followers.

---

## Why Heartbeats?

Followers periodically receive heartbeat messages from the leader.

If heartbeats stop arriving:

```text
Leader Failure
      │
      ▼
 Heartbeat Timeout
      │
      ▼
 New Leader Election
```

This allows the cluster to recover from node failures automatically.

---

# Running Locally

## Clone Repository

```bash
git clone https://github.com/<your-username>/kvault.git
cd kvault
```

---

## Install Dependencies

```bash
go mod tidy
```

---

## Start Leader Node

```bash
go run main.go --port=8080 --role=leader
```

---

## Start Follower Node 1

```bash
go run main.go --port=8081 --role=follower
```

---

## Start Follower Node 2

```bash
go run main.go --port=8082 --role=follower
```

---

# API Documentation

## Store a Value

### Request

```http
PUT /store/name
Content-Type: application/json
```

```json
{
  "value": "rishav"
}
```

### Example

```bash
curl -X PUT http://localhost:8080/store/name \
-H "Content-Type: application/json" \
-d '{"value":"rishav"}'
```

### Response

```json
{
  "message": "key stored successfully"
}
```

---

## Retrieve a Value

### Request

```http
GET /store/name
```

### Example

```bash
curl http://localhost:8080/store/name
```

### Response

```json
{
  "value": "rishav"
}
```

---

## Delete a Value

### Request

```http
DELETE /store/name
```

### Example

```bash
curl -X DELETE http://localhost:8080/store/name
```

### Response

```json
{
  "message": "key deleted successfully"
}
```

---

# Persistence Demo

## Store Data

```bash
curl -X PUT http://localhost:8080/store/user \
-H "Content-Type: application/json" \
-d '{"value":"rishav"}'
```

---

## Stop Server

```bash
CTRL + C
```

---

## Restart Server

```bash
go run main.go --port=8080 --role=leader
```

---

## Verify Recovery

```bash
curl http://localhost:8080/store/user
```

Expected:

```json
{
  "value": "rishav"
}
```

The value is restored by replaying the WAL.

---

# Replication Demo

## Write to Leader

```bash
curl -X PUT http://localhost:8080/store/name \
-H "Content-Type: application/json" \
-d '{"value":"rishav"}'
```

---

## Read from Follower

```bash
curl http://localhost:8081/store/name
```

Expected:

```json
{
  "value": "rishav"
}
```

The follower has received the replicated update from the leader.

---

# Failover Demo

## Start Cluster

```bash
go run main.go --port=8080 --role=leader
go run main.go --port=8081 --role=follower
go run main.go --port=8082 --role=follower
```

---

## Kill Leader

Stop the leader process:

```bash
CTRL + C
```

---

## Observe Election

Followers detect heartbeat timeout.

Expected logs:

```text
Heartbeat timeout detected
Node http://localhost:8081 became leader
```

---

## Continue Writing

```bash
curl -X PUT http://localhost:8081/store/age \
-H "Content-Type: application/json" \
-d '{"value":"22"}'
```

The cluster continues serving writes after leader failure.

---

# Future Improvements

Potential enhancements include:

- Snapshotting
- WAL compaction
- Quorum-based replication
- Raft consensus
- Consistent hashing
- Data sharding
- Metrics and observability
- Kubernetes deployment
- Distributed tracing
- Authentication and authorization

---

# Key Learnings

Building KVault provided practical experience with:

- Concurrent programming in Go
- Storage engine fundamentals
- Write-Ahead Logging
- Crash recovery mechanisms
- Distributed replication
- Failure detection
- Leader election
- Fault tolerance
- Distributed systems design
- Backend architecture patterns

---

# Resume Bullet

```text
Designed and implemented a fault-tolerant distributed key-value store in Go featuring WAL-based durability, crash recovery, leader-follower replication, heartbeat-driven leader election, and automatic failover handling.
```

---

# License

This project is licensed under the MIT License.

---

## Author

**Rishav Sinha**

Backend Engineer | Golang Developer | Distributed Systems Enthusiast

package store

import "sync"

type Store struct {
	data map[string]string
	mu   sync.RWMutex
	wal  *WAL
}

func NewStore(wal *WAL) (*Store, error) {
	s := &Store{
		data: make(map[string]string),
		wal:  wal,
	}

	err := s.recover()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Put(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.wal.Write(WALRecord{
		Operation: "PUT",
		Key:       key,
		Value:     value,
	})

	if err != nil {
		return err
	}

	s.data[key] = value

	return nil
}

func (s *Store) ReplicatedPut(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

func (s *Store) Delete(key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]

	if !ok {
		return false, nil
	}

	err := s.wal.Write(WALRecord{
		Operation: "DELETE",
		Key:       key,
	})

	if err != nil {
		return false, err
	}
	delete(s.data, key)

	return true, nil
}
func (s *Store) ReplicatedDelete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

func (s *Store) recover() error {
	records, err := s.wal.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		switch record.Operation {
		case "PUT":
			s.data[record.Key] = record.Value

		case "DELETE":
			delete(s.data, record.Key)
		}
	}

	return nil
}

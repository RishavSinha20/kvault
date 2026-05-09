package store

import (
	"bufio"
	"encoding/json"
	"os"
)

type WALRecord struct {
	Operation string `json:"operation"`
	Key       string `json:"key"`
	Value     string `json:"value,omitempty"`
}

type WAL struct {
	file *os.File
}

func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return &WAL{file: file}, nil
}

func (w *WAL) Write(record WALRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = w.file.Write(data)
	if err != nil {
		return err
	}

	return w.file.Sync()
}

func (w *WAL) ReadAll() ([]WALRecord, error) {
	var records []WALRecord

	_, err := w.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(w.file)

	for scanner.Scan() {
		var record WALRecord

		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, scanner.Err()
}

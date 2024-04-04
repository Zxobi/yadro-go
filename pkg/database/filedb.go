package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type FileDatabase struct {
	filename string
	records  RecordMap
}

func NewFileDatabase(fName string) (*FileDatabase, error) {
	fdb := &FileDatabase{filename: fName}
	if err := fdb.init(); err != nil {
		return nil, err
	}

	records, err := fdb.readFromFile()
	if err != nil {
		return nil, err
	}

	fdb.records = records

	fmt.Println("Database initialized:", len(records), "records loaded")
	return fdb, nil
}

func (db *FileDatabase) init() error {
	if _, err := os.Stat(db.filename); !errors.Is(err, os.ErrNotExist) {
		return err
	}

	file, err := os.Create(db.filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func (db *FileDatabase) Read() RecordMap {
	return makeCopy(db.records)
}

func (db *FileDatabase) Write(records RecordMap) error {
	cpy := makeCopy(records)

	data, err := json.Marshal(cpy)
	if err != nil {
		return err
	}

	if err = os.WriteFile(db.filename, data, 0644); err != nil {
		return err
	}

	db.records = cpy
	return nil
}

func (db *FileDatabase) readFromFile() (RecordMap, error) {
	data, err := os.ReadFile(db.filename)
	if err != nil {
		return nil, err
	}

	records := make(RecordMap)
	if len(data) == 0 {
		return records, nil
	}

	if err = json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}

func makeCopy(original RecordMap) RecordMap {
	records := make(RecordMap, len(original))
	for k, v := range original {
		records[k] = v
	}

	return records
}

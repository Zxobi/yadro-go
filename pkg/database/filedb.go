package database

import (
	"encoding/json"
	"errors"
	"os"
)

type FileDatabase struct {
	filename string
}

type IDatabase interface {
	Init() error
	Read() (RecordMap, error)
	Write(hash RecordMap) error
}

func NewFileDatabase(fName string) IDatabase {
	return &FileDatabase{filename: fName}
}

func (db *FileDatabase) Init() error {
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

func (db *FileDatabase) Read() (RecordMap, error) {
	data, err := os.ReadFile(db.filename)
	if err != nil {
		return nil, err
	}

	hash := make(RecordMap)
	if len(data) == 0 {
		return hash, nil
	}

	if err = json.Unmarshal(data, &hash); err != nil {
		return nil, err
	}

	return hash, nil
}

func (db *FileDatabase) Write(hash RecordMap) error {
	data, err := json.Marshal(hash)
	if err != nil {
		return err
	}

	return os.WriteFile(db.filename, data, 0644)
}

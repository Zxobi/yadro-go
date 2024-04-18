package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
)

type FileDatabase struct {
	log           *slog.Logger
	dbFilename    string
	indexFilename string
	records       RecordMap
	index         IndexMap
}

func NewFileDatabase(log *slog.Logger, dbFilename string, indexFilename string) (*FileDatabase, error) {
	fdb := &FileDatabase{
		log:           log,
		dbFilename:    dbFilename,
		indexFilename: indexFilename,
	}
	if err := fdb.init(); err != nil {
		return nil, err
	}

	log.Info(
		fmt.Sprintf("database initialized: %d records loaded, index size %d", len(fdb.records), len(fdb.index)))
	return fdb, nil
}

func (db *FileDatabase) init() error {
	if err := errors.Join(
		checkCreateFile(db.dbFilename),
		checkCreateFile(db.indexFilename),
	); err != nil {
		return err
	}

	records, err := db.readRecordsFromFile()
	if err != nil {
		return err
	}
	index, err := db.readIndexFromFile()
	if err != nil {
		return err
	}

	db.records = records
	db.index = index
	return nil
}

func checkCreateFile(path string) error {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_ = file.Close()
	return nil
}

func (db *FileDatabase) Records() RecordMap {
	return maps.Clone(db.records)
}

func (db *FileDatabase) Index() IndexMap {
	return maps.Clone(db.index)
}

func (db *FileDatabase) Save(records RecordMap) error {
	db.log.Debug("saving records")

	cpy := maps.Clone(records)

	recordJson, err := json.Marshal(cpy)
	if err != nil {
		return err
	}
	if err = os.WriteFile(db.dbFilename, recordJson, 0644); err != nil {
		return err
	}

	db.records = cpy

	db.log.Debug("records save complete")

	index := db.buildIndex()
	indexJson, err := json.Marshal(index)
	if err != nil {
		return err
	}
	if err = os.WriteFile(db.indexFilename, indexJson, 0644); err != nil {
		return err
	}

	db.index = index

	db.log.Debug("index save complete")

	return nil
}

func (db *FileDatabase) buildIndex() IndexMap {
	db.log.Debug("building index")

	index := make(IndexMap)
	for num, record := range db.records {
		for _, keyword := range record.Keywords {
			nums, ok := index[keyword]
			if ok {
				nums = append(nums, num)
			} else {
				nums = []int{num}
			}

			index[keyword] = nums
		}
	}

	db.log.Debug("index build complete")

	return index
}

func (db *FileDatabase) readRecordsFromFile() (RecordMap, error) {
	data, err := os.ReadFile(db.dbFilename)
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

func (db *FileDatabase) readIndexFromFile() (IndexMap, error) {
	data, err := os.ReadFile(db.indexFilename)
	if err != nil {
		return nil, err
	}

	index := make(IndexMap)
	if len(data) == 0 {
		return index, nil
	}

	if err = json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	return index, nil
}

package service

import (
	"errors"
	"fmt"
	"yadro-go/pkg/database"
	"yadro-go/pkg/words"
	"yadro-go/pkg/xkcd"
)

const BatchSize = 100

type ComicsService struct {
	db    IDatabase
	c     IClient
	limit int
}

type IDatabase interface {
	Read() database.RecordMap
	Write(records database.RecordMap) error
}

type IClient interface {
	GetById(id int) (xkcd.Comic, error)
}

func NewComicsService(c IClient, db IDatabase, limit int) *ComicsService {
	return &ComicsService{
		db:    db,
		c:     c,
		limit: limit,
	}
}

func (c *ComicsService) Fetch() error {
	records := c.db.Read()

	for count, i := 1, len(records)+1; i <= c.limit; count, i = count+1, i+1 {
		// comic with id 404 is always not found, skipping
		if i == 404 {
			records[i] = database.NewRecord("", []string{})
			continue
		}

		comic, err := c.c.GetById(i)
		if err != nil {
			if errors.Is(err, xkcd.NotFound) {
				break
			}
			return err
		}

		records[i] = makeRecord(comic)
		if count%BatchSize == 0 {
			fmt.Println("Fetched", count, "new records")
			if err = c.db.Write(records); err != nil {
				return err
			}
		}
	}

	fmt.Println("Fetch finished:", len(records), "records total")
	return c.db.Write(records)
}

func makeRecord(comic xkcd.Comic) database.Record {
	stemmed := words.Stem(comic.Title + " " + comic.Alt + " " + comic.Transcript)
	return database.NewRecord(comic.Img, stemmed)
}

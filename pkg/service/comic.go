package service

import (
	"errors"
	"fmt"
	"yadro-go/pkg/database"
	"yadro-go/pkg/words"
	"yadro-go/pkg/xkcd"
)

type ComicsService struct {
	db    database.IDatabase
	c     xkcd.IClient
	limit int
}

type IComicsService interface {
	Fetch() error
}

func NewComicsService(c xkcd.IClient, db database.IDatabase, limit int) IComicsService {
	return &ComicsService{
		db:    db,
		c:     c,
		limit: limit,
	}
}

func (c *ComicsService) Fetch() error {
	data, err := c.db.Read()
	if err != nil {
		return err
	}

	for i := 1; c.limit <= 0 || i <= c.limit; i++ {
		// comic with id 404 is always not found, skipping
		if i == 404 {
			continue
		}

		if _, ok := data[i]; ok {
			continue
		}

		comic, err := c.c.GetById(i)
		if err != nil {
			if errors.Is(err, xkcd.NotFound) {
				break
			}
			return err
		}

		data[i] = makeEntity(comic)
	}

	return c.db.Write(data)
}

func makeEntity(comic xkcd.Comic) database.Entity {
	stemmed := words.Stem(fmt.Sprintf("%s %s %s", comic.Title, comic.Alt, comic.Transcript))
	return database.NewEntity(comic.Img, stemmed)
}

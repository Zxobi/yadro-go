package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"yadro-go/pkg/database"
	"yadro-go/pkg/words"
	"yadro-go/pkg/xkcd"
)

type ComicsService struct {
	db       IDatabase
	c        IClient
	limit    int
	parallel int
}

type IDatabase interface {
	Read() database.RecordMap
	Write(records database.RecordMap) error
}

type IClient interface {
	GetById(id int) (*xkcd.Comic, error)
}

func NewComicsService(c IClient, db IDatabase, limit int, parallel int) *ComicsService {
	return &ComicsService{
		db:       db,
		c:        c,
		limit:    limit,
		parallel: parallel,
	}
}

func (c *ComicsService) Fetch(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	ids := make(chan int, c.parallel)
	comics := make(chan *xkcd.Comic, c.parallel)
	errs := make(chan error, c.parallel)

	records := c.db.Read()
	fmt.Println("start fetching with initial db size", len(records))

	fetchId := 1
	pushId := func() bool {
		for ; fetchId <= c.limit; fetchId++ {
			_, ok := records[fetchId]
			if !ok {
				ids <- fetchId
				fetchId++
				return true
			}
		}

		return false
	}

	for i := 0; i < c.parallel; i++ {
		if !pushId() {
			if i == 0 {
				fmt.Println("fetch finished: nothing to fetch")
				cancel()
				return nil
			}
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.fetchJob(ctx, ids, comics, errs)
		}()
	}

	var err error
	newCount := 0
	loop := true
	for loop {
		select {
		case comic := <-comics:
			newCount++
			records[comic.Num] = *makeRecord(comic)
			if !pushId() {
				loop = false
			}
		case workerErr := <-errs:
			if !errors.Is(workerErr, xkcd.NotFound) {
				fmt.Println("finishing fetching due to worker error")
				err = workerErr
			}
			loop = false
		case <-ctx.Done():
			fmt.Println("finishing fetching due to context closure")
			loop = false
		}
	}

	cancel()
	wg.Wait()
	close(comics)
	for comic := range comics {
		newCount++
		records[comic.Num] = *makeRecord(comic)
	}

	fmt.Println("fetch finished with", newCount, "new records")
	if newCount == 0 {
		return err
	}

	fmt.Println("saving records")
	dbErr := c.db.Write(records)
	return errors.Join(err, dbErr)
}

func (c *ComicsService) fetchJob(ctx context.Context, ids <-chan int, comics chan<- *xkcd.Comic, errs chan<- error) {
	for {
		select {
		case id := <-ids:
			if id == 404 {
				comics <- &xkcd.Comic{Num: 404}
				continue
			}
			comic, err := c.c.GetById(id)
			if err != nil {
				errs <- err
				return
			}
			comics <- comic
		case <-ctx.Done():
			return
		}
	}
}

func makeRecord(comic *xkcd.Comic) *database.Record {
	stemmed := words.Stem(comic.Title + " " + comic.Alt + " " + comic.Transcript)
	return database.NewRecord(comic.Img, stemmed)
}

package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
	"yadro-go/internal/database"
	"yadro-go/internal/xkcd"
	"yadro-go/pkg/logger"
	"yadro-go/pkg/stemmer"
)

var ErrUpdateInProgress = errors.New("update already in progress")

type Updater struct {
	log      *slog.Logger
	db       RecordRepository
	c        ComicFetcher
	limit    int
	parallel int
	mu       *sync.Mutex
}

type ComicFetcher interface {
	GetById(id int) (*xkcd.Comic, error)
}

type RecordRepository interface {
	Records() database.RecordMap
	Save(records database.RecordMap) error
}

func NewUpdater(log *slog.Logger, c ComicFetcher, db RecordRepository, limit int, parallel int) *Updater {
	return &Updater{
		log:      log,
		db:       db,
		c:        c,
		limit:    limit,
		parallel: parallel,
		mu:       &sync.Mutex{},
	}
}

func (u *Updater) StartScheduler(ctx context.Context, hour int, minute int) {
	const op = "updater.StartScheduler"
	log := u.log.With(slog.String("op", op))

	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}

	duration := next.Sub(now)
	timer := time.NewTimer(duration)

	log.Debug(fmt.Sprintf("scheduler started: next schedule time %v", next))

	for {
		select {
		case <-timer.C:
			log.Debug("update by scheduler")
			if _, err := u.Update(ctx); err != nil {
				log.Error("scheduled update error", logger.Err(err))
			}

			next = next.Add(24 * time.Hour)
			timer.Reset(next.Sub(time.Now()))
			log.Debug(fmt.Sprintf("next schedule time %v", next))
		case <-ctx.Done():
			log.Debug("scheduler stopped")
			timer.Stop()
			return
		}
	}
}

func (u *Updater) Update(ctx context.Context) (int, error) {
	const op = "updater.Update"
	log := u.log.With(slog.String("op", op))

	if !u.mu.TryLock() {
		log.Warn("update already in progress")
		return 0, ErrUpdateInProgress
	}
	defer u.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	ids := make(chan int, u.parallel)
	comics := make(chan *xkcd.Comic, u.parallel)
	errs := make(chan error, u.parallel)

	records := u.db.Records()
	log.Debug(fmt.Sprintf("start fetching with initial db size %d", len(records)))

	fetchId := 1
	pushId := func() bool {
		for ; fetchId <= u.limit; fetchId++ {
			_, ok := records[fetchId]
			if !ok {
				ids <- fetchId
				fetchId++
				return true
			}
		}

		return false
	}

	for i := 0; i < u.parallel; i++ {
		if !pushId() {
			if i == 0 {
				log.Debug("fetch finished: nothing to fetch")
				cancel()
				return len(records), nil
			}
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.fetchJob(ctx, ids, comics, errs)
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
				log.Debug("finishing fetching due to worker error")
				err = workerErr
			}
			loop = false
		case <-ctx.Done():
			log.Debug("finishing fetching due to context closure")
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

	log.Info(fmt.Sprintf("fetch finished: %d new records", newCount))
	if newCount == 0 {
		return len(records), err
	}

	dbErr := u.db.Save(records)
	return len(records), errors.Join(err, dbErr)
}

func (u *Updater) fetchJob(ctx context.Context, ids <-chan int, comics chan<- *xkcd.Comic, errs chan<- error) {
	for {
		select {
		case id := <-ids:
			if id == 404 {
				comics <- &xkcd.Comic{Num: 404}
				continue
			}
			comic, err := u.c.GetById(id)
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
	stemmed := stemmer.Stem(comic.Title + " " + comic.Alt + " " + comic.Transcript)
	return database.NewRecord(comic.Img, stemmed)
}

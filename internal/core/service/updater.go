package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"log/slog"
	"sync"
	"time"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	"yadro-go/pkg/logger"
)

type Updater struct {
	log         *slog.Logger
	stemmer     Stemmer
	comicRepo   ComicRepository
	keywordRepo KeywordRepository
	cp          ComicProvider
	limit       int
	parallel    int
	mu          *sync.Mutex
}

func NewUpdater(
	log *slog.Logger,
	stemmer Stemmer,
	comicRepo ComicRepository,
	keywordRepo KeywordRepository,
	c ComicProvider,
	limit int,
	parallel int,
) *Updater {
	return &Updater{
		log:         log,
		stemmer:     stemmer,
		comicRepo:   comicRepo,
		keywordRepo: keywordRepo,
		cp:          c,
		limit:       limit,
		parallel:    parallel,
		mu:          &sync.Mutex{},
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
		return 0, fmt.Errorf("%s: %w", op, ErrUpdateInProgress)
	}
	defer u.mu.Unlock()

	log.Debug("updating")

	comics, err := u.comicRepo.All(ctx)
	if err != nil {
		log.Error("failed to get all comics")
		return 0, fmt.Errorf("%s: %w", op, ErrInternal)
	}

	comicsMap := make(map[int]*domain.Comic, len(comics))
	for _, comic := range comics {
		comicsMap[comic.Num] = comic
	}

	log.Debug(fmt.Sprintf("start fetching with initial comics size %d", len(comicsMap)))

	jobCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	ids := make(chan int, u.parallel)
	res := make(chan *domain.Comic, u.parallel)
	errs := make(chan error, u.parallel)

	fetchId := 1
	pushId := func() bool {
		for ; fetchId <= u.limit; fetchId++ {
			_, ok := comicsMap[fetchId]
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
				return len(comicsMap), nil
			}
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.fetchJob(jobCtx, ids, res, errs)
		}()
	}

	newCount := 0
	loop := true
	for loop {
		select {
		case comic := <-res:
			newCount++
			comicsMap[comic.Num] = comic
			if !pushId() {
				loop = false
			}
		case workerErr := <-errs:
			if !errors.Is(workerErr, secondary.ErrComicNotFound) {
				log.Error("finishing update due to worker error", logger.Err(workerErr))
				err = workerErr
			}
			loop = false
		case <-ctx.Done():
			log.Debug("finishing update due to context closure")
			loop = false
		}
	}

	cancel()
	wg.Wait()
	close(res)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, ErrInternal)
	}

	for comic := range res {
		newCount++
		comicsMap[comic.Num] = comic
	}

	if newCount == 0 {
		log.Debug("update finished, no new records")
		return len(comicsMap), err
	}

	comics = maps.Values(comicsMap)
	if err = u.comicRepo.Save(ctx, comics); err != nil {
		log.Error("failed to save comics", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, ErrInternal)
	}

	if err = u.updateKeywords(ctx, comics); err != nil {
		log.Error("failed to update keywords", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, ErrInternal)
	}

	log.Debug(fmt.Sprintf("update finished: %d new comics", newCount))
	return len(comics), nil
}

func (u *Updater) updateKeywords(ctx context.Context, comics []*domain.Comic) error {
	const op = "updater.updateIndex"
	log := u.log.With(slog.String("op", op))

	log.Debug("updating keywords")
	keywordsMap := make(map[string]*domain.ComicKeyword)

	for _, comic := range comics {
		stemmed := u.stemmer.StemComic(comic)
		for _, word := range stemmed {
			keyword, ok := keywordsMap[word]
			if ok {
				keyword.Nums = append(keyword.Nums, comic.Num)
			} else {
				keywordsMap[word] = &domain.ComicKeyword{Word: word, Nums: []int{comic.Num}}
			}
		}
	}

	if err := u.keywordRepo.Save(ctx, maps.Values(keywordsMap)); err != nil {
		log.Error("failed to save keywords", logger.Err(err))
		return err
	}

	log.Debug("update keywords finished")
	return nil
}

func (u *Updater) fetchJob(ctx context.Context, ids <-chan int, comics chan<- *domain.Comic, errs chan<- error) {
	for {
		select {
		case id := <-ids:
			if id == 404 {
				comics <- &domain.Comic{Num: 404}
				continue
			}
			comic, err := u.cp.GetById(id)
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

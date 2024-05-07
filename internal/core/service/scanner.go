package service

import (
	"context"
	"fmt"
	"golang.org/x/exp/maps"
	"log/slog"
	"slices"
	"yadro-go/internal/core/domain"
	"yadro-go/pkg/logger"
)

type Scanner struct {
	log         *slog.Logger
	stemmer     Stemmer
	comicRepo   ComicRepository
	keywordRepo KeywordRepository
}

type NumMatch struct {
	num   int
	match int
}

func NewScanner(log *slog.Logger, stemmer Stemmer, comicRepo ComicRepository, keywordRepo KeywordRepository) *Scanner {
	return &Scanner{
		log:         log,
		stemmer:     stemmer,
		comicRepo:   comicRepo,
		keywordRepo: keywordRepo,
	}
}

func (s *Scanner) Scan(ctx context.Context, query string, useIndex bool) ([]string, error) {
	words := s.stemmer.StemString(query)

	if useIndex {
		return s.scanKeywords(ctx, words)
	}

	return s.scanComics(ctx, words)
}

func (s *Scanner) scanComics(ctx context.Context, words []string) ([]string, error) {
	const op = "scanner.scanComics"
	log := s.log.With(slog.String("op", op))

	log.Debug("scanning comics")

	comics, err := s.comicRepo.All(ctx)
	if err != nil {
		log.Error("failed to get comics", logger.Err(err))
		return nil, err
	}

	wordsSet := make(map[string]bool, len(words))
	for _, word := range words {
		wordsSet[word] = true
	}

	matches := make([]*NumMatch, 0)
	for num, comic := range comics {
		matchCount := 0
		select {
		case <-ctx.Done():
			log.Warn("scanning stopped, finishing")
			return nil, ctx.Err()

		default:
			for _, keyword := range s.stemmer.StemComic(comic) {
				if wordsSet[keyword] {
					matchCount++
				}
			}
			if matchCount > 0 {
				matches = append(matches, &NumMatch{num: num, match: matchCount})
			}
		}
	}

	log.Debug(fmt.Sprintf("scan finished: found %d matches", len(matches)))
	return finalizeResult(comics, matches), nil
}

func (s *Scanner) scanKeywords(ctx context.Context, words []string) ([]string, error) {
	const op = "scanner.scanKeywords"
	log := s.log.With(slog.String("op", op))

	log.Debug("scanning index")

	keywords, err := s.keywordRepo.Keywords(ctx, words)
	if err != nil {
		log.Error("failed to get keywords", logger.Err(err))
		return nil, err
	}

	matches := make(map[int]*NumMatch)
	for _, keyword := range keywords {
		if len(keyword.Nums) == 0 {
			continue
		}

		select {
		case <-ctx.Done():
			log.Warn("scanning stopped, finishing")
			return nil, ctx.Err()

		default:
			for _, num := range keyword.Nums {
				numMatch, ok := matches[num]
				if !ok {
					matches[num] = &NumMatch{num: num, match: 1}
				} else {
					numMatch.match++
				}
			}
		}
	}

	comics, err := s.comicRepo.Comics(ctx, maps.Keys(matches))
	if err != nil {
		log.Error("failed to get comics")
		return nil, err
	}

	log.Debug(fmt.Sprintf("scan finished: found %d matches", len(matches)))
	return finalizeResult(comics, maps.Values(matches)), nil
}

func finalizeResult(comics []*domain.Comic, matches []*NumMatch) []string {
	slices.SortFunc(matches, func(a, b *NumMatch) int {
		return b.match - a.match
	})

	comicMap := make(map[int]*domain.Comic, len(comics))
	for _, comic := range comics {
		comicMap[comic.Num] = comic
	}

	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = comicMap[match.num].Img
	}

	return result
}

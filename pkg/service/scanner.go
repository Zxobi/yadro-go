package service

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"yadro-go/pkg/database"
	"yadro-go/pkg/stem"
)

type Scanner struct {
	log *slog.Logger
	rp  RecordProvider
	ip  IndexProvider
}

type NumMatch struct {
	num   int
	match int
}

type RecordProvider interface {
	Records() database.RecordMap
}

type IndexProvider interface {
	Index() database.IndexMap
}

func NewScanner(log *slog.Logger, rp RecordProvider, ip IndexProvider) *Scanner {
	return &Scanner{
		log: log,
		rp:  rp,
		ip:  ip,
	}
}

func (s *Scanner) Scan(ctx context.Context, query string, useIndex bool) []string {
	words := stem.Stem(query)

	if useIndex {
		return s.scanIndex(ctx, words)
	}

	return s.scanRecords(ctx, words)
}

func (s *Scanner) scanRecords(ctx context.Context, words []string) []string {
	s.log.Info("scanning records")

	records := s.rp.Records()

	wordsSet := make(map[string]bool, len(words))
	for _, word := range words {
		wordsSet[word] = true
	}

	matches := make([]NumMatch, 0)
	for num, record := range records {
		matchCount := 0
		select {
		case <-ctx.Done():
			s.log.Info("scanning stopped, finishing")
			return nil

		default:
			for _, keyword := range record.Keywords {
				if wordsSet[keyword] {
					matchCount++
				}
			}
			if matchCount > 0 {
				matches = append(matches, NumMatch{num: num, match: matchCount})
			}
		}
	}

	s.log.Info(fmt.Sprintf("scan finished: found %d matches", len(matches)))
	return finalizeResult(records, matches)
}

func (s *Scanner) scanIndex(ctx context.Context, words []string) []string {
	s.log.Info("scanning index")
	index := s.ip.Index()

	matches := make(map[int]int)
	for _, word := range words {
		nums := index[word]
		if len(nums) == 0 {
			continue
		}

		select {
		case <-ctx.Done():
			s.log.Info("scanning stopped, finishing")
			return nil

		default:
			for _, num := range nums {
				match := matches[num]
				if match == 0 {
					matches[num] = 1
				} else {
					matches[num]++
				}
			}
		}
	}

	numMatches := make([]NumMatch, 0, len(matches))
	for num, match := range matches {
		numMatches = append(numMatches, NumMatch{num, match})
	}

	s.log.Info(fmt.Sprintf("scan finished: found %d matches", len(numMatches)))
	return finalizeResult(s.rp.Records(), numMatches)
}

func finalizeResult(records database.RecordMap, matches []NumMatch) []string {
	slices.SortFunc(matches, func(a, b NumMatch) int {
		return b.match - a.match
	})

	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = records[match.num].Url
	}

	return result
}

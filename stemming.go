package main

import (
	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
	"strings"
)

func stem(s string) ([]string, error) {
	words := strings.Split(strings.ToLower(strings.TrimSpace(s)), " ")

	stemmedWordsSet := make(map[string]struct{})
	for _, v := range words {
		if shouldIgnore(v) {
			continue
		}

		stemmed, err := snowball.Stem(v, "english", false)
		if err != nil {
			return nil, err
		}

		stemmedWordsSet[stemmed] = struct{}{}
	}

	var stemmedWordsSlice []string
	for v, _ := range stemmedWordsSet {
		stemmedWordsSlice = append(stemmedWordsSlice, v)
	}

	return stemmedWordsSlice, nil
}

func shouldIgnore(s string) bool {
	if isStopWord(s) {
		return true
	}

	ind := strings.IndexFunc(s, func(r rune) bool {
		return r == '\''
	})
	if ind >= 0 {
		if isStopWord(s[:ind]) {
			return true
		}
	}

	return false
}

func isStopWord(s string) bool {
	return len(s) <= 2 || english.IsStopWord(s)
}

package main

import (
	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

func stem(s string) ([]string, error) {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !(unicode.IsLetter(r) || r == '\'' || r == '-')
	})

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

	return ind > 0
}

func isStopWord(s string) bool {
	return len(s) <= 2 || english.IsStopWord(s)
}

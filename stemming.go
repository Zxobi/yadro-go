package main

import (
	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

var stemChain = []func(string) (string, bool, error){
	stemIgnore,
	stemApostrophe,
	stemHyphen,
	stemInternal,
}

func stem(s string) ([]string, error) {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !(unicode.IsLetter(r) || r == '\'' || r == '-')
	})

	stemmedWordsSet := make(map[string]struct{})
	for _, v := range words {
		var res string
		for _, f := range stemChain {
			if stemmed, ok, err := f(v); err != nil {
				return []string{}, err
			} else if ok {
				res = stemmed
				break
			}
		}

		if len(res) == 0 {
			continue
		}
		stemmedWordsSet[res] = struct{}{}
	}

	stemmedWordsSlice := make([]string, 0, len(stemmedWordsSet))
	for v := range stemmedWordsSet {
		stemmedWordsSlice = append(stemmedWordsSlice, v)
	}

	return stemmedWordsSlice, nil
}

func stemIgnore(word string) (string, bool, error) {
	if isStopWord(word) {
		return "", true, nil
	}

	return word, false, nil
}

func stemApostrophe(word string) (string, bool, error) {
	ind := strings.Index(word, "'")
	if ind < 0 {
		return word, false, nil
	}

	_, ignore, err := stemIgnore(word[:ind])
	if ignore || err != nil {
		return "", true, err
	}

	return word, true, nil
}

func stemHyphen(word string) (string, bool, error) {
	ind := strings.LastIndex(word, "-")
	if ind < 0 {
		return word, false, nil
	}

	stemmed, _, err := stemInternal(word[ind+1:])
	if err != nil {
		return "", true, err
	}

	return word[:ind+1] + stemmed, true, nil
}

func stemInternal(word string) (string, bool, error) {
	stemmed, err := snowball.Stem(word, "english", false)
	if err != nil {
		return "", true, err
	}

	return stemmed, true, nil
}

func isStopWord(s string) bool {
	return len(s) <= 2 || english.IsStopWord(s)
}

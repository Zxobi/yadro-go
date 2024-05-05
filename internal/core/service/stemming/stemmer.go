package stemming

import (
	"github.com/kljensen/snowball/english"
	"strings"
	"unicode"
	"yadro-go/internal/core/domain"
)

type Stemmer struct {
}

func New() *Stemmer {
	return &Stemmer{}
}

func (s *Stemmer) StemComic(comic *domain.Comic) []string {
	return s.StemString(comic.Title + " " + comic.Alt + " " + comic.Transcript)
}

func (s *Stemmer) StemString(str string) []string {
	words := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	stemmedWordsSet := make(map[string]struct{})
	for _, v := range words {
		stemmed := english.Stem(v, false)
		if shouldIgnore(stemmed) {
			continue
		}

		stemmedWordsSet[stemmed] = struct{}{}
	}

	stemmedWordsSlice := make([]string, 0, len(stemmedWordsSet))
	for v := range stemmedWordsSet {
		stemmedWordsSlice = append(stemmedWordsSlice, v)
	}

	return stemmedWordsSlice

}

func shouldIgnore(s string) bool {
	return len(s) <= 2 || english.IsStopWord(s)
}

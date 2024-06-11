package stemming

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var stemmer *Stemmer

func init() {
	stemmer = New()
}

func TestStemmer_StemString(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"follower brings bunch of questions", []string{"follow", "bring", "bunch", "question"}},
		{"i'll follow you as long as you are following me", []string{"follow", "long"}},
		{"foll[o]wer,,brin'gs       bunch-of-questions",
			[]string{"foll", "wer", "brin", "bunch", "question"}},
	}

	for _, testCase := range testTable {
		assert.ElementsMatch(t, stemmer.StemString(testCase.input), testCase.expected)
	}
}

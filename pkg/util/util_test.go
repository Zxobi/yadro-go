package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePlaceholders(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		count    int
		expected string
	}{
		{0, ""},
		{1, "?"},
		{2, "?,?"},
		{5, "?,?,?,?,?"},
	}

	for _, testCase := range testTable {
		assert.Equal(t, testCase.expected, GeneratePlaceholders(testCase.count))
	}
}

func TestSliceToAny(t *testing.T) {
	t.Parallel()

	s := make([]string, 0)
	assert.ElementsMatch(t, s, SliceToAny(s))
	s = []string{"str1", "str2", "str3"}
	assert.ElementsMatch(t, s, SliceToAny(s))
	i := []int{1, 2, 3, 4, 5}
	assert.ElementsMatch(t, i, SliceToAny(i))
}

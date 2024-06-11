package matcher

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"yadro-go/internal/core/domain"
)

type keywordPtrSliceMatcher struct {
	expected []*domain.ComicKeyword
}

func (p keywordPtrSliceMatcher) Matches(x interface{}) bool {
	actual, ok := x.([]*domain.ComicKeyword)
	if !ok {
		return false
	}

	if len(p.expected) != len(actual) {
		return false
	}

	expectedMap := make(map[string]*domain.ComicKeyword)
	for _, v := range p.expected {
		expectedMap[v.Word] = v
	}

	for _, v := range actual {
		expected := expectedMap[v.Word]
		if expected == nil || len(expected.Nums) != len(v.Nums) {
			return false
		}
		delete(expectedMap, v.Word)
	}

	return len(expectedMap) == 0
}

func (p keywordPtrSliceMatcher) String() string {
	return fmt.Sprintf("is equal to %v", p.expected)
}

func KeywordPtrSliceEqual(expected []*domain.ComicKeyword) gomock.Matcher {
	return keywordPtrSliceMatcher{expected: expected}
}

package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	mock_service "yadro-go/internal/core/service/mocks"
	"yadro-go/test/logger"
	"yadro-go/test/matcher"
)

func TestUpdater_Update(t *testing.T) {
	t.Parallel()

	var comic1 = &domain.Comic{Num: 1, Title: "test"}
	var comic2 = &domain.Comic{Num: 2, Title: "test", Alt: "test_alt"}
	var comic3 = &domain.Comic{Num: 3, Title: "test", Transcript: "test_transcript"}
	var keyword1 = &domain.ComicKeyword{Word: "test", Nums: []int{1, 2, 3}}
	var keyword1Limited = &domain.ComicKeyword{Word: "test", Nums: []int{1, 2}}
	var keyword2 = &domain.ComicKeyword{Word: "test_alt", Nums: []int{2}}
	var keyword3 = &domain.ComicKeyword{Word: "test_transcript", Nums: []int{3}}

	testTable := []struct {
		name                       string
		parallel                   int
		limit                      int
		comicProviderBehaviour     func(provider *mock_service.MockComicProvider)
		comicRepositoryBehaviour   func(repo *mock_service.MockComicRepository)
		keywordRepositoryBehaviour func(repo *mock_service.MockKeywordRepository)
		stemmerBehaviour           func(stemmer *mock_service.MockStemmer)
		expectedCount              int
		expectedError              error
	}{
		{
			name:     "SuccessNoParallel",
			parallel: 1,
			limit:    1000,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				gomock.InOrder(
					provider.EXPECT().GetById(1).Return(comic1, nil),
					provider.EXPECT().GetById(2).Return(comic2, nil),
					provider.EXPECT().GetById(3).Return(comic3, nil),
					provider.EXPECT().GetById(4).Return(nil, secondary.ErrComicNotFound),
				)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{}, nil)
				repo.EXPECT().Save(gomock.Any(), gomock.InAnyOrder([]*domain.Comic{comic1, comic2, comic3})).
					Return(nil)
			},
			keywordRepositoryBehaviour: func(repo *mock_service.MockKeywordRepository) {
				repo.EXPECT().Save(gomock.Any(), matcher.KeywordPtrSliceEqual([]*domain.ComicKeyword{
					keyword1, keyword2, keyword3,
				})).Return(nil)
			},
			stemmerBehaviour: func(stemmer *mock_service.MockStemmer) {
				stemmer.EXPECT().StemComic(comic1).Return([]string{"test"})
				stemmer.EXPECT().StemComic(comic2).Return([]string{"test", "test_alt"})
				stemmer.EXPECT().StemComic(comic3).Return([]string{"test", "test_transcript"})
			},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:     "SuccessParallel",
			parallel: 2,
			limit:    1000,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				provider.EXPECT().GetById(1).Return(comic1, nil)
				provider.EXPECT().GetById(2).Return(comic2, nil)
				provider.EXPECT().GetById(3).Return(comic3, nil)
				provider.EXPECT().GetById(gomock.Any()).AnyTimes().Return(nil, secondary.ErrComicNotFound)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{}, nil)
				repo.EXPECT().Save(gomock.Any(), gomock.InAnyOrder([]*domain.Comic{comic1, comic2, comic3})).
					Return(nil)
			},
			keywordRepositoryBehaviour: func(repo *mock_service.MockKeywordRepository) {
				repo.EXPECT().Save(gomock.Any(), matcher.KeywordPtrSliceEqual([]*domain.ComicKeyword{
					keyword1, keyword2, keyword3,
				})).Return(nil)
			},
			stemmerBehaviour: func(stemmer *mock_service.MockStemmer) {
				stemmer.EXPECT().StemComic(comic1).Return([]string{"test"})
				stemmer.EXPECT().StemComic(comic2).Return([]string{"test", "test_alt"})
				stemmer.EXPECT().StemComic(comic3).Return([]string{"test", "test_transcript"})
			},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:     "SuccessOnlyNewComicsFetched",
			parallel: 1,
			limit:    1000,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				gomock.InOrder(
					provider.EXPECT().GetById(2).Return(comic2, nil),
					provider.EXPECT().GetById(3).Return(comic3, nil),
					provider.EXPECT().GetById(4).Return(nil, secondary.ErrComicNotFound),
				)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{comic1}, nil)
				repo.EXPECT().Save(gomock.Any(), gomock.InAnyOrder([]*domain.Comic{comic1, comic2, comic3})).
					Return(nil)
			},
			keywordRepositoryBehaviour: func(repo *mock_service.MockKeywordRepository) {
				repo.EXPECT().Save(gomock.Any(), matcher.KeywordPtrSliceEqual([]*domain.ComicKeyword{
					keyword1, keyword2, keyword3,
				})).Return(nil)
			},
			stemmerBehaviour: func(stemmer *mock_service.MockStemmer) {
				stemmer.EXPECT().StemComic(comic1).Return([]string{"test"})
				stemmer.EXPECT().StemComic(comic2).Return([]string{"test", "test_alt"})
				stemmer.EXPECT().StemComic(comic3).Return([]string{"test", "test_transcript"})
			},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:     "LimitReached",
			parallel: 1,
			limit:    2,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				provider.EXPECT().GetById(1).Return(comic1, nil)
				provider.EXPECT().GetById(2).Return(comic2, nil)
				provider.EXPECT().GetById(gomock.Any()).AnyTimes().Return(comic3, nil)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{}, nil)
				repo.EXPECT().Save(gomock.Any(), gomock.InAnyOrder([]*domain.Comic{comic1, comic2})).
					Return(nil)
			},
			keywordRepositoryBehaviour: func(repo *mock_service.MockKeywordRepository) {
				repo.EXPECT().Save(gomock.Any(), matcher.KeywordPtrSliceEqual([]*domain.ComicKeyword{
					keyword1Limited, keyword2,
				})).Return(nil)
			},
			stemmerBehaviour: func(stemmer *mock_service.MockStemmer) {
				stemmer.EXPECT().StemComic(comic1).Return([]string{"test"})
				stemmer.EXPECT().StemComic(comic2).Return([]string{"test", "test_alt"})
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name:     "NothingToUpdate",
			parallel: 1,
			limit:    1000,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				provider.EXPECT().GetById(4).Return(nil, secondary.ErrComicNotFound)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{comic1, comic2, comic3}, nil)
			},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:     "NothingToUpdateLimitReached",
			parallel: 1,
			limit:    3,
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{comic1, comic2, comic3}, nil)
			},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:     "ComicRepositoryAllError",
			parallel: 1,
			limit:    1000,
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return(nil, secondary.ErrInternal)
			},
			expectedError: ErrInternal,
		},
		{
			name:     "ComicRepositorySaveError",
			parallel: 1,
			limit:    1000,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				gomock.InOrder(
					provider.EXPECT().GetById(2).Return(comic2, nil),
					provider.EXPECT().GetById(3).Return(comic3, nil),
					provider.EXPECT().GetById(4).Return(nil, secondary.ErrComicNotFound),
				)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{comic1}, nil)
				repo.EXPECT().Save(gomock.Any(), gomock.InAnyOrder([]*domain.Comic{comic1, comic2, comic3})).
					Return(secondary.ErrInternal)
			},
			expectedError: ErrInternal,
		},
		{
			name:     "KeywordRepositorySaveError",
			parallel: 1,
			limit:    1000,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				gomock.InOrder(
					provider.EXPECT().GetById(2).Return(comic2, nil),
					provider.EXPECT().GetById(3).Return(comic3, nil),
					provider.EXPECT().GetById(4).Return(nil, secondary.ErrComicNotFound),
				)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{comic1}, nil)
				repo.EXPECT().Save(gomock.Any(), gomock.InAnyOrder([]*domain.Comic{comic1, comic2, comic3})).
					Return(nil)
			},
			keywordRepositoryBehaviour: func(repo *mock_service.MockKeywordRepository) {
				repo.EXPECT().Save(gomock.Any(), matcher.KeywordPtrSliceEqual([]*domain.ComicKeyword{
					keyword1, keyword2, keyword3,
				})).Return(secondary.ErrInternal)
			},
			stemmerBehaviour: func(stemmer *mock_service.MockStemmer) {
				stemmer.EXPECT().StemComic(comic1).Return([]string{"test"})
				stemmer.EXPECT().StemComic(comic2).Return([]string{"test", "test_alt"})
				stemmer.EXPECT().StemComic(comic3).Return([]string{"test", "test_transcript"})
			},
			expectedError: ErrInternal,
		},
		{
			name:     "ComicProviderUnexpectedError",
			parallel: 1,
			limit:    1000,
			comicProviderBehaviour: func(provider *mock_service.MockComicProvider) {
				gomock.InOrder(
					provider.EXPECT().GetById(2).Return(nil, secondary.ErrInternal),
				)
			},
			comicRepositoryBehaviour: func(repo *mock_service.MockComicRepository) {
				repo.EXPECT().All(gomock.Any()).Return([]*domain.Comic{comic1}, nil)
			},
			expectedError: ErrInternal,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			comicProvider := mock_service.NewMockComicProvider(c)
			comicRepo := mock_service.NewMockComicRepository(c)
			keywordRepo := mock_service.NewMockKeywordRepository(c)
			stemmer := mock_service.NewMockStemmer(c)

			if testCase.comicProviderBehaviour != nil {
				testCase.comicProviderBehaviour(comicProvider)
			}
			if testCase.comicRepositoryBehaviour != nil {
				testCase.comicRepositoryBehaviour(comicRepo)
			}
			if testCase.stemmerBehaviour != nil {
				testCase.stemmerBehaviour(stemmer)
			}
			if testCase.keywordRepositoryBehaviour != nil {
				testCase.keywordRepositoryBehaviour(keywordRepo)
			}

			u := NewUpdater(slog.New(logger.EmptyHandler{}), stemmer, comicRepo, keywordRepo, comicProvider,
				testCase.limit, testCase.parallel)
			count, err := u.Update(context.Background())
			require.ErrorIs(t, err, testCase.expectedError)
			if testCase.expectedError == nil {
				assert.Equal(t, testCase.expectedCount, count)
			}
		})
	}
}

func TestUpdater_UpdateInProgress(t *testing.T) {
	t.Parallel()

	c := gomock.NewController(t)
	comicProvider := mock_service.NewMockComicProvider(c)
	comicRepo := mock_service.NewMockComicRepository(c)
	keywordRepo := mock_service.NewMockKeywordRepository(c)
	stemmer := mock_service.NewMockStemmer(c)

	ctx, cancel := context.WithCancel(context.Background())

	waitChan := make(chan bool)
	comicRepo.EXPECT().All(gomock.Any()).DoAndReturn(func(ctx context.Context) ([]*domain.Comic, error) {
		waitChan <- true
		for {
			select {
			case <-ctx.Done():
				return nil, secondary.ErrInternal
			default:
			}
		}
	})

	u := NewUpdater(slog.New(logger.EmptyHandler{}), stemmer, comicRepo, keywordRepo, comicProvider, 1000, 1)
	go func() {
		_, _ = u.Update(ctx)
	}()
	<-waitChan

	_, err := u.Update(ctx)
	require.ErrorIs(t, err, ErrUpdateInProgress)
	cancel()
}

func TestUpdater_StartSchedulerNotPanic(t *testing.T) {
	t.Parallel()

	c := gomock.NewController(t)
	comicProvider := mock_service.NewMockComicProvider(c)
	comicRepo := mock_service.NewMockComicRepository(c)
	keywordRepo := mock_service.NewMockKeywordRepository(c)
	stemmer := mock_service.NewMockStemmer(c)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	u := NewUpdater(slog.New(logger.EmptyHandler{}), stemmer, comicRepo, keywordRepo, comicProvider, 1000, 1)
	assert.NotPanics(t, func() { u.StartScheduler(ctx, 0, 0) })
}

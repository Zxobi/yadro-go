package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	mock_service "yadro-go/internal/core/service/mocks"
	"yadro-go/test/logger"
)

var testPasswordHash = "$2a$10$ObD/z7BvOq2AKUcJZ9FhOeTNF2SKFRqbGw7TwbyiDoQvs7Fsq42H." //testpassword

func TestAuth_Login(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name              string
		username          string
		password          string
		userRepoBehaviour func(repo *mock_service.MockUserRepository)
		tmBehaviour       func(tm *mock_service.MockTokenManager)
		expectedToken     string
		expectedError     error
	}{
		{
			name:     "Success",
			username: "test_user",
			password: "testpassword",
			userRepoBehaviour: func(repo *mock_service.MockUserRepository) {
				repo.EXPECT().UserByUsername(gomock.Any(), "test_user").Return(
					&domain.User{Username: "test_user", PassHash: []byte(testPasswordHash)}, nil,
				)
			},
			tmBehaviour: func(tm *mock_service.MockTokenManager) {
				tm.EXPECT().Token("test_user").Return("token", nil)
			},
			expectedToken: "token",
			expectedError: nil,
		},
		{
			name:     "WrongCredentials",
			username: "test_user",
			password: "wrong",
			userRepoBehaviour: func(repo *mock_service.MockUserRepository) {
				repo.EXPECT().UserByUsername(gomock.Any(), "test_user").Return(
					&domain.User{Username: "test_user", PassHash: []byte(testPasswordHash)}, nil,
				)
			},
			expectedError: ErrWrongCredentials,
		},
		{
			name:     "UserNotFound",
			username: "not_exists",
			password: "testpassword",
			userRepoBehaviour: func(repo *mock_service.MockUserRepository) {
				repo.EXPECT().UserByUsername(gomock.Any(), "not_exists").Return(
					nil, secondary.ErrUserNotFound,
				)
			},
			expectedError: ErrWrongCredentials,
		},
		{
			name:     "TokenFails",
			username: "test_user",
			password: "testpassword",
			userRepoBehaviour: func(repo *mock_service.MockUserRepository) {
				repo.EXPECT().UserByUsername(gomock.Any(), "test_user").Return(
					&domain.User{Username: "test_user", PassHash: []byte(testPasswordHash)}, nil,
				)
			},
			tmBehaviour: func(tm *mock_service.MockTokenManager) {
				tm.EXPECT().Token("test_user").Return("", errors.New("token failed"))
			},
			expectedError: ErrInternal,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			userRepo := mock_service.NewMockUserRepository(c)
			tm := mock_service.NewMockTokenManager(c)
			if testCase.userRepoBehaviour != nil {
				testCase.userRepoBehaviour(userRepo)
			}
			if testCase.tmBehaviour != nil {
				testCase.tmBehaviour(tm)
			}

			auth := NewAuth(slog.New(logger.EmptyHandler{}), tm, userRepo)
			token, err := auth.Login(context.Background(), testCase.username, testCase.password)
			require.ErrorIs(t, err, testCase.expectedError)
			if testCase.expectedError == nil {
				assert.Equal(t, testCase.expectedToken, token)
			}
		})
	}
}

func TestAuth_Verify(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name              string
		token             string
		userRepoBehaviour func(repo *mock_service.MockUserRepository)
		tmBehaviour       func(tm *mock_service.MockTokenManager)
		expectedUser      *domain.User
		expectedError     error
	}{
		{
			name:  "Success",
			token: "token",
			userRepoBehaviour: func(repo *mock_service.MockUserRepository) {
				repo.EXPECT().UserByUsername(gomock.Any(), "test_user").Return(
					&domain.User{Username: "test_user"}, nil,
				)
			},
			tmBehaviour: func(tm *mock_service.MockTokenManager) {
				tm.EXPECT().Verify("token").Return("test_user", nil)
			},
			expectedUser:  &domain.User{Username: "test_user"},
			expectedError: nil,
		},
		{
			name:  "BadToken",
			token: "bad_token",
			tmBehaviour: func(tm *mock_service.MockTokenManager) {
				tm.EXPECT().Verify("bad_token").Return("", ErrBadToken)
			},
			expectedUser:  nil,
			expectedError: ErrBadToken,
		},
		{
			name:  "UserNotFound",
			token: "token",
			userRepoBehaviour: func(repo *mock_service.MockUserRepository) {
				repo.EXPECT().UserByUsername(gomock.Any(), "test_user").Return(
					nil, secondary.ErrUserNotFound,
				)
			},
			tmBehaviour: func(tm *mock_service.MockTokenManager) {
				tm.EXPECT().Verify("token").Return("test_user", nil)
			},
			expectedUser:  nil,
			expectedError: ErrInternal,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			userRepo := mock_service.NewMockUserRepository(c)
			tm := mock_service.NewMockTokenManager(c)
			if testCase.userRepoBehaviour != nil {
				testCase.userRepoBehaviour(userRepo)
			}
			if testCase.tmBehaviour != nil {
				testCase.tmBehaviour(tm)
			}

			auth := NewAuth(slog.New(logger.EmptyHandler{}), tm, userRepo)
			user, err := auth.Authenticate(context.Background(), testCase.token)
			require.ErrorIs(t, err, testCase.expectedError)
			if testCase.expectedError == nil {
				assert.Equal(t, testCase.expectedUser, user)
			}
		})
	}
}

package token

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"time"
	"yadro-go/internal/core/service"
	"yadro-go/test/logger"
)

func TestJwtTokenManager_TokenVerify(t *testing.T) {
	t.Parallel()

	tm := NewJwtTokenManager(slog.New(logger.EmptyHandler{}), []byte("secret"), time.Minute)

	testTable := []string{"", "test_user", "test_user_with_very_long_username_string"}

	for _, testCase := range testTable {
		token, err := tm.Token(testCase)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		verifiedUsername, err := tm.Verify(token)
		require.NoError(t, err)
		assert.Equal(t, testCase, verifiedUsername)
	}
}

func TestJwtTokenManager_VerifyExpiredTokenError(t *testing.T) {
	t.Parallel()

	tm := NewJwtTokenManager(slog.New(logger.EmptyHandler{}), []byte("secret"), 0)

	token, err := tm.Token("test_user")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, err = tm.Verify(token)
	require.ErrorIs(t, err, service.ErrBadToken)
}

func TestJwtTokenManager_VerifyTokenSignedWithDifferentSecretError(t *testing.T) {
	t.Parallel()

	tmToken := NewJwtTokenManager(slog.New(logger.EmptyHandler{}), []byte("secret_1"), time.Minute)
	tmVerify := NewJwtTokenManager(slog.New(logger.EmptyHandler{}), []byte("secret_2"), time.Minute)

	token, err := tmToken.Token("test_user")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, err = tmVerify.Verify(token)
	require.ErrorIs(t, err, service.ErrBadToken)
}

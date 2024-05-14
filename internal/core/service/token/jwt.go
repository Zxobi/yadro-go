package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
	"yadro-go/internal/core/service"
	"yadro-go/pkg/logger"
)

type JwtTokenManager struct {
	log      *slog.Logger
	secret   []byte
	tokenTTL time.Duration
}

func NewJwtTokenManager(log *slog.Logger, secret []byte, tokenTTL time.Duration) *JwtTokenManager {
	return &JwtTokenManager{log: log, secret: secret, tokenTTL: tokenTTL}
}

func (t *JwtTokenManager) Token(username string) (string, error) {
	const op = "token.Token"
	log := t.log.With(slog.String("op", op))

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   username,
	}).SignedString(t.secret)

	if err != nil {
		log.Error("failed to make token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", service.ErrInternal)
	}

	return tokenString, nil
}

func (t *JwtTokenManager) Verify(tokenString string) (string, error) {
	const op = "token.Verify"
	log := t.log.With(slog.String("op", op))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: invalid token signing method %v", op, token.Header["alg"])
		}

		return t.secret, nil
	})
	if err != nil {
		log.Error("failed to parse token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, service.ErrBadToken)
	}

	username, err := token.Claims.GetSubject()
	if err != nil {
		log.Error("subject is missing")
		return "", fmt.Errorf("%s: %w", op, service.ErrBadToken)
	}

	return username, nil
}

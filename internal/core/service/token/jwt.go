package token

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"log/slog"
	"yadro-go/internal/core/service"
	"yadro-go/pkg/logger"
)

type JwtTokenManager struct {
	log    *slog.Logger
	secret []byte
}

func NewJwtTokenManager(log *slog.Logger, secret []byte) *JwtTokenManager {
	return &JwtTokenManager{log: log, secret: secret}
}

func (t *JwtTokenManager) Token(username string) (string, error) {
	const op = "token.Token"
	log := t.log.With(slog.String("op", op))

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Error("failed to get token claims")
		return "", fmt.Errorf("%s: %w", op, service.ErrBadToken)
	}

	username, ok := claims["sub"].(string)
	if !ok {
		log.Error("failed to get token sub")
		return "", fmt.Errorf("%s: %w", op, service.ErrBadToken)
	}

	return username, nil
}

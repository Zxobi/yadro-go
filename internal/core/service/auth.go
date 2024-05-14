package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	"yadro-go/pkg/logger"
)

type Auth struct {
	log          *slog.Logger
	tokenManager TokenManager
	userRepo     UserRepository
}

func NewAuth(log *slog.Logger, tokenManager TokenManager, userRepo UserRepository) *Auth {
	return &Auth{
		log:          log,
		tokenManager: tokenManager,
		userRepo:     userRepo,
	}
}

func (a *Auth) Login(ctx context.Context, username string, password string) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op), slog.String("uname", username))

	log.Debug("logging in")

	user, err := a.userRepo.UserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, secondary.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrWrongCredentials)
		}

		log.Error("failed to get user", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInternal)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", fmt.Errorf("%s: %w", op, ErrWrongCredentials)
		}

		log.Error("failed to compare passwords", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrWrongCredentials)
	}

	token, err := a.tokenManager.Token(user.Username)
	if err != nil {
		log.Error("failed to make token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInternal)
	}

	log.Debug("logged in")

	return token, nil
}

func (a *Auth) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	const op = "auth.Authenticate"
	log := a.log.With(slog.String("op", op))

	log.Debug("authenticating")

	username, err := a.tokenManager.Verify(token)
	if err != nil {
		log.Error("failed to verify token", logger.Err(err))
		return nil, ErrBadToken
	}

	user, err := a.userRepo.UserByUsername(ctx, username)
	if err != nil {
		log.Error("failed to get user", logger.Err(err))
		return nil, ErrInternal
	}

	log.Debug("authenticated", slog.String("uname", user.Username))

	return user, nil
}

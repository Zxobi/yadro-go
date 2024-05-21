package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	"yadro-go/pkg/logger"
)

const (
	statementSelectUserByUsername = "SELECT username, role, pass_hash FROM users WHERE username=?"
)

type UserRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewUserRepository(log *slog.Logger, db *sql.DB) *UserRepository {
	return &UserRepository{log: log, db: db}
}

func (r *UserRepository) UserByUsername(ctx context.Context, username string) (*domain.User, error) {
	const op = "user.UserByUsername"
	log := r.log.With(slog.String("op", op), slog.String("uname", username))

	log.Debug("fetching user")

	stmt, err := r.db.PrepareContext(ctx, statementSelectUserByUsername)
	if err != nil {
		log.Error("failed to prepare statement", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		log.Error("failed to query users", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer rows.Close()

	if !rows.Next() {
		log.Debug("user not found")
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrUserNotFound)
	}

	var user domain.User
	if err = rows.Scan(&user.Username, &user.Role, &user.PassHash); err != nil {
		log.Error("failed to decode user", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	return &user, nil
}

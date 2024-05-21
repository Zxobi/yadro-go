package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	"yadro-go/pkg/logger"
	"yadro-go/pkg/util"
)

const (
	querySelectAllComics          = "SELECT * FROM comics"
	statementInsertOrReplaceComic = "INSERT OR REPLACE INTO comics(num, title, transcript, alt, img) VALUES (?, ?, ?, ?, ?)"
	formatStatementSelectComics   = "SELECT * FROM comics WHERE num IN (%s)"
)

type ComicRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewComicRepository(log *slog.Logger, db *sql.DB) *ComicRepository {
	return &ComicRepository{log: log, db: db}
}

func (r *ComicRepository) Comics(ctx context.Context, nums []int) ([]*domain.Comic, error) {
	const op = "comic.Comics"
	log := r.log.With(slog.String("op", op))

	log.Debug("fetching comics")

	stmt, err := r.db.PrepareContext(ctx, fmt.Sprintf(formatStatementSelectComics, util.GeneratePlaceholders(len(nums))))
	if err != nil {
		log.Error("failed to prepare statement", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, util.SliceToAny(nums)...)
	if err != nil {
		log.Error("failed to query comics", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer rows.Close()

	comics := make([]*domain.Comic, 0)

	for rows.Next() {
		var comic domain.Comic
		err = rows.Scan(&comic.Num, &comic.Title, &comic.Transcript, &comic.Alt, &comic.Img)
		if err != nil {
			log.Error("failed to decode comic", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
		}

		comics = append(comics, &comic)
	}

	if err = rows.Err(); err != nil {
		log.Error("error during rows iteration", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	log.Debug("fetch comics complete")

	return comics, nil
}

func (r *ComicRepository) All(ctx context.Context) ([]*domain.Comic, error) {
	const op = "comic.All"
	log := r.log.With(slog.String("op", op))

	log.Debug("fetching all comics")

	rows, err := r.db.QueryContext(ctx, querySelectAllComics)
	if err != nil {
		log.Error("failed to query all comics")
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer rows.Close()

	res := make([]*domain.Comic, 0)

	for rows.Next() {
		var comic domain.Comic

		err = rows.Scan(&comic.Num, &comic.Title, &comic.Transcript, &comic.Alt, &comic.Img)
		if err != nil {
			log.Error("failed to decode comic", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
		}

		res = append(res, &comic)
	}

	if err = rows.Err(); err != nil {
		log.Error("error during rows iteration", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	log.Debug("fetch all comics complete")

	return res, nil
}

func (r *ComicRepository) Save(ctx context.Context, comics []*domain.Comic) error {
	const op = "comic.Save"
	log := r.log.With(slog.String("op", op))

	log.Debug("saving comics")

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to start a transaction", logger.Err(err))
		return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	stmt, err := tx.PrepareContext(ctx, statementInsertOrReplaceComic)
	if err != nil {
		log.Error("failed to prepare statement", logger.Err(err))
		return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer stmt.Close()

	for _, comic := range comics {
		_, err = stmt.ExecContext(ctx, comic.Num, comic.Title, comic.Transcript, comic.Alt, comic.Img)
		if err != nil {
			log.Error("failed to execute statement", logger.Err(err))
			if err = tx.Rollback(); err != nil {
				log.Error("tx rollback failed", logger.Err(err))
			}
			return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx commit failed", logger.Err(err))
		return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	log.Debug("save comics complete")
	return nil
}

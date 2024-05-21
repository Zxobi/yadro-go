package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/exp/maps"
	"log/slog"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	"yadro-go/pkg/logger"
	"yadro-go/pkg/util"
)

const (
	formantStatementSelectKeywords  = "SELECT * FROM keywords WHERE word IN (%s)"
	statementInsertOrReplaceKeyword = "INSERT OR REPLACE INTO keywords(word, num) VALUES (?, ?)"
)

type KeywordRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewKeywordRepository(log *slog.Logger, db *sql.DB) *KeywordRepository {
	return &KeywordRepository{log: log, db: db}
}

func (r *KeywordRepository) Keywords(ctx context.Context, keywords []string) ([]*domain.ComicKeyword, error) {
	const op = "keyword.Keywords"
	log := r.log.With(slog.String("op", op))

	log.Debug("fetching keywords")

	stmt, err := r.db.PrepareContext(ctx, fmt.Sprintf(formantStatementSelectKeywords, util.GeneratePlaceholders(len(keywords))))
	if err != nil {
		log.Error("failed to prepare statement", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, util.SliceToAny(keywords)...)
	if err != nil {
		log.Error("failed to query keywords", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer rows.Close()

	keywordsMap := make(map[string]*domain.ComicKeyword)

	for rows.Next() {
		var word string
		var num int

		err = rows.Scan(&word, &num)
		if err != nil {
			log.Error("failed to decode keyword", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
		}

		keyword, ok := keywordsMap[word]
		if !ok {
			keyword = &domain.ComicKeyword{Word: word}
			keywordsMap[word] = keyword
		}

		keyword.Nums = append(keyword.Nums, num)
	}

	if err = rows.Err(); err != nil {
		log.Error("error during rows iteration", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	log.Debug("fetch keywords complete")

	return maps.Values(keywordsMap), nil
}

func (r *KeywordRepository) Save(ctx context.Context, keywords []*domain.ComicKeyword) error {
	const op = "keyword.Save"
	log := r.log.With(slog.String("op", op))

	log.Debug("saving keywords")

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to start a transaction", logger.Err(err))
		return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	stmt, err := tx.PrepareContext(ctx, statementInsertOrReplaceKeyword)
	if err != nil {
		log.Error("failed to prepare statement", logger.Err(err))
		return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer stmt.Close()

	for _, keyword := range keywords {
		for _, num := range keyword.Nums {
			_, err = stmt.ExecContext(ctx, keyword.Word, num)
			if err != nil {
				log.Error("failed to execute statement", logger.Err(err))
				if err = tx.Rollback(); err != nil {
					log.Error("tx rollback failed", logger.Err(err))
				}
				return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx commit failed", logger.Err(err))
		return fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	log.Debug("save keywords complete")
	return nil
}

package postgres

import (
	"fmt"

	"github.com/dnahurnyi/proxybot/bot"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type Option func(*Repository)

func New(db *gorm.DB) (*Repository, error) {
	return &Repository{db: db}, nil
}

func (r *Repository) Transaction(action func(repo bot.Repository) error) error {
	tx := r.db.Begin()

	if err := action(&Repository{db: tx}); err != nil {
		err = fmt.Errorf("transaction error: %w", err)
		if rErr := tx.Rollback().Error; rErr != nil {
			err = fmt.Errorf("tranaction rollback: %s: %w", rErr, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("tranaction commit: %w", err)
	}
	return nil
}

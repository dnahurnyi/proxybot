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

func (r *Repository) SaveSubscription(sub *bot.Subscription) error {
	if sub == nil {
		return fmt.Errorf("subscription input is empty")
	}

	if err := r.db.Create(sub).Error; err != nil {
		return fmt.Errorf("create subscription record: %w", err)
	}
	return nil
}

func (r *Repository) GetSubscription(subID int64) (*bot.Subscription, error) {
	if subID == 0 {
		return nil, fmt.Errorf("subscription ID is empty")
	}
	var sub bot.Subscription
	res := r.db.Find(&sub, "id = ?", subID)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("get subscription by channel id: %d", sub.ID)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &sub, nil
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

func (r *Repository) TagSubscription(tag *bot.Tag, subscriptionID int64) error { return nil }
func (r *Repository) CreateTagForChatID(chatID int64, tag string) error        { return nil }
func (r *Repository) CreateChannelForTag(tag string, channelID int64) error    { return nil }
func (r *Repository) ListTags() ([]string, error)                              { return nil, nil }
func (r *Repository) GetChatIDsByTag(tag string) ([]int64, error)              { return nil, nil }
func (r *Repository) ChannelForChat(chatID int64) (int64, error)               { return 0, nil }

package postgres

import (
	"errors"
	"fmt"

	"github.com/dnahurnyi/proxybot/bot"
	"gorm.io/gorm"
)

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
	res := r.db.Preload("Tag").Find(&sub, "id = ?", subID)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get subscription by channel id: %d", subID)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &sub, nil
}

func (r *Repository) ListSubscriptions() ([]bot.Subscription, error) {
	var subs []bot.Subscription
	if err := r.db.Preload("Tag").Find(&subs).Error; err != nil {
		return nil, fmt.Errorf("get all subscriptions: %w", err)
	}
	return subs, nil
}

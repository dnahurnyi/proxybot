package postgres

import (
	"fmt"

	"github.com/dnahurnyi/proxybot/bot"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (r *Repository) SaveTag(tag *bot.Tag) error {
	if tag == nil {
		return fmt.Errorf("tag input is empty")
	}

	if err := r.db.Create(tag).Error; err != nil {
		return fmt.Errorf("create tag record: %w", err)
	}
	return nil
}

func (r *Repository) GetTagByName(name string) (*bot.Tag, error) {
	if name == "" {
		return nil, fmt.Errorf("tag name is empty")
	}
	var tag bot.Tag
	res := r.db.Find(&tag, "name = ?", name)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("get tag by name: %s", name)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &tag, nil
}

func (r *Repository) getTagByID(id uuid.UUID) (*bot.Tag, error) {
	var tag bot.Tag
	res := r.db.Find(&tag, "id = ?", id)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("get tag by id: %s", id)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &tag, nil
}

func (r *Repository) TagSubscription(tagID uuid.UUID, subID int64) error {
	tag, err := r.getTagByID(tagID)
	if err != nil {
		return fmt.Errorf("get tag by id %s: %w", tagID.String(), err)
	}
	if tag == nil {
		return fmt.Errorf("tag with id %s doesn't exist, can't tag", tagID.String())
	}
	sub, err := r.GetSubscription(subID)
	if err != nil {
		return fmt.Errorf("get subscription by id %d: %w", subID, err)
	}
	if sub == nil {
		return fmt.Errorf("subscription with id %d doesn't exist, can't tag it", subID)
	}
	sub.TagID = &tagID
	if err := r.db.Save(sub).Error; err != nil {
		return fmt.Errorf("save subscription record: %w", err)
	}
	return nil
}

func (r *Repository) ListTags() ([]bot.Tag, error) {
	var tags []bot.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("get all tags: %w", err)
	}
	return tags, nil
}

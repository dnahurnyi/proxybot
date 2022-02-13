package bot

import uuid "github.com/satori/go.uuid"

type Tag struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	ChannelID int64     `db:"channel_id"`
}

type Subscription struct {
	ID    int64      `db:"id"`
	Name  string     `db:"name"`
	TagID *uuid.UUID `db:"tag_id"`
	Tag   Tag        `db:"-"`
}

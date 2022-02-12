package bot

import uuid "github.com/satori/go.uuid"

type Tag struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	ChannelID uuid.UUID `db:"channel_id"`
}

type Subscription struct {
	ID    int64      `db:"id"`
	TagID *uuid.UUID `db:"tag_id"`
}

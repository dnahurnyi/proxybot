package bot

import uuid "github.com/satori/go.uuid"

type Client interface {
	MarkAsRead(chatID int64) error
	ForwardMsgToMaster(fromChatID, msgID int64) error
	ForwardMsgTo(fromChatID, msgID, toChatID int64) error
	SubscribeToChannel(channelID int64) error
	MessageToMaster(masterChatID int64, msg string) error
	CreateChannelForTag(tag string) (int64, error)
	GetChannelTitle(channelID int64) (string, error)
	ListChannels() ([]Channel, error)
}

type Repository interface {
	Transaction(func(Repository) error) error
	SaveSubscription(sub *Subscription) error
	GetSubscription(subID int64) (*Subscription, error)
	ListSubscriptions() ([]Subscription, error)
	GetTagByName(tag string) (*Tag, error)
	SaveTag(tag *Tag) error
	TagSubscription(tagID uuid.UUID, subID int64) error
	ListTags() ([]Tag, error)
}

type Channel struct {
	ID   int64
	Name string
}

type Message struct {
	ID              int64
	ChatID          int64
	Content         string
	IsPicture       bool
	IsChannel       bool
	Type            string
	IsPendingStatus bool
	IsForwarded     bool
	ForwardedFromID int64
	UserID          int64
}

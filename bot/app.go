package bot

// App is a package that handles business actions like:
// - new message
//

type Client interface {
	MarkAsRead(chatID int64) error
	ForwardMsgToMaster(fromChatID, msgID int64) error
	ForwardMsgTo(fromChatID, msgID, toChatID int64) error
	SubscribeToChannel(channelID int64) error
	MessageToMaster(masterChatID int64, msg string) error
	CreateChannelForTag(tag string) (int64, error)
	StatsClient
}

type StatsClient interface {
	ListChannels() ([]Channel, error)
}

type Repository interface {
	TagSubscription(tag *Tag, subscriptionID int64) error
	CreateTagForChatID(chatID int64, tag string) error
	CreateChannelForTag(tag string, channelID int64) error
	ListTags() ([]string, error)
	GetChatIDsByTag(tag string) ([]int64, error)
	ChannelForChat(chatID int64) (int64, error)
	SaveSubscription(sub *Subscription) error
	GetSubscription(subID int64) (*Subscription, error)
	Transaction(func(Repository) error) error
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

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
	StatsClient
}

type StatsClient interface {
	ListChannels() ([]Channel, error)
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
}

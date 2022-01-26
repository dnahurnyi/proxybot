package app

// App is a package that handles business actions like:
// - new message
//

type Client interface {
	MarkAsRead(chatID int64) error
	ForwardMsgToMaster(fromChatID, msgID int64) error
	JoinChat(chatID int64) error
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

package client

import (
	"fmt"

	"github.com/zelenin/go-tdlib/client"
)

func (c *Client) MarkAsRead(chatID int64) error {
	_, err := c.tgClient.ViewMessages(&client.ViewMessagesRequest{
		ChatId:    chatID,
		ForceRead: true,
	})
	if err != nil {
		return fmt.Errorf("view messages: %w", err)
	}
	return nil
}

func (c *Client) ForwardMsgToMaster(fromChatID, msgID int64) error {
	_, err := c.tgClient.ForwardMessages(&client.ForwardMessagesRequest{
		ChatId:     c.masterChatID,
		FromChatId: fromChatID,
		MessageIds: []int64{msgID},
	})
	if err != nil {
		return fmt.Errorf("forward message to master: %w", err)
	}
	return nil
}

func (c *Client) JoinChat(chatID int64) error {
	ok, err := c.tgClient.JoinChat(&client.JoinChatRequest{
		ChatId: chatID,
	})
	if err != nil {
		return fmt.Errorf("join chat by chat id: %w", err)
	}
	fmt.Println("Join chat: ", ok, err)
	return nil
}

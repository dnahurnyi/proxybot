package client

import (
	"fmt"
	"math"

	"github.com/dnahurnyi/proxybot/bot"
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
	_, err := c.tgClient.JoinChat(&client.JoinChatRequest{
		ChatId: chatID,
	})
	if err != nil {
		return fmt.Errorf("join chat by chat id: %w", err)
	}
	return nil
}

func (c *Client) SubscribeToChannel(channelID int64) error {
	err := c.joinChat(channelID)
	if err != nil {
		return fmt.Errorf("join channel: %w", err)
	}
	err = c.muteChat(channelID)
	if err != nil {
		return fmt.Errorf("mute channel: %w", err)
	}
	return nil
}

func (c *Client) muteChat(chatID int64) error {
	_, err := c.tgClient.SetChatNotificationSettings(&client.SetChatNotificationSettingsRequest{
		ChatId: chatID,
		NotificationSettings: &client.ChatNotificationSettings{
			MuteFor: math.MaxInt32,
		},
	})
	if err != nil {
		return fmt.Errorf("mute chat by chat id: %w", err)
	}
	return nil
}

func (c *Client) ListChannels() ([]bot.Channel, error) {
	chats, err := c.tgClient.GetChats(&client.GetChatsRequest{
		Limit: math.MaxInt32,
	})
	if err != nil {
		return nil, fmt.Errorf("get chats: %w", err)
	}

	res := []bot.Channel{}
	// TODO: fetch in parallel
	for _, chatID := range chats.ChatIds {
		chat, err := c.tgClient.GetChat(&client.GetChatRequest{
			ChatId: chatID,
		})
		if err != nil {
			return nil, fmt.Errorf("get chat by id %d: %w", chatID, err)
		}
		if chat.Type.ChatTypeType() != client.TypeChatTypeSupergroup {
			// we need only channels
			continue
		}
		res = append(res, bot.Channel{
			ID:   chat.Id,
			Name: chat.Title,
		})
	}

	return res, nil
}

func (c *Client) MessageToMaster(masterChatID int64, msg string) error {
	_, err := c.tgClient.SendMessage(&client.SendMessageRequest{
		ChatId: masterChatID,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: msg,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("send message to master: %w", err)
	}
	return nil
}

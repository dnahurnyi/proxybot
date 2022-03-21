package client

import (
	"fmt"
	"math"
	"strings"

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

func (c *Client) ForwardMsgTo(fromChatID, msgID, toChatID int64) error {
	_, err := c.tgClient.ForwardMessages(&client.ForwardMessagesRequest{
		ChatId:     toChatID,
		FromChatId: fromChatID,
		MessageIds: []int64{msgID},
	})
	if err != nil {
		return fmt.Errorf("forward message to master: %w", err)
	}
	return nil
}

func (c *Client) joinChat(chatID int64) error {
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
		if strings.Contains(err.Error(), "400 CHANNEL_PRIVATE") {
			err = bot.PrivateChannel
		}
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

func (c *Client) CreateChannelForTag(tag string) (int64, error) {
	chat, err := c.tgClient.CreateNewSupergroupChat(&client.CreateNewSupergroupChatRequest{
		Title:       fmt.Sprintf("%s_proxybot", tag),
		IsChannel:   true,
		Description: fmt.Sprintf("Channel for rerouting messages from subscribed channels tagged as %s", tag),
	})
	if err != nil {
		return 0, fmt.Errorf("create channel for tag %s, %w", tag, err)
	}
	// Don't need to create invite with ReplacePrimaryChatInviteLink since SetChatMemberStatus automatically adds user to the channel
	_, err = c.tgClient.SetChatMemberStatus(&client.SetChatMemberStatusRequest{
		ChatId: chat.Id,
		MemberId: &client.MessageSenderUser{
			UserId: c.masterChatID,
		},
		Status: &client.ChatMemberStatusAdministrator{
			CustomTitle:         "master",
			CanBeEdited:         true,
			CanManageChat:       true,
			CanChangeInfo:       true,
			CanPostMessages:     true,
			CanEditMessages:     true,
			CanDeleteMessages:   true,
			CanInviteUsers:      true,
			CanRestrictMembers:  true,
			CanPinMessages:      true,
			CanPromoteMembers:   true,
			CanManageVideoChats: true,
			IsAnonymous:         true,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("make master admin of %s channel, %w", tag, err)
	}
	// TODO: allow reactions in new chat, once it will be implemented in tdlib

	return chat.Id, nil
}

func (c *Client) GetChannelTitle(channelID int64) (string, error) {
	chat, err := c.tgClient.GetChat(&client.GetChatRequest{
		ChatId: channelID,
	})
	if err != nil {
		return "", fmt.Errorf("get chat with id %d: %w", channelID, err)
	}

	return chat.Title, nil
}

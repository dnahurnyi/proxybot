package client

import (
	"fmt"

	"github.com/dnahurnyi/proxybot/bot"
	"github.com/zelenin/go-tdlib/client"
)

type UpdatesHandler interface {
	Handle(msg bot.Message) error
}

func (c *Client) Start(handler UpdatesHandler) error {
	listener := c.tgClient.GetListener()
	defer listener.Close()

	for update := range listener.Updates {
		if update.GetType() == client.TypeUpdateNewMessage {
			msg := update.(*client.UpdateNewMessage)
			err := handler.Handle(transformMsg(msg.Message))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func transformMsg(msg *client.Message) bot.Message {
	chatID := msg.ChatId
	contentText := ""
	if content, ok := msg.Content.(*client.MessageText); ok {
		contentText = content.Text.Text
	}

	res := bot.Message{
		ID:        msg.Id,
		ChatID:    chatID,
		Content:   contentText,
		IsPicture: contentText == "",
		IsChannel: msg.IsChannelPost,
		Type:      "text message",
	}
	if msg.SendingState != nil {
		if msg.SendingState.MessageSendingStateType() == client.TypeMessageSendingStatePending {
			res.IsPendingStatus = true
		}
	}
	if msg.ForwardInfo != nil {
		res.IsForwarded = true
		if msg.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginChannel {
			channelOriginMsg := msg.ForwardInfo.Origin.(*client.MessageForwardOriginChannel)
			fmt.Printf("origin msg chatID: %d\n", channelOriginMsg.ChatId)
			res.ForwardedFromID = channelOriginMsg.ChatId
		}
	}
	if msg.SenderId != nil {
		if msg.SenderId.MessageSenderType() == client.TypeMessageSenderUser {
			sender := msg.SenderId.(*client.MessageSenderUser)
			res.UserID = sender.UserId
		}
	}
	return res
}

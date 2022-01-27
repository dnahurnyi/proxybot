package bot

import (
	"fmt"
)

type UpdatesHandler struct {
	client       Client
	masterChatID int64
}

func NewUpdatesHandler(client Client, masterChatID int64) (*UpdatesHandler, error) {
	return &UpdatesHandler{
		client:       client,
		masterChatID: masterChatID,
	}, nil
}

func (h *UpdatesHandler) Handle(msg Message) error {
	if msg.IsPendingStatus {
		// we don't need to handle pending status notifications
		return nil
	}
	fmt.Printf("----------New Message----------\n")
	fmt.Printf("%s\n", msg.Content)
	if msg.ChatID == h.masterChatID {
		// handle command from master
		err := h.masterCommand(msg)
		if err != nil {
			return fmt.Errorf("handle master command: %w", err)
		}
	} else {
		err := h.client.ForwardMsgToMaster(msg.ChatID, msg.ID)
		if err != nil {
			return fmt.Errorf("forward message to master: %w", err)
		}
	}

	err := h.client.MarkAsRead(msg.ChatID)
	if err != nil {
		return fmt.Errorf("mark as read: %w", err)
	}

	return nil
}

func (h *UpdatesHandler) masterCommand(msg Message) error {
	if msg.IsForwarded {
		// subscribe to channel command
		err := h.client.JoinChat(msg.ForwardedFromID)
		if err != nil {
			return fmt.Errorf("join chat: %w", err)
		}
		err = h.client.MuteChat(msg.ForwardedFromID)
		if err != nil {
			return fmt.Errorf("join chat: %w", err)
		}
	}
	return nil
}

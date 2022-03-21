package bot

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type UpdatesHandler struct {
	client       Client
	repo         Repository
	masterChatID int64
	id           IDGenerator
}

type IDGenerator interface {
	New() uuid.UUID
}

type idGenerator struct{}

func NewIDGenerator() IDGenerator {
	return idGenerator{}
}

func (idGenerator) New() uuid.UUID {
	return uuid.NewV4()
}

func NewUpdatesHandler(client Client, repo Repository, masterChatID int64, idGen IDGenerator) (*UpdatesHandler, error) {
	return &UpdatesHandler{
		client:       client,
		repo:         repo,
		masterChatID: masterChatID,
		id:           idGen,
	}, nil
}

func (h *UpdatesHandler) Handle(msg Message) error {
	if msg.IsPendingStatus {
		// we don't need to handle pending status notifications
		return nil
	}
	fmt.Println("----------New Message----------")
	fmt.Printf("%s\n", msg.Content)
	if msg.ChatID == h.masterChatID {
		// handle command from master
		err := h.MasterCommand(msg)
		if err != nil {
			return fmt.Errorf("handle master command: %w", err)
		}
	} else {
		// reroute to appropriate channel
		err := h.processOuterMessage(msg)
		if err != nil {
			return fmt.Errorf("handle master command: %w", err)
		}
	}

	err := h.client.MarkAsRead(msg.ChatID)
	if err != nil {
		return fmt.Errorf("mark as read: %w", err)
	}

	return nil
}

func (h *UpdatesHandler) MasterCommand(msg Message) error {
	if msg.IsForwarded {
		err := h.saveSubscription(msg.ForwardedFromID)
		if err != nil {
			return fmt.Errorf("save subscription: %w", err)
		}
	}
	// check for commands
	if strings.Contains(msg.Content, "/list_channels") {
		err := h.listSubscriptions()
		if err != nil {
			return fmt.Errorf("list subscriptions: %w", err)
		}
	}
	if strings.Contains(msg.Content, "/tag ") {
		err := h.tagChat(msg.Content)
		if err != nil {
			return fmt.Errorf("tag subscription: %w", err)
		}
	}
	if strings.Contains(msg.Content, "/tags") {
		err := h.listTags()
		if err != nil {
			return fmt.Errorf("list tags: %w", err)
		}
	}
	return nil
}

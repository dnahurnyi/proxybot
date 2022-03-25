package bot

import (
	"errors"
	"fmt"
)

var ErrPrivateChannel = fmt.Errorf("can't subscribe, channel is private")

func (h *UpdatesHandler) saveSubscription(subID int64) error {
	err := h.repo.Transaction(func(repo Repository) (err error) {
		existingSub, err := repo.GetSubscription(subID)
		if err != nil {
			return fmt.Errorf("get subscription by id %d: %w", subID, err)
		}
		if existingSub != nil {
			err = h.client.MessageToMaster(h.masterChatID, fmt.Sprintf("Channel with id %d already subscribed", subID))
			if err != nil {
				return fmt.Errorf("send message to master: %w", err)
			}
			return nil
		}
		subName, err := h.client.GetChannelTitle(subID)
		if err != nil {
			return fmt.Errorf("get subscription name by id %d: %w", subID, err)
		}

		err = repo.SaveSubscription(&Subscription{
			ID:   subID,
			Name: subName,
		})
		if err != nil {
			return fmt.Errorf("repo save subscription: %w", err)
		}

		err = h.client.SubscribeToChannel(subID)
		if err != nil {
			if errors.Is(err, ErrPrivateChannel) {
				errM := h.client.MessageToMaster(h.masterChatID, "Channel is private, can't subscribe")
				if errM != nil {
					return fmt.Errorf("send message to master: %w", errM)
				}
				// return error anyway to cancel transaction
				return err
			}
			return fmt.Errorf("subscribe to channel: %w", err)
		}

		err = h.client.MessageToMaster(h.masterChatID, fmt.Sprintf("Subscribed to %s", subName))
		if err != nil {
			return fmt.Errorf("send message to master: %w", err)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, ErrPrivateChannel) {
			return nil
		}
		return err
	}
	return nil
}

func (h *UpdatesHandler) listSubscriptions() error {
	subs, err := h.repo.ListSubscriptions()
	if err != nil {
		return fmt.Errorf("repository ListSubscriptions: %w", err)
	}
	msg := "Channels that I listen:\n------------------------\n"
	for _, sub := range subs {
		subRow := fmt.Sprintf(" - %s (id: %d)", sub.Name, sub.ID)
		if sub.Tag.Name != "" {
			subRow += fmt.Sprintf(", %s", sub.Tag.Name)
		}
		msg += subRow + "\n"
	}
	err = h.client.MessageToMaster(h.masterChatID, msg)
	if err != nil {
		return fmt.Errorf("send message to master: %w", err)
	}
	return nil
}

func (h *UpdatesHandler) processOuterMessage(msg *Message) error {
	if msg == nil {
		return nil
	}
	sub, err := h.repo.GetSubscription(msg.ChatID)
	if err != nil {
		return fmt.Errorf("get subscription by id %d: %w", msg.ChatID, err)
	}
	if sub != nil && sub.Tag.ChannelID != 0 {
		err = h.client.ForwardMsgTo(msg.ChatID, msg.ID, sub.Tag.ChannelID)
		if err != nil {
			return fmt.Errorf("forward message to tag channel: %w", err)
		}
		return nil
	}
	// message from untagged channel, send to master
	err = h.client.ForwardMsgToMaster(msg.ChatID, msg.ID)
	if err != nil {
		return fmt.Errorf("forward message to tag channel: %w", err)
	}
	return nil
}

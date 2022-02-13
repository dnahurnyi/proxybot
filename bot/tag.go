package bot

import (
	"fmt"
	"strconv"
	"strings"
)

func parseTagCommand(in string) (int64, string, error) {
	commandParts := strings.Split(in, " ")
	if len(commandParts) != 3 {
		return 0, "", fmt.Errorf("tag command should contain chatID and tag name, not: %s", in)
	}
	subIDRaw, err := strconv.Atoi(commandParts[1])
	if err != nil {
		return 0, "", fmt.Errorf("can't convert cahtID to int: %s", commandParts[1])
	}
	return int64(subIDRaw), commandParts[2], nil
}

func (h *UpdatesHandler) tagChat(command string) error {
	subID, tag, err := parseTagCommand(command)
	if err != nil {
		return fmt.Errorf(`parse tag command "%s" to arguments: %w`, command, err)
	}
	err = h.repo.Transaction(func(repo Repository) (err error) {
		// check if chat is subscribed
		existingSub, err := repo.GetSubscription(subID)
		if err != nil {
			return fmt.Errorf("get subscription: %w", err)
		}
		if existingSub == nil {
			err = h.client.MessageToMaster(h.masterChatID, fmt.Sprintf("Can't find subscription by id %d in storage", subID))
			if err != nil {
				return fmt.Errorf("send message to master: %w", err)
			}
			return nil
		}

		existingTag, err := repo.GetTagByName(tag)
		if err != nil {
			return fmt.Errorf("get tag by name %s: %w", tag, err)
		}
		if existingTag == nil {
			channelID, err := h.client.CreateChannelForTag(tag)
			if err != nil {
				return fmt.Errorf("create chat for tag %s: %w", tag, err)
			}
			existingTag = &Tag{
				ID:        h.id.New(),
				Name:      tag,
				ChannelID: channelID,
			}
			err = repo.SaveTag(existingTag)
			if err != nil {
				return fmt.Errorf("save tag %s: %w", tag, err)
			}

			// TODO: make _proxybot sufix constant or env variable
			err = h.client.MessageToMaster(h.masterChatID, fmt.Sprintf("New channel %s_proxybot for your tag  created, you have been invited there as admin", tag))
			if err != nil {
				return fmt.Errorf("send message to master: %w", err)
			}
		}

		err = repo.TagSubscription(existingTag.ID, subID)
		if err != nil {
			return fmt.Errorf("tag subscription %d with tag %s: %w", subID, tag, err)
		}
		// respond to master
		err = h.client.MessageToMaster(h.masterChatID, fmt.Sprintf(`Subscription successfully tagged under %s`, tag))
		if err != nil {
			return fmt.Errorf("send message to master: %w", err)
		}
		return nil
	})
	return err
}

func (h *UpdatesHandler) listTags() error {
	tags, err := h.repo.ListTags()
	if err != nil {
		return fmt.Errorf("can't get list of tags: %w", err)
	}
	var msg string
	if len(tags) == 0 {
		msg = "You have no channels tagged, tag any chat by running `/tag <chatID> <tagName>`"
	} else {
		msg = "List of existing tags: \n"
		for _, tag := range tags {
			msg += fmt.Sprintf(" - %s\n", tag.Name)
		}
	}

	err = h.client.MessageToMaster(h.masterChatID, msg)
	if err != nil {
		return fmt.Errorf("send message to master: %w", err)
	}
	return nil
}

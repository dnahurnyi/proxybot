//go:build unit
// +build unit

package bot_test

import (
	"fmt"
	"testing"

	"github.com/dnahurnyi/proxybot/bot"
	"github.com/dnahurnyi/proxybot/bot/mock"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func uuidToPtrUUID(in uuid.UUID) *uuid.UUID {
	return &in
}

type transaction func(bot.Repository) error

func Test_Handle_outer_message(t *testing.T) {
	t.Parallel()
	masterChatID := int64(25)
	subID := int64(123)
	subName := "subscription"
	tagID := uuid.NewV4()
	tagChannel := int64(234)

	existingSub := &bot.Subscription{
		ID:    subID,
		Name:  subName,
		TagID: &tagID,
		Tag: bot.Tag{
			ID:        tagID,
			Name:      "tag#1",
			ChannelID: tagChannel,
		},
	}

	msg := &bot.Message{
		ChatID:          23,
		ForwardedFromID: subID,
	}

	clientErr := fmt.Errorf("client error")
	repoErr := fmt.Errorf("repo error")

	t.Run("process_outer_message_nil_sub", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().GetSubscription(msg.ChatID).Return(nil, nil)

		client.EXPECT().ForwardMsgToMaster(msg.ChatID, msg.ID).Return(nil)
		client.EXPECT().MarkAsRead(msg.ChatID).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Nil(t, gotErr)
	})

	t.Run("process_outer_message_nil_sub_client_err", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().GetSubscription(msg.ChatID).Return(nil, nil)

		client.EXPECT().ForwardMsgToMaster(msg.ChatID, msg.ID).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Equal(t,
			fmt.Errorf("process outer message: %w",
				fmt.Errorf("forward message to tag channel: %w", clientErr),
			),
			gotErr)
	})

	t.Run("process_outer_message_can't_get_subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().GetSubscription(msg.ChatID).Return(nil, repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Equal(t,
			fmt.Errorf("process outer message: %w",
				fmt.Errorf("get subscription by id %d: %w", msg.ChatID, repoErr),
			),
			gotErr)
	})

	t.Run("process_outer_message_sub_exist", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().GetSubscription(msg.ChatID).Return(existingSub, nil)

		client.EXPECT().ForwardMsgTo(msg.ChatID, msg.ID, tagChannel).Return(nil)
		client.EXPECT().MarkAsRead(msg.ChatID).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Nil(t, gotErr)
	})

	t.Run("process_outer_message_sub_exist_client_err", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().GetSubscription(msg.ChatID).Return(existingSub, nil)

		client.EXPECT().ForwardMsgTo(msg.ChatID, msg.ID, tagChannel).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Equal(t,
			fmt.Errorf("process outer message: %w",
				fmt.Errorf("forward message to tag channel: %w", clientErr),
			),
			gotErr)
	})
}

func Test_Handle_MasterCommand(t *testing.T) {
	t.Parallel()
	masterChatID := int64(25)
	subID := int64(123)

	clientErr := fmt.Errorf("client error")
	repoErr := fmt.Errorf("repo error")

	t.Run("empty_master_command", func(t *testing.T) {
		t.Parallel()

		msg := &bot.Message{
			ChatID:          masterChatID,
			ForwardedFromID: subID,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		client.EXPECT().MarkAsRead(msg.ChatID).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Nil(t, gotErr)
	})

	t.Run("empty_master_command_client_error", func(t *testing.T) {
		t.Parallel()

		msg := &bot.Message{
			ChatID:          masterChatID,
			ForwardedFromID: subID,
			Content:         "/tags",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListTags().Return(nil, repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Equal(t, fmt.Errorf("handle master command: %w",
			fmt.Errorf("list tags: %w",
				fmt.Errorf("repository.ListTags: %w", repoErr),
			),
		), gotErr)
	})

	t.Run("empty_master_command_ack_failure", func(t *testing.T) {
		t.Parallel()

		msg := &bot.Message{
			ChatID:          masterChatID,
			ForwardedFromID: subID,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		client.EXPECT().MarkAsRead(msg.ChatID).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Equal(t, fmt.Errorf("mark as read: %w", clientErr), gotErr)
	})

	t.Run("pending_state_msg", func(t *testing.T) {
		t.Parallel()

		msg := &bot.Message{
			IsPendingStatus: true,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.Handle(msg)
		require.Nil(t, gotErr)
	})
}

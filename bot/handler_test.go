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

func Test_MasterCommand(t *testing.T) {
	t.Parallel()

	var masterChatID int64

	msg := bot.Message{
		ID:      23,
		ChatID:  23,
		Content: "/list_channels",
	}

	subs := []bot.Subscription{
		{
			ID:    1,
			Name:  "channel#1",
			TagID: uuidToPtrUUID(uuid.FromStringOrNil("57545a0c-ef74-413c-88f2-f8d42e5285e6")),
		},
		{
			ID:    2,
			Name:  "channel#2",
			TagID: uuidToPtrUUID(uuid.FromStringOrNil("57545a0c-ef74-413c-88f2-f8d42e5285e6")),
			Tag: bot.Tag{
				ID:        uuid.FromStringOrNil("17545a0c-ef74-413c-88f2-f8d42e5285e6"),
				Name:      "tag#1",
				ChannelID: 23,
			},
		},
	}

	t.Run("list_channels_success", func(t *testing.T) {
		t.Parallel()

		masterMsg := "Channels that I listen:\n------------------------\n" +
			" - channel#1 (id: 1)\n" +
			" - channel#2 (id: 2), tag#1\n"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListSubscriptions().Return(subs, nil)
		client.EXPECT().MessageToMaster(masterChatID, masterMsg).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("repo.ListSubscriptions_fails", func(t *testing.T) {
		t.Parallel()

		repoErr := fmt.Errorf("repo error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListSubscriptions().Return(nil, repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf("list subscriptions: %w", fmt.Errorf("repository ListSubscriptions: %w", repoErr)), gotErr)
	})

	t.Run("repo.ListSubscriptions_fails", func(t *testing.T) {
		t.Parallel()

		masterMsg := "Channels that I listen:\n------------------------\n" +
			" - channel#1 (id: 1)\n" +
			" - channel#2 (id: 2), tag#1\n"

		clientErr := fmt.Errorf("client error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListSubscriptions().Return(subs, nil)
		client.EXPECT().MessageToMaster(masterChatID, masterMsg).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf("list subscriptions: %w", fmt.Errorf("send message to master: %w", clientErr)), gotErr)
	})
}

func uuidToPtrUUID(in uuid.UUID) *uuid.UUID {
	return &in
}

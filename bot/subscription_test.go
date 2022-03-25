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

func Test_MasterCommand_listSubscriptions(t *testing.T) {
	t.Parallel()

	var masterChatID int64

	msg := &bot.Message{
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

	repoErr := fmt.Errorf("repo error")

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

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("repo.ListSubscriptions_fails", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListSubscriptions().Return(nil, repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
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

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf("list subscriptions: %w", fmt.Errorf("send message to master: %w", clientErr)), gotErr)
	})
}

func Test_MasterCommand_subscribe(t *testing.T) {
	t.Parallel()

	var masterChatID int64

	subID := int64(23)
	subName := "sub1"

	msg := &bot.Message{
		IsForwarded:     true,
		ForwardedFromID: subID,
	}

	existingSub := &bot.Subscription{
		ID:    subID,
		Name:  subName,
		TagID: nil,
	}

	repoErr := fmt.Errorf("repo error")
	clientErr := fmt.Errorf("client error")

	t.Run("already_subscribed_to_channel_success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf("Channel with id %d already subscribed", subID)).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("can't_get_subscription_to_check_whether_already_subscribed", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf(
			"save subscription: %w", fmt.Errorf(
				"get subscription by id %d: %w", subID, repoErr,
			),
		), gotErr)
	})

	t.Run("already_subscribed_but_can't_notify_master", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf("Channel with id %d already subscribed", subID)).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf(
			"save subscription: %w", fmt.Errorf(
				"send message to master: %w", clientErr,
			),
		), gotErr)
	})

	t.Run("subscribed_to_the_new_channel_success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)

		client.EXPECT().SubscribeToChannel(subID).Return(nil)
		client.EXPECT().GetChannelTitle(subID).Return(subName, nil)
		repo.EXPECT().SaveSubscription(&bot.Subscription{
			ID:   subID,
			Name: subName,
		}).Return(nil)
		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf("Subscribed to %s", subName)).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("private_channel_success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)
		client.EXPECT().GetChannelTitle(subID).Return(subName, nil)
		repo.EXPECT().SaveSubscription(&bot.Subscription{
			ID:   subID,
			Name: subName,
		}).Return(nil)

		client.EXPECT().SubscribeToChannel(subID).Return(bot.ErrPrivateChannel)
		client.EXPECT().MessageToMaster(masterChatID, "Channel is private, can't subscribe").Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("private_channel_but_can't_notify_master", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)
		client.EXPECT().GetChannelTitle(subID).Return(subName, nil)
		repo.EXPECT().SaveSubscription(&bot.Subscription{
			ID:   subID,
			Name: subName,
		}).Return(nil)

		client.EXPECT().SubscribeToChannel(subID).Return(bot.ErrPrivateChannel)
		client.EXPECT().MessageToMaster(masterChatID, "Channel is private, can't subscribe").Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf(
			"save subscription: %w", fmt.Errorf(
				"send message to master: %w", clientErr,
			),
		), gotErr)
	})

	t.Run("can't_subscribe_to_channel", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)
		client.EXPECT().GetChannelTitle(subID).Return(subName, nil)
		repo.EXPECT().SaveSubscription(&bot.Subscription{
			ID:   subID,
			Name: subName,
		}).Return(nil)

		client.EXPECT().SubscribeToChannel(subID).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf(
			"save subscription: %w", fmt.Errorf(
				"subscribe to channel: %w", clientErr,
			),
		), gotErr)
	})

	t.Run("can't_get_channel_title", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)
		client.EXPECT().GetChannelTitle(subID).Return("", clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf(
			"save subscription: %w", fmt.Errorf(
				"get subscription name by id %d: %w", subID, clientErr,
			),
		), gotErr)
	})

	t.Run("can't_save_subscription_in_DB", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)
		client.EXPECT().GetChannelTitle(subID).Return(subName, nil)
		repo.EXPECT().SaveSubscription(&bot.Subscription{
			ID:   subID,
			Name: subName,
		}).Return(repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf(
			"save subscription: %w", fmt.Errorf(
				"repo save subscription: %w", repoErr,
			),
		), gotErr)
	})

	t.Run("subscribed_to_the_new_channel_but_can't_notify_master", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)

		client.EXPECT().SubscribeToChannel(subID).Return(nil)
		client.EXPECT().GetChannelTitle(subID).Return(subName, nil)
		repo.EXPECT().SaveSubscription(&bot.Subscription{
			ID:   subID,
			Name: subName,
		}).Return(nil)
		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf("Subscribed to %s", subName)).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf(
			"save subscription: %w", fmt.Errorf(
				"send message to master: %w", clientErr,
			),
		), gotErr)
	})
}

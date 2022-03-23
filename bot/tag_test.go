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

func Test_MasterCommand_listTags(t *testing.T) {
	t.Parallel()

	var masterChatID int64

	msg := bot.Message{
		ID:      23,
		ChatID:  23,
		Content: "/tags",
	}

	tags := []bot.Tag{
		{
			ID:        uuid.FromStringOrNil("17545a0c-ef74-413c-88f2-f8d42e5285e6"),
			Name:      "tag#1",
			ChannelID: 23,
		},
		{
			ID:        uuid.FromStringOrNil("27545a0c-ef74-413c-88f2-f8d42e5285e6"),
			Name:      "tag#2",
			ChannelID: 24,
		},
	}

	tagsMsg := "List of existing tags: \n" +
		" - tag#1\n" +
		" - tag#2\n"

	noTagsMsg := "You have no channels tagged, tag any chat by running `/tag <chatID> <tagName>`"
	repoErr := fmt.Errorf("repo error")
	clientErr := fmt.Errorf("client error")

	t.Run("listTags_success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListTags().Return(tags, nil)
		client.EXPECT().MessageToMaster(masterChatID, tagsMsg).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("listTags_no_tags", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListTags().Return([]bot.Tag{}, nil)
		client.EXPECT().MessageToMaster(masterChatID, noTagsMsg).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("listTags_repo_error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListTags().Return(nil, repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf("list tags: %w",
			fmt.Errorf("repository.ListTags: %w", repoErr),
		), gotErr)
	})

	t.Run("listTags_client_error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().ListTags().Return(tags, nil)
		client.EXPECT().MessageToMaster(masterChatID, tagsMsg).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, fmt.Errorf("list tags: %w", fmt.Errorf("send message to master: %w", clientErr)), gotErr)
	})
}

func Test_MasterCommand_tag(t *testing.T) {
	t.Parallel()

	var masterChatID int64

	msg := bot.Message{
		Content: "/tag 24 #tag1",
	}

	subID := int64(24)
	tagName := "#tag1"

	existingSub := &bot.Subscription{
		ID:    subID,
		Name:  "sub1",
		TagID: nil,
	}

	existingTag := &bot.Tag{
		ID:        uuid.FromStringOrNil("27545a0c-ef74-413c-88f2-f8d42e5285e6"),
		Name:      tagName,
		ChannelID: 12,
	}

	newTag := &bot.Tag{
		ID:        uuid.FromStringOrNil("37545a0c-ef74-413c-88f2-f8d42e5285e6"),
		Name:      tagName,
		ChannelID: 13,
	}

	repoErr := fmt.Errorf("repo error")
	clientErr := fmt.Errorf("client error")

	t.Run("tag_with_existing_tag_success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(existingTag, nil)
		repo.EXPECT().TagSubscription(existingTag.ID, subID).Return(nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf(`Subscription successfully tagged under %s`, tagName)).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("tag_with_new_tag_success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)
		idGen := mock.NewMockIDGenerator(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(nil, nil)

		client.EXPECT().CreateChannelForTag(tagName).Return(newTag.ChannelID, nil)
		idGen.EXPECT().New().Return(newTag.ID)
		repo.EXPECT().SaveTag(newTag).Return(nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf(`New channel %s_proxybot for your tag  created, you have been invited there as admin`, tagName)).Return(nil)

		repo.EXPECT().TagSubscription(newTag.ID, subID).Return(nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf(`Subscription successfully tagged under %s`, tagName)).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, idGen)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Nil(t, gotErr)
	})

	t.Run("tag_command_has_no_tag", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		badMsg := "/tag 24"
		gotErr := h.MasterCommand(bot.Message{
			Content: badMsg,
		})
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf(`parse tag command "%s" to arguments: %w`, badMsg,
					fmt.Errorf("tag command should contain chatID and tag name, not: %s", badMsg)),
			), gotErr,
		)
	})

	t.Run("tag_command_has_bad_sub_id", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		badMsg := "/tag badSubID #tag1"
		gotErr := h.MasterCommand(bot.Message{
			Content: badMsg,
		})
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf(`parse tag command "%s" to arguments: %w`, badMsg,
					fmt.Errorf("can't convert cahtID to int: %s", "badSubID")),
			), gotErr,
		)
	})

	t.Run("can't_get_subscriptions", func(t *testing.T) {
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
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("get subscription: %w",
					repoErr,
				),
			), gotErr,
		)
	})

	t.Run("subscription_is_not_in_repo", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf("Can't find subscription by id %d in storage", subID)).Return(nil)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t, nil, gotErr)
	})

	t.Run("subscription_is_not_in_repo_but_can't_notify_master", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(nil, nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf("Can't find subscription by id %d in storage", subID)).
			Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("send message to master: %w",
					clientErr,
				),
			), gotErr,
		)
	})

	t.Run("can't_get_tag", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(nil, repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("get tag by name %s: %w", tagName,
					repoErr,
				),
			), gotErr,
		)
	})

	t.Run("can't_create_channel_for_new_tag", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)
		idGen := mock.NewMockIDGenerator(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(nil, nil)

		client.EXPECT().CreateChannelForTag(tagName).Return(int64(0), clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, idGen)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("create channel for tag %s: %w", tagName,
					fmt.Errorf("create chat for tag %s: %w", tagName,
						clientErr,
					),
				),
			), gotErr,
		)
	})

	t.Run("can't_save_new_tag", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)
		idGen := mock.NewMockIDGenerator(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(nil, nil)

		client.EXPECT().CreateChannelForTag(tagName).Return(newTag.ChannelID, nil)
		idGen.EXPECT().New().Return(newTag.ID)
		repo.EXPECT().SaveTag(newTag).Return(repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, idGen)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("create channel for tag %s: %w", tagName,
					fmt.Errorf("save tag %s: %w", tagName, repoErr),
				),
			), gotErr,
		)
	})

	t.Run("create_new_tag_but_can't_notify_master", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)
		idGen := mock.NewMockIDGenerator(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(nil, nil)

		client.EXPECT().CreateChannelForTag(tagName).Return(newTag.ChannelID, nil)
		idGen.EXPECT().New().Return(newTag.ID)
		repo.EXPECT().SaveTag(newTag).Return(nil)
		client.EXPECT().
			MessageToMaster(masterChatID,
				fmt.Sprintf(`New channel %s_proxybot for your tag  created, you have been invited there as admin`, tagName),
			).
			Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, idGen)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("create channel for tag %s: %w", tagName,
					fmt.Errorf("send message to master: %w", clientErr),
				),
			), gotErr,
		)
	})

	t.Run("can't_tag_subscription_because_of_repo_err", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(existingTag, nil)
		repo.EXPECT().TagSubscription(existingTag.ID, subID).Return(repoErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("tag subscription %d with tag %s: %w", subID, tagName, repoErr),
			), gotErr,
		)
	})

	t.Run("tag_with_existing_tag_but_can't_notify_master", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		client := mock.NewMockClient(ctrl)

		repo.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f transaction) error { return f(repo) })
		repo.EXPECT().GetSubscription(subID).Return(existingSub, nil)
		repo.EXPECT().GetTagByName(tagName).Return(existingTag, nil)
		repo.EXPECT().TagSubscription(existingTag.ID, subID).Return(nil)

		client.EXPECT().MessageToMaster(masterChatID, fmt.Sprintf(`Subscription successfully tagged under %s`, tagName)).Return(clientErr)

		h, err := bot.NewUpdatesHandler(client, repo, masterChatID, nil)
		require.Nil(t, err)
		gotErr := h.MasterCommand(msg)
		require.Equal(t,
			fmt.Errorf("tag subscription: %w",
				fmt.Errorf("send message to master: %w", clientErr),
			), gotErr,
		)
	})
}

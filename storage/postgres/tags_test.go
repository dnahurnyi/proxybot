//go:build integration
// +build integration

package postgres_test

import (
	"fmt"
	"testing"

	"github.com/dnahurnyi/proxybot/bot"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

var (
	tagID1   = uuid.FromStringOrNil("691f06fd-863c-49fb-9972-5d07fcb2ea81")
	tagID2   = uuid.FromStringOrNil("691f06fd-863c-49fb-9972-5d07fcb2ea82")
	tagID3   = uuid.FromStringOrNil("691f06fd-863c-49fb-9972-5d07fcb2ea83")
	seedTags = []bot.Tag{
		{
			ID:        tagID1,
			Name:      "tag#1",
			ChannelID: 121,
		},
		{
			ID:        tagID2,
			Name:      "tag#2",
			ChannelID: 122,
		},
		{
			ID:        tagID3,
			Name:      "tag#3",
			ChannelID: 123,
		},
	}
)

func (s RepoTestSuite) seedTags(tags []bot.Tag) {
	for _, tag := range tags {
		err := s.repo.SaveTag(&tag)
		if err != nil {
			s.T().Fatal(err)
		}
	}
}

func (s *RepoTestSuite) Test_ListTags() {
	s.seedTags(seedTags)

	s.T().Run("list-seeded-tags", func(t *testing.T) {
		got, err := s.repo.ListTags()
		require.NoError(t, err)
		require.ElementsMatch(t, seedTags, got)
	})
	s.TearDownTest()
}

func (s *RepoTestSuite) Test_GetTagByName() {
	s.seedTags(seedTags)

	tests := []struct {
		name string
		in   string
		want *bot.Tag
		err  error
	}{
		{
			name: "name-is-empty",
			in:   "",
			want: nil,
			err:  fmt.Errorf("tag name is empty"),
		},
		{
			name: "tag-doesn't-exist",
			in:   "unknown-tag",
			want: nil,
			err:  nil,
		},
		{
			name: "success",
			in:   "tag#1",
			want: &seedTags[0],
			err:  nil,
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			got, err := s.repo.GetTagByName(tc.in)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.err, err)
		})
	}

	s.TearDownTest()
}

func (s *RepoTestSuite) Test_TagSubscription() {
	s.seedTags(seedTags)
	s.seedSubscriptions(seedSubs)

	tests := []struct {
		name            string
		tagID           uuid.UUID
		subID           int64
		taggingExpected bool
		err             error
	}{
		{
			name:  "tag-doesn't-exist",
			tagID: uuid.FromStringOrNil("7334b195-612a-4954-8431-384d25da5dad"),
			subID: 1121,
			err:   fmt.Errorf("tag with id 7334b195-612a-4954-8431-384d25da5dad doesn't exist, can't tag"),
		},
		{
			name:  "subscription-doesn't-exist",
			tagID: tagID1,
			subID: -1,
			err:   fmt.Errorf("subscription with id -1 doesn't exist, can't tag it"),
		},
		{
			name:  "successful-tagging",
			tagID: tagID1,
			subID: 1121,
			err:   nil,
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.repo.TagSubscription(tc.tagID, tc.subID)
			require.Equal(t, tc.err, err)
			if tc.taggingExpected {
				got, err := s.repo.GetSubscription(tc.subID)
				require.Equal(t, tc.tagID, got.TagID)
				require.Equal(t, tc.err, err)
			}
		})
	}

	s.TearDownTest()
}

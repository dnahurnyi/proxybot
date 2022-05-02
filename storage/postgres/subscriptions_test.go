//go:build integration
// +build integration

package postgres_test

import (
	"fmt"
	"testing"

	"github.com/dnahurnyi/proxybot/bot"
	"github.com/stretchr/testify/require"
)

var (
	seedSubs = []bot.Subscription{
		{
			ID:    1121,
			Name:  "sub#1",
			TagID: nil,
		},
		{
			ID:    1122,
			Name:  "sub#2",
			TagID: nil,
		},
		{
			ID:    1123,
			Name:  "sub#3",
			TagID: nil,
		},
	}
)

func (s RepoTestSuite) seedSubscriptions(subs []bot.Subscription) {
	for _, sub := range subs {
		err := s.repo.SaveSubscription(&sub)
		if err != nil {
			s.T().Fatal(err)
		}
	}
}

func (s *RepoTestSuite) Test_ListSubscriptions() {
	s.seedSubscriptions(seedSubs)

	s.T().Run("list-seeded-subscriptions", func(t *testing.T) {
		got, err := s.repo.ListSubscriptions()
		require.NoError(t, err)
		require.ElementsMatch(t, seedSubs, got)
	})
	s.TearDownTest()
}

func (s *RepoTestSuite) Test_GetSubscription() {
	s.seedSubscriptions(seedSubs)

	tests := []struct {
		name string
		in   int64
		want *bot.Subscription
		err  error
	}{
		{
			name: "id-is-empty",
			in:   0,
			want: nil,
			err:  fmt.Errorf("subscription ID is empty"),
		},
		{
			name: "subsccription-doesn't-exist",
			in:   -1,
			want: nil,
			err:  nil,
		},
		{
			name: "success",
			in:   1121,
			want: &seedSubs[0],
			err:  nil,
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			got, err := s.repo.GetSubscription(tc.in)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.err, err)
		})
	}

	s.TearDownTest()
}

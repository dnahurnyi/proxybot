package client

import (
	"testing"

	"github.com/dnahurnyi/proxybot/bot"
	"github.com/stretchr/testify/require"
	"github.com/zelenin/go-tdlib/client"
)

func Test_transformMsg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *client.Message
		want *bot.Message
	}{
		{
			name: "values",
			in: &client.Message{
				Id:     23,
				ChatId: 24,
				Content: &client.MessageText{
					Text: &client.FormattedText{
						Text: "content",
					},
				},
				SendingState: &client.MessageSendingStatePending{},
				ForwardInfo: &client.MessageForwardInfo{
					Origin: &client.MessageForwardOriginChannel{
						ChatId: 25,
					},
				},
				SenderId: &client.MessageSenderUser{
					UserId: 26,
				},
				IsChannelPost: true,
			},
			want: &bot.Message{
				ID:              23,
				ChatID:          24,
				Content:         "content",
				IsPicture:       false,
				IsChannel:       true,
				Type:            "text message",
				ForwardedFromID: 25,
				UserID:          26,
				IsPendingStatus: true,
				IsForwarded:     true,
			},
		},
		{
			name: "empty",
			in:   nil,
			want: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := transformMsg(tc.in)
			require.Equal(t, tc.want, got)
		})
	}
}

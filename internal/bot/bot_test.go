package bot

import (
	"context"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	bot := New()
	bot.Respond("ping", pong)
	bot.Respond("remember (.+) is (.+)", remember)
	bot.Respond("say (.+)", say)

	var cases = []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "Ping pong test",
			in:       "ping",
			expected: "pong",
		},
		{
			name:     "Simple sentence",
			in:       "remember abc is 123",
			expected: "123",
		},
		{
			name:     "Simple sentence with uppercase",
			in:       "reMemBer abc IS 123",
			expected: "123",
		},
		{
			name:     "Bot answer with a constructed sentence",
			in:       "say coucou",
			expected: "Oh okay, coucou",
		},
	}

	ctx := context.TODO()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			answer := bot.Sentence(ctx, tc.in)
			if tc.expected != answer {
				t.Errorf("got %s, want %s", answer, tc.expected)
			}
		})
	}

}

func pong(msg Message) string {
	return "pong"
}

func remember(msg Message) string {
	return msg.Matches[1]
}

func say(msg Message) string {
	return fmt.Sprintf("Oh okay, %s", msg.Matches[0])
}

package domain_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/mackerel-server-go/domain"
)

func TestChannel(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"email", "*domain.EmailChannel"},
		{"slack", "*domain.SlackChannel"},
		{"webhook", "*domain.WebhookChannel"},
		{"other", "*domain.Channel"},
	}

	for _, c := range cases {
		ch := &domain.Channel{
			Type: c.in,
		}

		got := fmt.Sprintf("%T", ch.Cast())
		if got != c.want {
			t.Fail()
		}
	}
}

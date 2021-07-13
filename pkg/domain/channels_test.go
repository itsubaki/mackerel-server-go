package domain_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

func TestChannel(t *testing.T) {
	cases := []struct {
		t string
		T string
	}{
		{"email", "*domain.EmailChannel"},
		{"slack", "*domain.SlackChannel"},
		{"webhook", "*domain.WebhookChannel"},
		{"other", "*domain.Channel"},
	}

	for _, c := range cases {
		ch := &domain.Channel{
			Type: c.t,
		}

		if fmt.Sprintf("%T", ch.Cast()) != c.T {
			t.Error()
		}
	}
}

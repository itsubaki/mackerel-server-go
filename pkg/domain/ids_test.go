package domain_test

import (
	"regexp"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

func TestNewRandomID(t *testing.T) {
	id := domain.NewRandomID()
	if len(id) != 11 {
		t.Error()
	}

	if !regexp.MustCompile(`[0-9A-Za-z]`).Match([]byte(id)) {
		t.Errorf("id=%s doesn't match [0-9A-Za-z]", id)
	}

}

func TestIDWith(t *testing.T) {
	regexp := regexp.MustCompile(`[0-9A-Za-z]`)
	cases := []struct {
		length int
		seed   []string
	}{
		{11, []string{"hoge", "hoge"}},
		{11, []string{"foobar"}},
	}

	for _, c := range cases {
		id := domain.NewIDWith(c.seed...)
		if len(id) != c.length {
			t.Error()
		}

		if !regexp.Match([]byte(id)) {
			t.Errorf("id=%s doesn't match [0-9A-Za-z]", id)
		}
	}
}

func TestNewID(t *testing.T) {
	regexp := regexp.MustCompile(`[0-9A-Za-z]`)
	cases := []struct {
		length int
		seed   []string
	}{
		{2, []string{"hoge"}},
		{3, []string{"foobar"}},
		{5, []string{"piyo"}},
		{8, []string{"fuga"}},
	}

	for _, c := range cases {
		id := domain.NewID(c.length, c.seed...)
		if len(id) != c.length {
			t.Errorf("id=%s, len(id)=%d", id, len(id))
		}

		if !regexp.Match([]byte(id)) {
			t.Errorf("id=%s doesn't match [0-9A-Za-z]", id)
		}
	}
}

package domain

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

type Org struct {
	Name string `json:"name"`
}

type XAPIKey struct {
	Org     string
	Name    string
	XAPIKey string
	Read    bool
	Write   bool
}

func NewXAPIKey(org, name string, write bool) *XAPIKey {
	sha := sha256.Sum256([]byte(uuid.Must(uuid.NewRandom()).String()))
	hash := hex.EncodeToString(sha[:])

	return &XAPIKey{
		Org:     org,
		Name:    name,
		XAPIKey: hash[:44],
		Read:    true,
		Write:   write,
	}
}

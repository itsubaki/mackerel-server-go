package domain

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func NewOrgID() string {
	return NewRandomID(11)
}

func NewHostID() string {
	return NewRandomID(11)
}

func NewAPIKey() string {
	return NewRandomID(44)
}

func NewMonitorID(seed ...string) string {
	return NewID(11, seed...)
}

func NewAlertID(seed ...string) string {
	return NewID(11, seed...)
}

func NewRandomID(digit int) string {
	return NewID(digit, uuid.Must(uuid.NewRandom()).String())
}

func NewID(digit int, seed ...string) string {
	var sum string
	for i := range seed {
		sum = sum + seed[i]
	}

	sha := sha256.Sum256([]byte(sum))
	hash := hex.EncodeToString(sha[:])

	return hash[:digit]
}

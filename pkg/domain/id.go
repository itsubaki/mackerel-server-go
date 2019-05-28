package domain

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func NewXAPIKey() string {
	return GenID(44, uuid.Must(uuid.NewRandom()).String())
}

func NewOrgID() string {
	return GenID(11, uuid.Must(uuid.NewRandom()).String())
}

func NewHostID() string {
	return GenID(11, uuid.Must(uuid.NewRandom()).String())
}

func NewMonitorID(seed ...string) string {
	return GenID(11, seed...)
}

func NewAlertID(seed ...string) string {
	return GenID(11, seed...)
}

func GenID(digit int, seed ...string) string {
	var sum string
	for i := range seed {
		sum = sum + seed[i]
	}

	sha := sha256.Sum256([]byte(sum))
	hash := hex.EncodeToString(sha[:])

	return hash[:digit]
}

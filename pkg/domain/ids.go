package domain

import (
	"github.com/google/uuid"
	"github.com/speps/go-hashids"
)

func NewHostID() string {
	return NewRandomID(11)
}

func NewAPIKey() string {
	return NewRandomID(44)
}

func NewIDWith(seed ...string) string {
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

	hd := hashids.NewData()
	hd.MinLength = digit
	hd.Salt = sum
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	id, err := h.Encode([]int{42})
	if err != nil {
		panic(err)
	}

	return id
}

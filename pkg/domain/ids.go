package domain

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/speps/go-hashids"
)

func NewRandomID() string {
	rand.Seed(time.Now().UnixNano())
	return NewID(11, strconv.Itoa(rand.Int()))
}

func NewIDWith(seed ...string) string {
	return NewID(11, seed...)
}

func NewID(digit int, seed ...string) string {
	var salt string
	for i := range seed {
		salt = salt + seed[i]
	}

	hd := hashids.NewData()
	hd.MinLength = digit
	hd.Salt = salt

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

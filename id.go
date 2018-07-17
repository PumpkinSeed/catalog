package catalog

import (
	"math/rand"
	"strconv"
	"time"
)

type identifier uint

func NewID() identifier {
	return identifier(random())
}

func NewIDFromString(id string) (identifier, error) {
	uintID, err := strconv.ParseUint(id, 10, 10)
	return identifier(uintID), err
}

func (i *identifier) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

func random() uint {
	nano := uint64(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	return uint(rand.Uint64() + nano)
}

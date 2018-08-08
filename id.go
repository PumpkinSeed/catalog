package catalog

import (
	"math/rand"
	"strconv"
	"time"
)

type Identifier uint

func NewID() Identifier {
	return Identifier(random())
}

func NewIDFromString(id string) (Identifier, error) {
	uintID, err := strconv.ParseUint(id, 10, 10)
	return Identifier(uintID), err
}

func (i *Identifier) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

func random() uint {
	nano := uint64(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	return uint(rand.Uint64() + nano)
}

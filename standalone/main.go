package main

import (
	"sync"

	"github.com/PumpkinSeed/catalog"
)

func main() {
	serv := catalog.NewServer("127.0.0.1:7777", nil, &sync.RWMutex{})
	panic(serv.Listen())
}

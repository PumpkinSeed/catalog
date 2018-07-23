package main

import "github.com/PumpkinSeed/catalog"

func main() {
	serv := catalog.NewServer("127.0.0.1:8080")
	panic(serv.Listen())
}

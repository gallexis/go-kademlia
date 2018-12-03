package main

import (
	"fmt"
	"kademlia/datastructure"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	n := datastructure.NewNodeID()
	fmt.Println(n)
}

package main

import (
    "fmt"
    "math/rand"
    "time"
    "kademlia/datastructure"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    dht := NewDHT()
    dht.eventDispatcher.Start()
    dht.Bootstrap(dht.bootstrapNodes[1])
    dht.Receiver()
    dht.PopulateRT()

    time.Sleep(time.Second * 10)

    dht.GetPeers(datastructure.NewNodeIdFromString("4EBF7D54EABA7380D46C05604B059FABAEA212F0"))

    fmt.Scanln()
}

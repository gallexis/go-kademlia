package main

import (
    "fmt"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

/*
    gerer les données entrantes (data + IP/port)
    gerer la structure des nodes (update)
 */

func main() {
    dht := NewDHT()
    dht.eventDispatcher.Start()
    dht.Bootstrap(dht.bootstrapNodes[1])
    dht.Receiver()
    dht.PopulateRT()

    //time.Sleep(time.Second * 10)

    //dht.GetPeers(datastructure.NewNodeIdFromString("4EBF7D54EABA7380D46C05604B059FABAEA212F0"))



    fmt.Scanln()
}

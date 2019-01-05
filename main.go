package main

import (
    "fmt"
    "math/rand"
    "time"
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

    //n := datastructure.NewNodeIdFromString("b2d76aa3bd8c3b8b755b29ed0d95b2ef65ae44b4")
    //c := datastructure.Contact{Port: 21456, IP: net.ParseIP("2.154.8.94"), NodeID: n}
    //
    //tx1 := message.NewTransactionId()
    //dht.Send(message.PingRequest{}.Encode(tx1, dht.selfNodeID), c)
    //log.Info("Sent ping: ", tx1.String(), " to: ", c.IP, c.Port)
    //
    //tx2 := message.NewTransactionId()
    //dht.Send(message.GetPeersRequest{}.Encode(tx2, dht.selfNodeID, n), c)
    //log.Info("Sent request: ", tx2.String(), " to: ", c.IP, c.Port)

    fmt.Scanln()
}

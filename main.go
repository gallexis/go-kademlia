package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "kademlia/datastructure"
    "kademlia/message"
    "math/rand"
    "net"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    dht := NewDHT()
    dht.Bootstrap(dht.bootstrapNodes[1])
    go dht.Receiver()
    //dht.PopulateRT()

    n := datastructure.NewNodeIDFromString("b2d76aa3bd8c3b8b755b29ed0d95b2ef65ae44b4")
    c := datastructure.Contact{Port: 21456, IP: net.ParseIP("2.154.8.94"), NodeID: n}

    tx1 := message.NewRandomBytes(2)
    dht.Send(message.PingRequest{}.Encode(tx1, dht.nodeID), c)
    log.Info("Sent ping: ", tx1.String(), " to: ", c.IP, c.Port)

    tx2 := message.NewRandomBytes(2)
    dht.Send(message.GetPeersRequest{}.Encode(tx2, dht.nodeID, n), c)
    log.Info("Sent request: ", tx2.String(), " to: ", c.IP, c.Port)

    fmt.Scanln()
}

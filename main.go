package main

import (
    "kademlia/message"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    dht := NewDHT()
    go dht.Receiver()
    dht.Bootstrap(dht.bootstrapNodes[0])

    dht.Send(message.PingRequest{}.Encode(message.NewRandomBytes(2), dht.nodeID), "74.69.68.188", 40107)
    dht.Send(message.PingRequest{}.Encode(message.NewRandomBytes(2), dht.nodeID), "177.16.122.194", 22440)

    dht.routingTable.Display()
}

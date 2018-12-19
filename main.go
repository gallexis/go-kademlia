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
    dht.Bootstrap(dht.bootstrapNodes[0])
    go dht.Receiver()
    dht.PopulateRT()

    dht.routingTable.Display()
    fmt.Scanln()
}

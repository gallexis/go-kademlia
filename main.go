package main

import (
    "kademlia/datastructure"
    "math/rand"
    "time"
    //log "github.com/sirupsen/logrus"
)

func init() {
    //log.
    rand.Seed(time.Now().UTC().UnixNano())
}

/*
    add statistics
    clean code
    write tests & comments
 */

func main() {
    dht := NewDHT()
    dht.eventDispatcher.Start()
    go dht.CallbackCaller()
    dht.Bootstrap(dht.bootstrapNodes[1])
    dht.Receiver()
    dht.Timer()
    dht.PopulateRT()

    time.Sleep(time.Second * 5)

    //dht.GetPeers(datastructure.NewNodeIdFromString("57537D93A76F574369DC2E573E99C3840A9FD89D"))
    //dht.GetPeers(datastructure.NewNodeIdFromString("FC0CCE628DBE7EEA0CF655A6A13336791021F25F"))
    //dht.GetPeers(datastructure.NewNodeIdFromString("23E7A4876B36CE427A847A827306B4B2DC67304A"))
    //dht.GetPeers(datastructure.NewNodeIdFromString("5EF929E35650741627DACA28E18A3DF0FC5A53DB"))
    dht.GetPeers(datastructure.NewNodeIdFromString("D58952BDBBBFBA9DA444F8FE99DCF2C7F2E4AB77"))
    //dht.GetPeers(datastructure.NewNodeIdFromString("4EBF7D54EABA7380D46C05604B059FABAEA212F0"))

    dht.PendingPingPool()
}

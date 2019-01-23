package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "math/rand"
    "time"
)

func init() {
    log.SetLevel(log.DebugLevel)
    rand.Seed(time.Now().UTC().UnixNano())
}

/*
    add statistics
    clean code
    write tests & comments
    deal with utorrent & bittorrent
    set max tries for getpeers queries
 */

func main() {
    dht := NewDHT()
    err := dht.Init()

    if err != nil{
        log.Error("DHT failure : ", err.Error())
        return
    }

    fmt.Scanln()
}

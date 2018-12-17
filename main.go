package main

import (
    "fmt"
    ds "kademlia/datastructure"
    "kademlia/message"
    "math/rand"
    "net"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    n := ds.NewNodeID()
    n2 := ds.NewNodeID()

    v := []string{
        "router.utorrent.com:6881",
        "router.bittorrent.com:6881",
    }

    raddr, err := net.ResolveUDPAddr("udp", v[1])
    if err != nil {
        fmt.Println("can't resolve")
        return
    }

    buffer := make([]byte, 600)

    conn, err := net.DialUDP("udp", nil, raddr)
    if err != nil {
        fmt.Println("can't dial")
    }
    defer conn.Close()

    deadline := time.Now().Add(time.Second * 10)
    err = conn.SetReadDeadline(deadline)
    if err != nil {
        fmt.Println("too long")
        return
    }

    conn.Write(message.FindNodeRequest{}.Encode(message.NewRandomBytes(2), n, n2))
    //conn.Write(message.PingRequest{}.Encode(message.NewRandomBytes(2), n))
    //conn.Write(message.GetPeersRequest{}.Encode(message.NewRandomBytes(2), n, n2))

    _, _, err = conn.ReadFrom(buffer)
    if err != nil {
        fmt.Println("can't read", err.Error())
        return
    }

    fmt.Println(string(buffer), "\n")

    //dht := NewDHT(ds.NewNodeID())

    Router(buffer)
}

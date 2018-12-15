package main

import (
    "fmt"
    "kademlia/core/handler"
    "kademlia/datastructure"
    "kademlia/network"
    "kademlia/network/krpc"
    "net"
    "time"
)

func init() {

}

func main() {
    pingChan := make(chan krpc.PingResponse)
    findNodesChan := make(chan krpc.FindNodeResponse)
    getPeersChan := make(chan krpc.GetPeersResponse)
    getPeersWithNodesChan := make(chan krpc.GetPeersResponseWithNodes)
    announcePeerChan := make(chan krpc.AnnouncePeersResponse)

    go handler.OnPing(pingChan)
    go handler.OnFindNodes(findNodesChan)
    go handler.OnGetPeers(getPeersChan)
    go handler.OnGetPeersWithNodes(getPeersWithNodesChan)
    go handler.OnAnnouncePeer(announcePeerChan)

    n := datastructure.NewNodeID()
    n2 := datastructure.NewNodeID()

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

    conn.Write(krpc.FindNodeRequest{}.Encode(krpc.NewRandomBytes(2), n, n2))
    conn.Write(krpc.PingRequest{}.Encode(krpc.NewRandomBytes(2), n))
    conn.Write(krpc.GetPeersRequest{}.Encode(krpc.NewRandomBytes(2), n, n2))

    _, _, err = conn.ReadFrom(buffer)
    if err != nil {
        fmt.Println("can't read", err.Error())
        return
    }

    fmt.Println(string(buffer), "\n")
    network.Router(buffer)

    _, _, err = conn.ReadFrom(buffer)
    if err != nil {
        fmt.Println("can't read", err.Error())
        return
    }

    fmt.Println(string(buffer), "\n")
    network.Router(buffer)

    _, _, err = conn.ReadFrom(buffer)
    if err != nil {
        fmt.Println("can't read", err.Error())
        return
    }

    fmt.Println(string(buffer), "\n")
    network.Router(buffer)

}


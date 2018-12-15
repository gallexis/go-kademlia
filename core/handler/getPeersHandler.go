package handler

import (
    "fmt"
    "kademlia/network/krpc"
)

func OnGetPeers(getPeersChan chan krpc.GetPeersResponse) {
    for {
        select {
        case getPeers := <-getPeersChan:
            fmt.Println(getPeers)
        }
    }
}

func OnGetPeersWithNodes(getPeersWithNodesChan chan krpc.GetPeersResponseWithNodes) {
    for {
        select {
        case getPeersWithNodes := <-getPeersWithNodesChan:
            fmt.Println(getPeersWithNodes)
        }
    }
}

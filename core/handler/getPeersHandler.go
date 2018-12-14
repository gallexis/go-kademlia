package handler

import (
    "fmt"
    "kademlia/network/messages"
)

func OnGetPeers(getPeersChan chan messages.GetPeersResponse) {
    for {
        select {
        case getPeers := <-getPeersChan:
            fmt.Println(getPeers)
        }
    }
}

func OnGetPeersWithNodes(getPeersWithNodesChan chan messages.GetPeersResponseWithNodes) {
    for {
        select {
        case getPeersWithNodes := <-getPeersWithNodesChan:
            fmt.Println(getPeersWithNodes)
        }
    }
}

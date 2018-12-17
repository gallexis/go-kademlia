package main

import (
    "fmt"
    "kademlia/message"
    ds "kademlia/datastructure"
)

type DHT struct {
    RoutingTable ds.RoutingTable
    //peerstore
}

func NewDHT(id ds.NodeID) DHT {
    return DHT{
        RoutingTable: ds.NewRoutingTable(id),
    }
}

func OnAnnouncePeerResponse(announcePeer message.AnnouncePeersResponse) {
    fmt.Println(announcePeer)
}

func OnFindNodesResponse(findNodes message.FindNodeResponse) {
    fmt.Printf("findNodes: %+v \n", findNodes)
}

func OnGetPeersResponse(getPeers message.GetPeersResponse) {
    fmt.Printf("getPeers: %+v \n", getPeers)
}

func OnGetPeersWithNodesResponse(getPeersWithNodes message.GetPeersResponseWithNodes) {
    fmt.Printf("getPeersWithNodes: %+v \n", getPeersWithNodes)
}

func OnPingResponse(ping message.PingResponse) {
    fmt.Printf("ping: %+v \n", ping)
}

func Router(data []byte) {
    g := message.BytesToMessage(data)
    switch g.Y {
    case "q":

        switch g.Q {
        case "ping":
            ping := message.PingRequest{}
            ping.Decode(g.T, g.A.Id)

        case "find_node":
            findNodesRequest := message.FindNodeRequest{}
            findNodesRequest.Decode(g.T, g.A)

        case "get_peers":
            getPeers := message.GetPeersRequest{}
            getPeers.Decode(g.T, g.A)

        case "announce_peer":
            announcePeers := message.AnnouncePeersRequest{}
            announcePeers.Decode(g.T, g.A)

        default:
            panic("q")
        }

    case "r":

        switch {
        case len(g.R.Values) > 0 && g.R.Token != "":
            getPeers := message.GetPeersResponse{}
            getPeers.Decode(g.T, g.R)
            OnGetPeersResponse(getPeers)

        case len(g.R.Nodes) > 0 && g.R.Token != "":
            getPeers := message.GetPeersResponseWithNodes{}
            getPeers.Decode(g.T, g.R)
            OnGetPeersWithNodesResponse(getPeers)

        case len(g.R.Nodes) > 0:
            findNodes := message.FindNodeResponse{}
            findNodes.Decode(g.T, g.R)
            OnFindNodesResponse(findNodes)

        case len(g.R.Id) > 0:
            ping := message.PingResponse{}
            ping.Decode(g.T, g.R.Id)
            OnPingResponse(ping)

            /*
            AnnouncePeersResponse == PingResponse

                announcePeers := AnnouncePeersResponse{}
                announcePeers.Decode(g.RandomBytes, g.R)
                fmt.Println(announcePeers)
             */

        default:
            panic("r")
        }

    case "e":
        fmt.Println("Error:", g.E)

    }
}

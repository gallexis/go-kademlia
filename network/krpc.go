package network

import (
    "fmt"
    "kademlia/network/krpc"
)

func Router(data []byte) {
    g := krpc.BytesToMessage(data)
    switch g.Y {
    case "q":

        switch g.Q {
        case "ping":
            ping := krpc.PingRequest{}
            ping.Decode(g.T, g.A.Id)

        case "find_node":
            findNodesRequest := krpc.FindNodeRequest{}
            findNodesRequest.Decode(g.T, g.A)

        case "get_peers":
            getPeers := krpc.GetPeersRequest{}
            getPeers.Decode(g.T, g.A)

        case "announce_peer":
            announcePeers := krpc.AnnouncePeersRequest{}
            announcePeers.Decode(g.T, g.A)

        default:
            panic("q")
        }

    case "r":

        switch {
        case len(g.R.Values) > 0 && g.R.Token !=  "": //&& g.R.Id != "":
            getPeers := krpc.GetPeersResponse{}
            getPeers.Decode(g.T, g.R)
            fmt.Printf("getPeers1: %+v \n", getPeers)

        case len(g.R.Nodes) > 0 && g.R.Token !=  "": //&& g.R.Id != "":
            getPeers := krpc.GetPeersResponseWithNodes{}
            getPeers.Decode(g.T, g.R)
            fmt.Printf("getPeers2: %+v \n", getPeers)

        case len(g.R.Nodes) > 0: // && g.R.Id != "":
            findNodes := krpc.FindNodeResponse{}
            findNodes.Decode(g.T, g.R)
            fmt.Printf("findNodes: %+v \n", findNodes)

        case len(g.R.Id) > 0:
            ping := krpc.PingResponse{}
            ping.Decode(g.T, g.R.Id)
            fmt.Printf("ping: %+v \n", ping)

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

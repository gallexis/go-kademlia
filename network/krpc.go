package network

import (
    "bytes"
    "fmt"
    "kademlia/network/messages"
)

func Router(data []byte) {
    g := messages.BytesToMessage(data)
    switch g.Y {
    case "q":

        switch g.Q {
        case "ping":
            ping := messages.PingRequest{}
            ping.Decode(g.T, g.A.Id)

        case "find_node":
            findNodesRequest := messages.FindNodeRequest{}
            findNodesRequest.Decode(g.T, g.A)

        case "get_peers":
            getPeers := messages.GetPeersRequest{}
            getPeers.Decode(g.T, g.A)

        case "announce_peer":
            announcePeers := messages.AnnouncePeersRequest{}
            announcePeers.Decode(g.T, g.A)

        default:
            panic("q")
        }

    case "r":

        switch {
        case len(g.R.Values) > 0 && !bytes.Equal(g.R.Token, []byte{}): //&& g.R.Id != "":
            getPeers := messages.GetPeersResponse{}
            getPeers.Decode(g.T, g.R)
            fmt.Printf("getPeers1: %+v \n", getPeers)

        case len(g.R.Nodes) > 0 && !bytes.Equal(g.R.Token, []byte{}): //&& g.R.Id != "":
            getPeers := messages.GetPeersResponseWithNodes{}
            getPeers.Decode(g.T, g.R)
            fmt.Printf("getPeers2: %+v \n", getPeers)

        case len(g.R.Nodes) > 0: // && g.R.Id != "":
            findNodes := messages.FindNodeResponse{}
            findNodes.Decode(g.T, g.R)
            fmt.Printf("findNodes: %+v \n", findNodes)

        case len(g.R.Id) > 0:
            ping := messages.PingResponse{}
            ping.Decode(g.T, g.R.Id)
            fmt.Printf("ping: %+v \n", ping)

            /*
            AnnouncePeersResponse == PingResponse

                announcePeers := AnnouncePeersResponse{}
                announcePeers.Decode(g.T, g.R)
                fmt.Println(announcePeers)
             */

        default:
            panic("r")
        }

    case "e":
        fmt.Println("Error:", g.E)

    }
}

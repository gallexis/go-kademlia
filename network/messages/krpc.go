package messages

import (
    "fmt"
    "github.com/ehmry/go-bencode"
    "log"
)

// GENERIC RECEIVE
type Response struct {
    Values []string `bencode:"values"`
    Id     string   `bencode:"id"`
    Nodes  string   `bencode:"nodes"`
    Nodes6 string   `bencode:"nodes6"`
    Token  string   `bencode:"token"`
}

type Answer struct {
    Id          string `bencode:"id"`
    Target      string `bencode:"target"`
    InfoHash    string `bencode:"info_hash"`
    Port        int    `bencode:"port"`
    Token       string `bencode:"token"`
    ImpliedPort int    `bencode:"implied_port"`
}

type GenericMessage struct {
    T string   `bencode:"t"`
    Y string   `bencode:"y"`
    Q string   `bencode:"q"`
    R Response `bencode:"r"`
    A Answer   `bencode:"a"`
    E []string `bencode:"e"`
}

// GENERIC SEND
type RequestMessage struct {
    T string                 `bencode:"t"`
    Y string                 `bencode:"y"`
    Q string                 `bencode:"q"`
    A map[string]interface{} `bencode:"a"`
}

type ResponseMessage struct {
    T string                 `bencode:"t"`
    Y string                 `bencode:"y"`
    R map[string]interface{} `bencode:"r"`
}

func (g GenericMessage) Dispatch(r []byte) {
    if err := bencode.Unmarshal(r, &g); err != nil {
        log.Fatalln(err.Error())
    }

    switch g.Y {
    case "r":

        switch g.Q {
        case "ping":
            ping := PingResponse{}
            ping.Decode(g.T, g.R.Id)

        case "find_node":
            findNodes := FindNodeResponse{}
            findNodes.Decode(g.T, g.R)

        case "get_peers":
            if len(g.R.Values) > 0 {
                getPeers := GetPeersResponse{}
                getPeers.Decode(g.T, g.R)

            } else if g.R.Nodes != "" {
                getPeers := GetPeersResponseWithNodes{}
                getPeers.Decode(g.T, g.R)
            }
        case "announce_peer":
            announcePeers := AnnouncePeersResponse{}
            announcePeers.Decode(g.T, g.R)

        default:
            panic("")
        }

    case "q":

        switch {
        case g.A.Target != "" && g.A.Id != "":
            findNodesRequest := FindNodeRequest{}
            findNodesRequest.Decode(g.T, g.A)

        case g.A.InfoHash != "" && g.A.Token != "" && g.A.Id != "":
            announcePeers := AnnouncePeersRequest{}
            announcePeers.Decode(g.T, g.A)

        case g.A.InfoHash != "" && g.A.Id != "":
            getPeers := GetPeersRequest{}
            getPeers.Decode(g.T, g.A)

        case g.A.Id != "":
            ping := PingRequest{}
            ping.Decode(g.T, g.A.Id)

        default:
            panic("")
        }

    case "e":
        fmt.Println("Error:", g.E)

    }

}

package messages

import (
    "github.com/ehmry/go-bencode"
    "kademlia/datastructure"
    "log"
)

type announcePeers struct {
    T  string
    Id datastructure.NodeID
}

type AnnouncePeersResponse struct {
    announcePeers
}

func (g *AnnouncePeersResponse) Decode(t string, r Response) {
    g.T = t
    g.Id = datastructure.StringToNodeID(r.Id)
}

func (_ AnnouncePeersResponse) Encode(t string, id datastructure.NodeID) []byte {
    q := ResponseMessage{}
    q.T = t
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": id.String(),
    }
    buffer, err := bencode.Marshal(q)

    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

type AnnouncePeersRequest struct {
    announcePeers
    ImpliedPort int
    InfoHash    datastructure.NodeID
    Port        int
    Token       string
}

func (g *AnnouncePeersRequest) Decode(t string, a Answer) {
    g.T = t
    g.Id = datastructure.StringToNodeID(a.Id)
    g.ImpliedPort = a.ImpliedPort
    g.InfoHash = datastructure.StringToNodeID(a.InfoHash)
    g.Port = a.Port
    g.Token = a.Token
}

func (_ AnnouncePeersRequest) Encode(t, token string, id, infoHash datastructure.NodeID, impliedPort, port int) []byte {
    q := RequestMessage{}
    q.T = t
    q.Y = "q"
    q.A = map[string]interface{}{
        "id":           id.String(),
        "implied_port": impliedPort,
        "info_hash":    infoHash.String(),
        "port":         port,
        "token":        token,
    }

    buffer, err := bencode.Marshal(q)
    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

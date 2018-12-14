package messages

import (
    "kademlia/datastructure"
)

type announcePeers struct {
    T  Token
    Id datastructure.NodeID
}

type AnnouncePeersResponse struct {
    announcePeers
}

func (g *AnnouncePeersResponse) Decode(t string, r Response) {
    g.T = NewTokenFromString(t)
    g.Id = datastructure.BytesToNodeID(r.Id)
}

func (_ AnnouncePeersResponse) Encode(t Token, id datastructure.NodeID) []byte {
    q := ResponseMessage{}
    q.T = t
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": id.Bytes(),
    }

    return MessageToBytes(q)
}

type AnnouncePeersRequest struct {
    announcePeers
    ImpliedPort int
    InfoHash    datastructure.NodeID
    Port        int
    Token       string
}

func (g *AnnouncePeersRequest) Decode(t string, a Answer) {
    g.T = NewTokenFromString(t)
    g.Id = datastructure.BytesToNodeID(a.Id)
    g.ImpliedPort = a.ImpliedPort
    g.InfoHash = datastructure.BytesToNodeID(a.InfoHash)
    g.Port = a.Port
    g.Token = string(a.Token)
}

func (_ AnnouncePeersRequest) Encode(t, token Token, id, infoHash datastructure.NodeID, impliedPort, port int) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.A = map[string]interface{}{
        "id":           id.Bytes(),
        "implied_port": impliedPort,
        "info_hash":    infoHash.Bytes(),
        "port":         port,
        "token":        token.String(),
    }

    return MessageToBytes(q)
}

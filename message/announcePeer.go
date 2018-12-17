package message

import ds "kademlia/datastructure"

type announcePeers struct {
    T  RandomBytes
    Id ds.NodeID
}

type AnnouncePeersResponse struct {
    announcePeers
}

func (g *AnnouncePeersResponse) Decode(t string, r Response) {
    g.T = NewRandomBytesFromString(t)
    g.Id.Decode(r.Id)
}

func (_ AnnouncePeersResponse) Encode(t RandomBytes, id ds.NodeID) []byte {
    q := ResponseMessage{}
    q.T = t.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": id.Encode(),
    }

    return MessageToBytes(q)
}

type AnnouncePeersRequest struct {
    announcePeers
    ImpliedPort int
    InfoHash    ds.NodeID
    Port        int
    Token       RandomBytes
}

func (g *AnnouncePeersRequest) Decode(t string, a Answer) {
    g.T = NewRandomBytesFromString(t)
    g.Id.Decode(a.Id)
    g.ImpliedPort = a.ImpliedPort
    g.InfoHash.Decode(a.InfoHash)
    g.Port = a.Port
    g.Token = NewRandomBytesFromString(a.Token)
}

func (_ AnnouncePeersRequest) Encode(t, token RandomBytes, id, infoHash ds.NodeID, impliedPort, port int) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.A = map[string]interface{}{
        "id":           id.Encode(),
        "implied_port": impliedPort,
        "info_hash":    infoHash.Encode(),
        "port":         port,
        "token":        token.String(),
    }

    return MessageToBytes(q)
}

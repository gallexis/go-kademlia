package messages

import (
    "github.com/ehmry/go-bencode"
    "kademlia/datastructure"
    "log"
)

type getPeers struct {
    T string
}

type GetPeersResponse struct {
    getPeers
    Id     datastructure.NodeID
    Token  string
    Values []string
}

func (g *GetPeersResponse) Decode(t string, r Response) {
    g.T = t
    g.Id = datastructure.StringToNodeID(r.Id)
    g.Token = r.Token
    g.Values = r.Values
}

func (_ GetPeersResponse) Encode(t string, id datastructure.NodeID, token string, values []string) []byte {
    q := ResponseMessage{}
    q.T = t
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":     id.String(),
        "token":  token,
        "values": values,
    }

    buffer, err := bencode.Marshal(q)
    if err != nil {
        log.Println(err.Error())
    }
    return buffer
}

type GetPeersResponseWithNodes struct {
    getPeers
    Id    datastructure.NodeID
    Token string
    Nodes []datastructure.NodeID
}

func (g *GetPeersResponseWithNodes) Decode(t string, r Response) {
    numberOfNodes := len(r.Nodes) / 40
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }
    g.T = t
    g.Id = datastructure.StringToNodeID(r.Id)
    g.Token = r.Token
    for i := 0; i < numberOfNodes; i++ {
        offset := i * 40
        g.Nodes = append(g.Nodes, datastructure.StringToNodeID(r.Nodes[offset:(offset + 40)]))
    }
}

func (_ GetPeersResponseWithNodes) Encode(t string, id datastructure.NodeID, token string, nodes []datastructure.NodeID) []byte {
    nodesToString := ""
    numberOfNodes := len(nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }
    for i := 0; i < numberOfNodes; i++ {
        nodesToString += nodes[i].String()
    }

    q := ResponseMessage{}
    q.T = t
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":    id.String(),
        "token": token,
        "nodes": nodesToString,
    }

    buffer, err := bencode.Marshal(q)
    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

type GetPeersRequest struct {
    findNode
    Id       datastructure.NodeID
    InfoHash datastructure.NodeID
}

func (g *GetPeersRequest) Decode(t string, a Answer) {
    g.T = t
    g.Id = datastructure.StringToNodeID(a.Id)
    g.InfoHash = datastructure.StringToNodeID(a.InfoHash)
}

func (_ GetPeersRequest) Encode(t string, id, infoHash datastructure.NodeID) []byte {
    q := RequestMessage{}
    q.T = t
    q.Y = "q"
    q.Q = "get_peers"
    q.A = map[string]interface{}{
        "id":        id.String(),
        "info_hash": infoHash.String(),
    }
    buffer, err := bencode.Marshal(q)

    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

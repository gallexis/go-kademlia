package messages

import (
    "kademlia/datastructure"
)

type getPeers struct {
    T Token
}

type GetPeersResponse struct {
    getPeers
    Id     datastructure.NodeID
    Token  Token
    Values []string
}

func (g *GetPeersResponse) Decode(t string, r Response) {
    g.T = NewTokenFromString(t)
    g.Id = datastructure.BytesToNodeID(r.Id)
    g.Token = r.Token
    g.Values = r.Values
}

func (_ GetPeersResponse) Encode(t Token, id datastructure.NodeID, token Token, values []string) []byte {
    q := ResponseMessage{}
    q.T = t
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":     id.Bytes(),
        "token":  token,
        "values": values,
    }

    return MessageToBytes(q)
}

type GetPeersResponseWithNodes struct {
    getPeers
    Id    datastructure.NodeID
    Token Token
    Nodes []datastructure.NodeID
}

func (g *GetPeersResponseWithNodes) Decode(t string, r Response) {
    lengthNodeID := datastructure.BytesInNodeID
    numberOfNodes := len(r.Nodes) / lengthNodeID
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }
    g.T = NewTokenFromString(t)
    g.Id = datastructure.BytesToNodeID(r.Id)
    g.Token = r.Token
    for i := 0; i < numberOfNodes; i++ {
        offset := i * lengthNodeID
        g.Nodes = append(g.Nodes, datastructure.BytesToNodeID(r.Nodes[offset:(offset + lengthNodeID)]))
    }
}

func (_ GetPeersResponseWithNodes) Encode(t Token, id datastructure.NodeID, token Token, nodes []datastructure.NodeID) []byte {
    var byteNodes []byte
    numberOfNodes := len(nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    for i := 0; i < numberOfNodes; i++ {
        byteNodes = append(byteNodes, nodes[i].Bytes()...)
    }

    q := ResponseMessage{}
    q.T = []byte(t)
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":    id.Bytes(),
        "token": token,
        "nodes": byteNodes,
    }

    return MessageToBytes(q)
}

type GetPeersRequest struct {
    findNode
    Id       datastructure.NodeID
    InfoHash datastructure.NodeID
}

func (g *GetPeersRequest) Decode(t string, a Answer) {
    g.T = NewTokenFromString(t)
    g.Id = datastructure.BytesToNodeID(a.Id)
    g.InfoHash = datastructure.BytesToNodeID(a.InfoHash)
}

func (_ GetPeersRequest) Encode(t Token, id, infoHash datastructure.NodeID) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.Q = "get_peers"
    q.A = map[string]interface{}{
        "id":        id.Bytes(),
        "info_hash": infoHash.Bytes(),
    }

    return MessageToBytes(q)
}

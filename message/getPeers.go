package message

import (
    ds "kademlia/datastructure"
)

type GetPeersResponse struct {
    T TransactionId
    Id     ds.NodeId
    Token  Token
    Values []string
}

func (g *GetPeersResponse) Decode(t string, r Response) {
    g.T = NewTransactionIdFromString(t)
    g.Id.Decode(r.Id)
    g.Token = r.Token
    g.Values = r.Values
}

func (g GetPeersResponse) Encode() []byte {
    q := ResponseMessage{}
    q.T = g.T.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":     g.Id.Encode(),
        "token":  g.Token.String(),
        "values": g.Values,
    }

    return MessageToBytes(q)
}

type GetPeersResponseWithNodes struct {
    T TransactionId
    Id    ds.NodeId
    Token Token
    Nodes []ds.NodeId
}

func (g *GetPeersResponseWithNodes) Decode(t string, r Response) {
    lengthNodeID := ds.BytesInNodeID
    numberOfNodes := len(r.Nodes) / lengthNodeID
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }
    g.T = NewTransactionIdFromString(t)
    g.Id.Decode(r.Id)
    g.Token = r.Token
    for i := 0; i < numberOfNodes; i++ {
        offset := i * lengthNodeID
        nid := ds.NodeId{}
        nid.Decode(r.Nodes[offset:(offset + lengthNodeID)])
        g.Nodes = append(g.Nodes, nid)
    }
}

func (g GetPeersResponseWithNodes) Encode() []byte {
    var byteNodes []byte
    numberOfNodes := len(g.Nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    for i := 0; i < numberOfNodes; i++ {
        byteNodes = append(byteNodes, g.Nodes[i].Encode()...)
    }

    q := ResponseMessage{}
    q.T = g.T.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":    g.Id.Encode(),
        "token": g.Token.String(),
        "nodes": byteNodes,
    }

    return MessageToBytes(q)
}

type GetPeersRequest struct {
    T TransactionId
    Id       ds.NodeId
    InfoHash ds.InfoHash
}

func (g *GetPeersRequest) Decode(t string, a Answer) {
    g.T = NewTransactionIdFromString(t)
    g.Id.Decode(a.Id)
    g.InfoHash.Decode(a.InfoHash)
}

func (g GetPeersRequest) Encode() []byte {
    q := RequestMessage{}
    q.T = g.T.String()
    q.Y = "q"
    q.Q = "get_peers"
    q.A = map[string]interface{}{
        "id":        g.Id.Encode(),
        "info_hash": g.InfoHash.Encode(),
    }

    return MessageToBytes(q)
}

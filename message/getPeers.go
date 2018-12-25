package message

import (
    ds "kademlia/datastructure"
)

type getPeers struct {
    T RandomBytes
}

type GetPeersResponse struct {
    getPeers
    Id     ds.NodeID
    Token  RandomBytes
    Values []string
}

func (g *GetPeersResponse) Decode(t string, r Response) {
    g.T = NewRandomBytesFromString(t)
    g.Id.Decode(r.Id)
    g.Token = r.Token
    g.Values = r.Values
}

func (_ GetPeersResponse) Encode(t RandomBytes, id ds.NodeID, token RandomBytes, values []string) []byte {
    q := ResponseMessage{}
    q.T = t.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":     id.Encode(),
        "token":  token.String(),
        "values": values,
    }

    return MessageToBytes(q)
}

type GetPeersResponseWithNodes struct {
    getPeers
    Id    ds.NodeID
    Token RandomBytes
    Nodes []ds.NodeID
}

func (g *GetPeersResponseWithNodes) Decode(t string, r Response) {
    lengthNodeID := ds.BytesInNodeID
    numberOfNodes := len(r.Nodes) / lengthNodeID
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }
    g.T = NewRandomBytesFromString(t)
    g.Id.Decode(r.Id)
    g.Token = r.Token
    for i := 0; i < numberOfNodes; i++ {
        offset := i * lengthNodeID
        nid := ds.NodeID{}
        nid.Decode(r.Nodes[offset:(offset + lengthNodeID)])
        g.Nodes = append(g.Nodes, nid)
    }
}

func (_ GetPeersResponseWithNodes) Encode(t RandomBytes, id ds.NodeID, token RandomBytes, nodes []ds.NodeID) []byte {
    var byteNodes []byte
    numberOfNodes := len(nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    for i := 0; i < numberOfNodes; i++ {
        byteNodes = append(byteNodes, nodes[i].Encode()...)
    }

    q := ResponseMessage{}
    q.T = t.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id":    id.Encode(),
        "token": token.String(),
        "nodes": byteNodes,
    }

    return MessageToBytes(q)
}

type GetPeersRequest struct {
    findNode
    Id       ds.NodeID
    InfoHash ds.NodeID
}

func (g *GetPeersRequest) Decode(t string, a Answer) {
    g.T = NewRandomBytesFromString(t)
    g.Id.Decode(a.Id)
    g.InfoHash.Decode(a.InfoHash)
}

func (_ GetPeersRequest) Encode(t RandomBytes, id, infoHash ds.NodeID) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.Q = "get_peers"
    q.A = map[string]interface{}{
        "id":        id.Encode(),
        "info_hash": infoHash.Encode(),
    }

    return MessageToBytes(q)
}

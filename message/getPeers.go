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

func (g *GetPeersResponse) Decode(message GenericMessage) {
    g.T = NewTransactionIdFromString(message.T)
    g.Id.Decode(message.R.Id)
    g.Token = message.R.Token
    g.Values = message.R.Values
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
    Nodes []ds.Node
}

func (g *GetPeersResponseWithNodes) Decode(message GenericMessage) {
    lengthNodeID := ds.BytesInNodeID
    numberOfNodes := len(message.R.Nodes) / lengthNodeID
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }
    g.T = NewTransactionIdFromString(message.T)
    g.Id.Decode(message.R.Id)
    g.Token = message.R.Token

    for i := 0; i < numberOfNodes; i++ {
        offset := i * lengthNodeID
        node := ds.Node{}
        node.ContactInfo.Decode(message.R.Nodes[offset:(offset + lengthNodeID)])
        g.Nodes = append(g.Nodes, node)
    }
}

func (g GetPeersResponseWithNodes) Encode() []byte {
    var byteNodes []byte
    numberOfNodes := len(g.Nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    for i := 0; i < numberOfNodes; i++ {
        byteNodes = append(byteNodes, g.Nodes[i].ContactInfo.Encode()...)
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

func (g *GetPeersRequest) Decode(message GenericMessage) {
    g.T = NewTransactionIdFromString(message.T)
    g.Id.Decode(message.A.Id)
    g.InfoHash.Decode(message.A.InfoHash)
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

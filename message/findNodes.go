package message

import (
    ds "kademlia/datastructure"
)

type findNode struct {
    T  RandomBytes
    Id ds.NodeID
}

type FindNodeResponse struct {
    findNode
    Nodes []ds.Contact
}

func (e *FindNodeResponse) Decode(t string, response Response) {
    lengthNodeID := 26
    numberOfNodes := len(response.Nodes) / lengthNodeID
    if numberOfNodes > 16 {
        numberOfNodes = 16
    }

    e.T = NewRandomBytesFromString(t)
    e.Id.Decode(response.Id)
    for i := 0; i < numberOfNodes; i++ {
        offset := i * lengthNodeID
        contact := ds.Contact{}
        contact.Decode(response.Nodes[offset:(offset + lengthNodeID)])
        e.Nodes = append(e.Nodes, contact)

    }
}

func (_ FindNodeResponse) Encode(t RandomBytes, id ds.NodeID, nodes []ds.NodeID) []byte {
    numberOfNodes := len(nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    q := ResponseMessage{}
    var byteNodes []byte
    q.T = t.String()
    q.Y = "r"

    for i := 0; i < numberOfNodes; i++ {
        byteNodes = append(byteNodes, nodes[i].Encode()...)
    }

    q.R = map[string]interface{}{
        "id":    id.Encode(),
        "nodes": byteNodes,
    }

    return MessageToBytes(q)
}

type FindNodeRequest struct {
    findNode
    Target ds.NodeID
}

func (e *FindNodeRequest) Decode(t string, a Answer) {
    e.T = NewRandomBytesFromString(t)
    e.Id.Decode(a.Id)
    e.Target.Decode(a.Target)
}

func (_ FindNodeRequest) Encode(t RandomBytes, id, target ds.NodeID) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.Q = "find_node"

    q.A = map[string]interface{}{
        "id":     id.Encode(),
        "target": target.Encode(),
    }

    return MessageToBytes(q)
}

package krpc

import (
    "kademlia/datastructure"
)

type findNode struct {
    T  RandomBytes
    Id datastructure.NodeID
}

type FindNodeResponse struct {
    findNode
    Nodes []datastructure.NodeID
}

func (e *FindNodeResponse) Decode(t string, response Response) {
    lengthNodeID := datastructure.BytesInNodeID
    numberOfNodes := len(response.Nodes) / lengthNodeID
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    e.T = NewRandomBytesFromString(t)
    e.Id = datastructure.BytesToNodeID(response.Id)
    for i := 0; i < numberOfNodes; i++ {
        offset := i * lengthNodeID
        e.Nodes = append(e.Nodes, datastructure.BytesToNodeID(response.Nodes[offset:(offset + lengthNodeID)]))
    }
}

func (_ FindNodeResponse) Encode(t RandomBytes, id datastructure.NodeID, nodes []datastructure.NodeID) []byte {
    numberOfNodes := len(nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    q := ResponseMessage{}
    var byteNodes []byte
    q.T = t.String()
    q.Y = "r"

    for i := 0; i < numberOfNodes; i++ {
        byteNodes = append(byteNodes, nodes[i].Bytes()...)
    }

    q.R = map[string]interface{}{
        "id":    id.Bytes(),
        "nodes": byteNodes,
    }

    return MessageToBytes(q)
}

type FindNodeRequest struct {
    findNode
    Target datastructure.NodeID
}

func (e *FindNodeRequest) Decode(t string, a Answer) {
    e.T = NewRandomBytesFromString(t)
    e.Id = datastructure.BytesToNodeID(a.Id)
    e.Target = datastructure.BytesToNodeID(a.Target)
}

func (_ FindNodeRequest) Encode(t RandomBytes, id, target datastructure.NodeID) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.Q = "find_node"

    q.A = map[string]interface{}{
        "id":     id.Bytes(),
        "target": target.Bytes(),
    }

    return MessageToBytes(q)
}

package message

import (
    ds "kademlia/datastructure"
)

type FindNodeResponse struct {
    T  TransactionId
    Id ds.NodeId
    Nodes []ds.Node
}

func (f *FindNodeResponse) Decode(message GenericMessage) {
    lengthNodeID := 26
    numberOfNodes := len(message.R.Nodes) / lengthNodeID
    if numberOfNodes > 16 {
        numberOfNodes = 16
    }

    f.T = NewTransactionIdFromString(message.T)
    f.Id.Decode(message.R.Id)

    for i := 0; i < numberOfNodes; i++ {
        offset := i * lengthNodeID
        node := ds.Node{}
        node.Decode(message.R.Nodes[offset:(offset + lengthNodeID)])
        f.Nodes = append(f.Nodes, node)
    }
}

func (f FindNodeResponse) Encode() []byte {
    numberOfNodes := len(f.Nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    q := ResponseMessage{}
    var byteNodes []byte
    q.T = f.T.String()
    q.Y = "r"

    for i := 0; i < numberOfNodes; i++ {
        byteNodes = append(byteNodes, f.Nodes[i].Encode()...)
    }

    q.R = map[string]interface{}{
        "id":    f.Id.Encode(),
        "nodes": byteNodes,
    }

    return MessageToBytes(q)
}

type FindNodeRequest struct {
    T  TransactionId
    Id ds.NodeId
    Target ds.NodeId
}

func (f *FindNodeRequest) Decode(message GenericMessage) {
    f.T = NewTransactionIdFromString(message.T)
    f.Id.Decode(message.A.Id)
    f.Target.Decode(message.A.Target)
}

func (f FindNodeRequest) Encode() []byte {
    q := RequestMessage{}
    q.T = f.T.String()
    q.Y = "q"
    q.Q = "find_node"

    q.A = map[string]interface{}{
        "id":     f.Id.Encode(),
        "target": f.Target.Encode(),
    }

    return MessageToBytes(q)
}

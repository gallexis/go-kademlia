package messages

import (
    "github.com/ehmry/go-bencode"
    "kademlia/datastructure"
    "log"
)

type findNode struct {
    T  string
    Id datastructure.NodeID
}

type FindNodeResponse struct {
    findNode
    Nodes []datastructure.NodeID
}

func (e *FindNodeResponse) Decode(t string, response Response) {
    numberOfNodes := len(response.Nodes) / 40
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    e.Id = datastructure.StringToNodeID(response.Id)
    for i := 0; i < numberOfNodes; i++ {
        start := i * 40
        e.Nodes = append(e.Nodes, datastructure.StringToNodeID(response.Nodes[start:(start + 40)]))

    }
}

func (_ FindNodeResponse) Encode(t string, id datastructure.NodeID, nodes []datastructure.NodeID) []byte {
    numberOfNodes := len(nodes)
    if numberOfNodes > 8 {
        numberOfNodes = 8
    }

    q := ResponseMessage{}
    nodesToString := ""
    q.T = t
    q.Y = "r"

    for i := 0; i < numberOfNodes; i++ {
        nodesToString += nodes[i].String()
    }

    q.R = map[string]interface{}{
        "id":    id.String(),
        "nodes": nodesToString,
    }

    buffer, err := bencode.Marshal(q)
    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

type FindNodeRequest struct {
    findNode
    Target datastructure.NodeID
}

func (e *FindNodeRequest) Decode(t string, a Answer) {
    e.T = t
    e.Id = datastructure.StringToNodeID(a.Id)
    e.Target = datastructure.StringToNodeID(a.Target)
}

func (_ FindNodeRequest) Encode(t string, id, target datastructure.NodeID) []byte {
    q := RequestMessage{}
    q.T = t
    q.Y = "q"
    q.Q = "find_node"
    q.A = map[string]interface{}{
        "id":     id.String(),
        "target": target.String(),
    }
    buffer, err := bencode.Marshal(q)

    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

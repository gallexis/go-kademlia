package messages

import (
    "github.com/ehmry/go-bencode"
    "kademlia/datastructure"
    "log"
)

type ping struct {
    T  string
    Id datastructure.NodeID
}

type PingRequest struct {
    ping
}

func (_ PingRequest) Encode(t string, nodeID datastructure.NodeID) []byte {
    q := RequestMessage{}
    q.T = t
    q.Y = "q"
    q.Q = "ping"
    q.A = map[string]interface{}{
        "id": nodeID.String(),
    }
    buffer, err := bencode.Marshal(q)

    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

func (p *PingRequest) Decode(t string, nodeID string) {
    p.T = t
    p.Id = datastructure.StringToNodeID(nodeID)
}


type PingResponse struct {
    ping
}

func (p *PingResponse) Decode(t string, nodeID string) {
    p.T = t
    p.Id = datastructure.StringToNodeID(nodeID)
}

func (_ PingResponse) Encode(t string, nodeID datastructure.NodeID) []byte {
    q := ResponseMessage{}
    q.T = t
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": nodeID.String(),
    }

    buffer, _ := bencode.Marshal(q)

    return buffer
}
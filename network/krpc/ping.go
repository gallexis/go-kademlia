package krpc

import (
    "kademlia/datastructure"
)

type ping struct {
    T  RandomBytes
    Id datastructure.NodeID
}

type PingRequest struct {
    ping
}

func (_ PingRequest) Encode(t RandomBytes, nodeID datastructure.NodeID) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.Q = "ping"

    q.A = map[string]interface{}{
        "id": nodeID.Bytes(),
    }

    return MessageToBytes(q)
}

func (p *PingRequest) Decode(t string, nodeID []byte) {
    p.T = NewRandomBytesFromString(t)
    p.Id = datastructure.BytesToNodeID(nodeID)
}


type PingResponse struct {
    ping
}

func (p *PingResponse) Decode(t string, nodeID []byte) {
    p.T = NewRandomBytesFromString(t)
    p.Id = datastructure.BytesToNodeID(nodeID)
}

func (_ PingResponse) Encode(t RandomBytes, nodeID datastructure.NodeID) []byte {
    q := ResponseMessage{}
    q.T = t.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": nodeID.Bytes(),
    }

    return MessageToBytes(q)
}

package message

import (
    ds "kademlia/datastructure"
)

type ping struct {
    T  RandomBytes
    Id ds.NodeID
}

type PingRequest struct {
    ping
}

func (_ PingRequest) Encode(t RandomBytes, nodeID ds.NodeID) []byte {
    q := RequestMessage{}
    q.T = t.String()
    q.Y = "q"
    q.Q = "ping"

    q.A = map[string]interface{}{
        "id": nodeID.Encode(),
    }

    return MessageToBytes(q)
}

func (p *PingRequest) Decode(t string, nodeID []byte) {
    p.T = NewRandomBytesFromString(t)
    p.Id.Decode(nodeID)
}


type PingResponse struct {
    ping
}

func (p *PingResponse) Decode(t string, nodeID []byte) {
    p.T = NewRandomBytesFromString(t)
    p.Id.Decode(nodeID)
}

func (_ PingResponse) Encode(t RandomBytes, nodeID ds.NodeID) []byte {
    q := ResponseMessage{}
    q.T = t.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": nodeID.Encode(),
    }

    return MessageToBytes(q)
}

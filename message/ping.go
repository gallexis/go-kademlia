package message

import (
    ds "kademlia/datastructure"
)

type PingRequest struct {
    T  TransactionId
    Id ds.NodeId
}

func (p *PingRequest) Decode(message GenericMessage) {
    p.T = NewTransactionIdFromString(message.T)
    p.Id.Decode(message.A.Id)
}

func (p PingRequest) Encode() []byte {
    q := RequestMessage{}
    q.T = p.T.String()
    q.Y = "q"
    q.Q = "ping"

    q.A = map[string]interface{}{
        "id": p.Id.Encode(),
    }

    return MessageToBytes(q)
}

type PingResponse struct {
    T  TransactionId
    Id ds.NodeId
}

func (p *PingResponse) Decode(message GenericMessage) {
    p.T = NewTransactionIdFromString(message.T)
    p.Id.Decode(message.R.Id)
}

func (p PingResponse) Encode() []byte {
    q := ResponseMessage{}
    q.T = p.T.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": p.Id.Encode(),
    }

    return MessageToBytes(q)
}

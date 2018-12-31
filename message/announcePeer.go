package message

import ds "kademlia/datastructure"

type AnnouncePeersResponse struct {
    T  TransactionId
    Id ds.NodeId
}

func (a *AnnouncePeersResponse) Decode(t string, r Response) {
    a.T = NewTransactionIdFromString(t)
    a.Id.Decode(r.Id)
}

func (a AnnouncePeersResponse) Encode() []byte {
    q := ResponseMessage{}
    q.T = a.T.String()
    q.Y = "r"
    q.R = map[string]interface{}{
        "id": a.Id.Encode(),
    }

    return MessageToBytes(q)
}

type AnnouncePeersRequest struct {
    T  TransactionId
    Id ds.NodeId
    ImpliedPort int
    InfoHash    ds.NodeId
    Port        int
    Token       Token
}

func (a *AnnouncePeersRequest) Decode(t string, answer Answer) {
    a.T = NewTransactionIdFromString(t)
    a.Id.Decode(answer.Id)
    a.ImpliedPort = answer.ImpliedPort
    a.InfoHash.Decode(answer.InfoHash)
    a.Port = answer.Port
    a.Token = answer.Token
}

func (a AnnouncePeersRequest) Encode() []byte {
    q := RequestMessage{}
    q.T = a.T.String()
    q.Y = "q"
    q.A = map[string]interface{}{
        "id":           a.Id.Encode(),
        "implied_port": a.ImpliedPort,
        "info_hash":    a.InfoHash.Encode(),
        "port":         a.Port,
        "token":        a.Token.String(),
    }

    return MessageToBytes(q)
}

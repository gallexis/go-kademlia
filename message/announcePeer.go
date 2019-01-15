package message

import ds "kademlia/datastructure"

type AnnouncePeersResponse struct {
    T  TransactionId
    Id ds.NodeId
}

func (a *AnnouncePeersResponse) Decode(message GenericMessage) {
    a.T = NewTransactionIdFromString(message.T)
    a.Id.Decode(message.R.Id)
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

func (a *AnnouncePeersRequest) Decode(message GenericMessage) {
    a.T = NewTransactionIdFromString(message.T)
    a.Id.Decode(message.A.Id)
    a.ImpliedPort = message.A.ImpliedPort
    a.InfoHash.Decode(message.A.InfoHash)
    a.Port = message.A.Port
    a.Token = message.A.Token
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

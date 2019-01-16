package test

import (
    "kademlia/datastructure"
    "kademlia/message"
    "testing"
)

func TestGetPeersResponse(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    tx := message.NewTransactionIdFromString("aaeebb")
    token := message.NewTransactionIdFromString("bbaaee")
    values := []string{"abc", "def"}
    encoded := message.GetPeersResponse{}.Encode(tx, randomNodeID, token, values)
    g := message.BytesToMessage(encoded)

    response := message.GetPeersResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        response.T.String() != tx.String() ||
        response.Token.String() != token.String() ||
        response.Peers[0] != values[0] ||
        response.Peers[1] != values[1]{
        t.Error("")
    }
}

func TestGetPeersResponseWithNodes(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    n2 := ds.FakeNodeID(0x13)
    n3 := ds.FakeNodeID(0x14)
    tx := message.NewTransactionIdFromString("aaeebb")
    token := message.NewTransactionIdFromString("bbaaee")
    nodes := []ds.NodeID{n2, n3}
    encoded := message.GetPeersResponseWithNodes{}.Encode(tx, randomNodeID, token, nodes)
    g := message.BytesToMessage(encoded)

    response := message.GetPeersResponseWithNodes{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        response.T.String() != tx.String() ||
        response.Token.String() != token.String() ||
        !response.Nodes[0].Equals(nodes[0])  ||
        !response.Nodes[1].Equals(nodes[1]){
        t.Error("")
    }
}

func TestGetPeersRequest(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    infohash := ds.FakeNodeID(0xF4)
    tx := message.NewTransactionIdFromString("aaeebb")
    encoded := message.GetPeersRequest{}.Encode(tx, randomNodeID, infohash)
    g := message.BytesToMessage(encoded)

    response := message.GetPeersRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.InfoHash.Equals(infohash) ||
        response.T.String() != tx.String(){
        t.Error()
    }
}

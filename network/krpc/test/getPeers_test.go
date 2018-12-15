package test

import (
    "kademlia/datastructure"
    "kademlia/network/krpc"
    "testing"
)

func TestGetPeersResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    token := krpc.NewRandomBytesFromString("bbaaee")
    values := []string{"abc", "def"}
    encoded := krpc.GetPeersResponse{}.Encode(tx, randomNodeID, token, values)
    g := krpc.BytesToMessage(encoded)

    response := krpc.GetPeersResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        response.T.String() != tx.String() ||
        response.Token.String() != token.String() ||
        response.Values[0] != values[0] ||
        response.Values[1] != values[1]{
        t.Error("")
    }
}

func TestGetPeersResponseWithNodes(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    n2 := datastructure.FakeNodeID(0x13)
    n3 := datastructure.FakeNodeID(0x14)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    token := krpc.NewRandomBytesFromString("bbaaee")
    nodes := []datastructure.NodeID{n2, n3}
    encoded := krpc.GetPeersResponseWithNodes{}.Encode(tx, randomNodeID, token, nodes)
    g := krpc.BytesToMessage(encoded)

    response := krpc.GetPeersResponseWithNodes{}
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
    randomNodeID := datastructure.FakeNodeID(0x12)
    infohash := datastructure.FakeNodeID(0xF4)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    encoded := krpc.GetPeersRequest{}.Encode(tx, randomNodeID, infohash)
    g := krpc.BytesToMessage(encoded)

    response := krpc.GetPeersRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.InfoHash.Equals(infohash) ||
        response.T.String() != tx.String(){
        t.Error()
    }
}

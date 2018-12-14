package test

import (
    "bytes"
    "kademlia/datastructure"
    "kademlia/network/messages"
    "testing"
)

func TestGetPeersResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    tx := "aaeebb"
    token := []byte("token")
    values := []string{"abc", "def"}
    encoded := messages.GetPeersResponse{}.Encode(tx, randomNodeID, token, values)
    g := messages.BytesToMessage(encoded)

    response := messages.GetPeersResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        response.T != tx ||
        !bytes.Equal(response.Token, token) ||
        response.Values[0] != values[0] ||
        response.Values[1] != values[1]{
        t.Error("")
    }
}

func TestGetPeersResponseWithNodes(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    n2 := datastructure.FakeNodeID(0x13)
    n3 := datastructure.FakeNodeID(0x14)
    tx := "aaeebb"
    token := []byte("token")
    nodes := []datastructure.NodeID{n2, n3}
    encoded := messages.GetPeersResponseWithNodes{}.Encode(tx, randomNodeID, token, nodes)
    g := messages.BytesToMessage(encoded)

    response := messages.GetPeersResponseWithNodes{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        response.T != tx ||
        !bytes.Equal(response.Token, token) ||
        !response.Nodes[0].Equals(nodes[0])  ||
        !response.Nodes[1].Equals(nodes[1]){
        t.Error("")
    }
}

func TestGetPeersRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    infohash := datastructure.FakeNodeID(0xF4)
    tx := "aaeebb"
    encoded := messages.GetPeersRequest{}.Encode(tx, randomNodeID, infohash)
    g := messages.BytesToMessage(encoded)

    response := messages.GetPeersRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.InfoHash.Equals(infohash) ||
        response.T != tx{
        t.Error()
    }
}

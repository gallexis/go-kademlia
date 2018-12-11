package test

import (
    "github.com/ehmry/go-bencode"
    "kademlia/datastructure"
    "kademlia/network/messages"
    "log"
    "testing"
)

func TestGetPeersResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    tx := "aaeebb"
    token := "aaeebb"
    values := []string{"abc", "def"}
    encoded := messages.GetPeersResponse{}.Encode(tx, randomNodeID, token, values)

    g := messages.GenericMessage{}
    if err := bencode.Unmarshal(encoded, &g); err != nil {
        log.Fatalln(err.Error())
    }

    response := messages.GetPeersResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        response.T != tx ||
        response.Token != token ||
        response.Values[0] != values[0] ||  response.Values[1] != values[1]{
        t.Error("")
    }
}

func TestGetPeersResponseWithNodes(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    n2 := datastructure.FakeNodeID(0x13)
    n3 := datastructure.FakeNodeID(0x14)
    tx := "aaeebb"
    token := "aaeebb"
    nodes := []datastructure.NodeID{n2, n3}
    encoded := messages.GetPeersResponseWithNodes{}.Encode(tx, randomNodeID, token, nodes)

    g := messages.GenericMessage{}
    if err := bencode.Unmarshal(encoded, &g); err != nil {
        log.Fatalln(err.Error())
    }

    response := messages.GetPeersResponseWithNodes{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        response.T != tx ||
        response.Token != token ||
        response.Nodes[0] != nodes[0] ||  response.Nodes[1] != nodes[1]{
        t.Error("")
    }
}

func TestGetPeersRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    infohash := datastructure.FakeNodeID(0xF4)
    tx := "aaeebb"
    encoded := messages.GetPeersRequest{}.Encode(tx, randomNodeID, infohash)

    g := messages.GenericMessage{}
    if err := bencode.Unmarshal(encoded, &g); err != nil {
        log.Fatalln(err.Error())
    }

    response := messages.GetPeersRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.InfoHash.Equals(infohash) ||
        response.T != tx{
        t.Error()
    }
}

package test

import (
    "fmt"
    "kademlia/datastructure"
    "kademlia/network/messages"
    "testing"
)

func TestAnnouncePeerResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    tx := "aaeebb"
    encoded := messages.AnnouncePeersResponse{}.Encode(tx, randomNodeID)
    g := messages.BytesToMessage(encoded)
    response := messages.AnnouncePeersResponse{}
    response.Decode(g.T, g.R)

    fmt.Println(response.T)

    if !response.Id.Equals(randomNodeID) || response.T != tx {
        t.Error("")
    }
}

func TestAnnouncePeerRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    infohash := datastructure.FakeNodeID(0xF4)
    tx := "aaeebb"
    token := "token"
    impliedPort := 1
    port := 1337
    encoded := messages.AnnouncePeersRequest{}.Encode(tx, token, randomNodeID, infohash, impliedPort, port)
    g := messages.BytesToMessage(encoded)

    response := messages.AnnouncePeersRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.InfoHash.Equals(infohash) ||
        response.T != tx ||
        response.Port != port ||
        response.ImpliedPort != impliedPort {
        t.Error()
    }
}

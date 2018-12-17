package test

import (
    "kademlia/datastructure"
    "kademlia/message"
    "testing"
)

func TestAnnouncePeerResponse(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    tx := message.NewRandomBytesFromString("aaeebb")
    encoded := message.AnnouncePeersResponse{}.Encode(tx, randomNodeID)
    g := message.BytesToMessage(encoded)
    response := message.AnnouncePeersResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) || response.T.String() != tx.String() {
        t.Error("")
    }
}

func TestAnnouncePeerRequest(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    infohash := ds.FakeNodeID(0xF4)
    tx := message.NewRandomBytesFromString("aaeebb")
    token := message.NewRandomBytesFromString("bbaaee")
    impliedPort := 1
    port := 1337
    encoded := message.AnnouncePeersRequest{}.Encode(tx, token, randomNodeID, infohash, impliedPort, port)
    g := message.BytesToMessage(encoded)

    response := message.AnnouncePeersRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.InfoHash.Equals(infohash) ||
        response.T.String() != tx.String() ||
        response.Port != port ||
        response.ImpliedPort != impliedPort {
        t.Error()
    }
}

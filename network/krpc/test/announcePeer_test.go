package test

import (
    "kademlia/datastructure"
    "kademlia/network/krpc"
    "testing"
)

func TestAnnouncePeerResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    encoded := krpc.AnnouncePeersResponse{}.Encode(tx, randomNodeID)
    g := krpc.BytesToMessage(encoded)
    response := krpc.AnnouncePeersResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) || response.T.String() != tx.String() {
        t.Error("")
    }
}

func TestAnnouncePeerRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    infohash := datastructure.FakeNodeID(0xF4)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    token := krpc.NewRandomBytesFromString("bbaaee")
    impliedPort := 1
    port := 1337
    encoded := krpc.AnnouncePeersRequest{}.Encode(tx, token, randomNodeID, infohash, impliedPort, port)
    g := krpc.BytesToMessage(encoded)

    response := krpc.AnnouncePeersRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.InfoHash.Equals(infohash) ||
        response.T.String() != tx.String() ||
        response.Port != port ||
        response.ImpliedPort != impliedPort {
        t.Error()
    }
}

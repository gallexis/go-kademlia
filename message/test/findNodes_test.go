package test

import (
    ds "kademlia/datastructure"
    "kademlia/message"
    "testing"
)

func TestFindNodeResponse(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    randomNodeID2 := ds.FakeNodeID(0xF4)
    tx := message.NewRandomBytesFromString("aaeebb")
    encoded := message.FindNodeResponse{}.Encode(tx, randomNodeID, []ds.NodeID{randomNodeID, randomNodeID2})
    g := message.BytesToMessage(encoded)

    response := message.FindNodeResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        !response.Contact[0].Equals(randomNodeID)||
        !response.Contact[1].Equals(randomNodeID2){
        t.Error("")
    }
}


func TestFindNodeRequest(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    randomNodeID2 := ds.FakeNodeID(0xF4)
    tx := message.NewRandomBytesFromString("aaeebb")
    encoded := message.FindNodeRequest{}.Encode(tx, randomNodeID, randomNodeID2)
    g := message.BytesToMessage(encoded)

    response := message.FindNodeRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.Target.Equals(randomNodeID2){
        t.Error()
    }
}

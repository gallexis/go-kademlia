package test

import (
    "kademlia/datastructure"
    "kademlia/network/krpc"
    "testing"
)

func TestFindNodeResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    randomNodeID2 := datastructure.FakeNodeID(0xF4)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    encoded := krpc.FindNodeResponse{}.Encode(tx, randomNodeID, []datastructure.NodeID{randomNodeID, randomNodeID2})
    g := krpc.BytesToMessage(encoded)

    response := krpc.FindNodeResponse{}
    response.Decode(g.T, g.R)

    if !response.Id.Equals(randomNodeID) ||
        !response.Nodes[0].Equals(randomNodeID)||
        !response.Nodes[1].Equals(randomNodeID2){
        t.Error("")
    }
}


func TestFindNodeRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    randomNodeID2 := datastructure.FakeNodeID(0xF4)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    encoded := krpc.FindNodeRequest{}.Encode(tx, randomNodeID, randomNodeID2)
    g := krpc.BytesToMessage(encoded)

    response := krpc.FindNodeRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.Target.Equals(randomNodeID2){
        t.Error()
    }
}

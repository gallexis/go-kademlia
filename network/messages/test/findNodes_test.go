package test

import (
    "kademlia/datastructure"
    "kademlia/network/messages"
    "testing"
)

func TestFindNodeResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    randomNodeID2 := datastructure.FakeNodeID(0xF4)
    encoded := messages.FindNodeResponse{}.Encode("aaee", randomNodeID, []datastructure.NodeID{randomNodeID, randomNodeID2})
    g := messages.BytesToMessage(encoded)

    response := messages.FindNodeResponse{}
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
    encoded := messages.FindNodeRequest{}.Encode("aaee", randomNodeID, randomNodeID2)
    g := messages.BytesToMessage(encoded)

    response := messages.FindNodeRequest{}
    response.Decode(g.T, g.A)

    if !response.Id.Equals(randomNodeID) ||
        !response.Target.Equals(randomNodeID2){
        t.Error()
    }
}

package test

import (
    "kademlia/datastructure"
    "kademlia/network/messages"
    "testing"
)

func TestPingResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    encoded := messages.PingResponse{}.Encode("aafst", randomNodeID)
    g := messages.BytesToMessage(encoded)
    response := messages.PingResponse{}
    response.Decode(g.T, g.R.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

func TestPingRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    encoded := messages.PingRequest{}.Encode("aafst", randomNodeID)
    g := messages.BytesToMessage(encoded)
    response := messages.PingRequest{}
    response.Decode(g.T, g.A.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

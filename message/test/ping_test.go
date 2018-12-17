package test

import (
    "kademlia/datastructure"
    "kademlia/message"
    "testing"
)

func TestPingResponse(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    tx := message.NewRandomBytesFromString("aaeebb")
    encoded := message.PingResponse{}.Encode(tx, randomNodeID)
    g := message.BytesToMessage(encoded)
    response := message.PingResponse{}
    response.Decode(g.T, g.R.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

func TestPingRequest(t *testing.T) {
    randomNodeID := ds.FakeNodeID(0x12)
    tx := message.NewRandomBytesFromString("aaeebb")
    encoded := message.PingRequest{}.Encode(tx, randomNodeID)
    g := message.BytesToMessage(encoded)
    response := message.PingRequest{}
    response.Decode(g.T, g.A.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

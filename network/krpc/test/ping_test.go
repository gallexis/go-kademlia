package test

import (
    "kademlia/datastructure"
    "kademlia/network/krpc"
    "testing"
)

func TestPingResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    encoded := krpc.PingResponse{}.Encode(tx, randomNodeID)
    g := krpc.BytesToMessage(encoded)
    response := krpc.PingResponse{}
    response.Decode(g.T, g.R.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

func TestPingRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    tx := krpc.NewRandomBytesFromString("aaeebb")
    encoded := krpc.PingRequest{}.Encode(tx, randomNodeID)
    g := krpc.BytesToMessage(encoded)
    response := krpc.PingRequest{}
    response.Decode(g.T, g.A.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

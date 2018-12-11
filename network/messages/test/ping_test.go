package test

import (
    "github.com/ehmry/go-bencode"
    "kademlia/datastructure"
    "kademlia/network/messages"
    "log"
    "testing"
)

func TestPingResponse(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    encoded := messages.PingResponse{}.Encode("aafst", randomNodeID)

    g := messages.GenericMessage{}
    if err := bencode.Unmarshal(encoded, &g); err != nil {
        log.Fatalln(err.Error())
    }

    response := messages.PingResponse{}
    response.Decode(g.T, g.R.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

func TestPingRequest(t *testing.T) {
    randomNodeID := datastructure.FakeNodeID(0x12)
    encoded := messages.PingRequest{}.Encode("aafst", randomNodeID)

    g := messages.GenericMessage{}
    if err := bencode.Unmarshal(encoded, &g); err != nil {
        log.Fatalln(err.Error())
    }

    response := messages.PingRequest{}
    response.Decode(g.T, g.A.Id)

    if !response.Id.Equals(randomNodeID){
        t.Error()
    }
}

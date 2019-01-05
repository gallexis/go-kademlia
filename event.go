package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "kademlia/message"
    "time"
)

func (d *DHT) OnAnnouncePeerResponse(announcePeer message.AnnouncePeersResponse) {
    log.Info(announcePeer)
}

func (d *DHT) OnGetPeersResponse(getPeers message.GetPeersResponse) {
    log.Infof("getPeers: %+v", getPeers)
}

func (d *DHT) OnGetPeersWithNodesResponse(getPeersWithNodes message.GetPeersResponseWithNodes) {
    log.Infof("getPeersWithNodes: %+v", getPeersWithNodes)
}

// FIND NODES
func (d *DHT) OnFindNodesResponse(findNodes message.FindNodeResponse) {
    //log.Println("findNodes")
    fmt.Println("findnodes")

    for _, c := range findNodes.Nodes {
        d.routingTable.Insert(c, d.PingRequest)
    }
}

func (d *DHT) PopulateRT() {
    closestNodes := d.routingTable.GetClosestNodes()
    tx := message.NewTransactionId()

    for _, node := range closestNodes {
        findNodeRequest := message.FindNodeRequest{
            T:      tx,
            Id:     d.selfNodeID,
            Target: d.selfNodeID,
        }
        d.Send(findNodeRequest.Encode(), node.ContactInfo)
    }

    if d.routingTable.ClosestBucketFilled < MinBucketFilled{
        d.eventDispatcher.AddEvent(tx.String(), Event{
            timeout:           time.Now(),
            maxTries:          1,
            duplicates:        8,
            CallbackOnTimeout: NewCallback(d.PopulateRT),
            Callback:          NewCallback(d.OnFindNodesResponse),
            Caller:            Callback{},
        })
    }
    fmt.Println(d.routingTable.ClosestBucketFilled, d.routingTable)
}

//----------------------------------------

// PING
func (d *DHT) OnPingResponse(ping message.PingResponse) {
    log.Infof("OnPingResponse: %+v", d.pingPool)

    exists := d.routingTable.UpdateNodeStatus(ping.Id)
    if !exists {
        return
    }

    if c, ok := d.pingPool[ping.T.String()]; ok {
        c <- true
    }
}

func (d *DHT) PingRequest(pingChan chan bool) {
    tx := message.NewTransactionId()
    pingRequest := message.PingRequest{
        T:  tx,
        Id: d.selfNodeID,
    }
    _, err := d.conn.Write(pingRequest.Encode())
    if err != nil {
        log.Error("Failed to send ping request")
        return
    }
    d.pingPool[tx.String()] = pingChan
}

//----------------------------------------

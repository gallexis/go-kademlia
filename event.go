package main

import (
    log "github.com/sirupsen/logrus"
    "kademlia/message"
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
    log.Infof("findNodes: %+v", findNodes)
    if !d.eventDispatcher.EventExists(findNodes.T.String()) {
        return
    }

    for _, c := range findNodes.Nodes {
        d.routingTable.Insert(c, d.PingRequest)
    }

    if d.latestBucketFilled <= d.routingTable.GetLatestBucketFilled() {
        d.PopulateRT()
    }
}

func (d *DHT) PopulateRT() {
    closestNodes := d.routingTable.GetClosestNodes()
    d.latestBucketFilled = d.routingTable.GetLatestBucketFilled()

    for _, node := range closestNodes {
        if !node.RequestFindNode() {
            continue
        }

        tx := message.NewTransactionId()
        findNodeRequest := message.FindNodeRequest{
            T:      tx,
            Id:     d.selfNodeID,
            Target: d.selfNodeID,
        }
        d.eventDispatcher.AddEvent(tx.String(), )
        d.Send(findNodeRequest.Encode(), node.ContactInfo)
    }
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

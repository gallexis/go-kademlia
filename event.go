package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "kademlia/datastructure"
    "kademlia/message"
    "time"
)

func (d *DHT) OnAnnouncePeerResponse(announcePeer *message.AnnouncePeersResponse) {
    log.Info(announcePeer)
}

// Get Peers
func (d *DHT) OnGetPeersResponse(infoHash datastructure.InfoHash, getPeers *message.GetPeersResponse) {
    log.Infof("OnGetPeersResponse----> : %+v", getPeers)

}

func (d *DHT) OnGetPeersWithNodesResponse(infoHash datastructure.InfoHash, getPeersWithNodes *message.GetPeersResponseWithNodes) {
    log.Infof("getPeersWithNodes")

    for _, c := range getPeersWithNodes.Nodes {
        d.routingTable.Insert(c, d.PingRequest)
    }

    if !d.peerStore.Contains(infoHash){
        
    }
}

func (d *DHT) OnGetPeers(infoHash datastructure.InfoHash, msg message.Message) {
    switch v := msg.(type) {
    case *message.GetPeersResponseWithNodes:
        d.OnGetPeersWithNodesResponse(infoHash, v)
    case *message.GetPeersResponse:
        d.OnGetPeersResponse(infoHash, v)
    default:
        log.Debug("Error : default case for OnGetPeers")
        return
    }
}

func (d *DHT) getPeersByNodes(infoHash datastructure.InfoHash, nodes []datastructure.Node) {
    fmt.Println("Inside Getpeers")
    tx := message.NewTransactionId()

    for _, node := range nodes {
        getPeersRequest := message.GetPeersRequest{
            T:        tx,
            Id:       d.selfNodeID,
            InfoHash: infoHash,
        }
        d.Send(getPeersRequest.Encode(), node.ContactInfo)
    }

    d.eventDispatcher.AddEvent(tx.String(), Event{
        timeout:           time.Now(),
        maxTries:          2,
        duplicates:        len(nodes),
        CallbackOnTimeout: NewCallback(d.getPeersByNodes, infoHash, nodes),
        Callback:          NewCallback(d.OnGetPeers, infoHash),
        Caller:            Callback{},
    })
}

func (d *DHT) GetPeers(infoHash datastructure.InfoHash) {
    fmt.Println("Inside Getpeers")
    nodes := d.routingTable.Get(infoHash)
    d.getPeersByNodes(infoHash, nodes)
}

//-------------------------

// FIND NODES
func (d *DHT) OnFindNodesResponse(findNodes *message.FindNodeResponse) {
    //log.Println("findNodes")
    //fmt.Println("findnodes")

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

    if d.routingTable.ClosestBucketFilled < MinBucketFilled {
        d.eventDispatcher.AddEvent(tx.String(), Event{
            timeout:           time.Now(),
            maxTries:          1,
            duplicates:        len(closestNodes),
            CallbackOnTimeout: NewCallback(d.PopulateRT),
            Callback:          NewCallback(d.OnFindNodesResponse),
            Caller:            Callback{},
        })
    }
    fmt.Println(d.routingTable.ClosestBucketFilled, d.routingTable)
}

//----------------------------------------

// PING
func (d *DHT) OnPingResponse(ping *message.PingResponse) {
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
        //log.Error("Failed to send ping request")
        return
    }
    d.pingPool[tx.String()] = pingChan
}

//----------------------------------------

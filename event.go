package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "kademlia/datastructure"
    "kademlia/message"
    "time"
)

func (d *DHT) OnAnnouncePeerResponse(announcePeer *message.AnnouncePeersResponse) {
    log.Info("OnAnnouncePeerResponse", announcePeer)
}

// GetK Peers
func (d *DHT) OnGetPeersResponse(infoHash datastructure.InfoHash, getPeers *message.GetPeersResponse) {
    log.Info("!!! OnGetPeersResponse !!!")

    d.peerStore.Add(infoHash, getPeers.Peers)
}

func (d *DHT) OnGetPeersWithNodesResponse(infoHash datastructure.InfoHash, getPeersWithNodes *message.GetPeersResponseWithNodes) {
    log.Debug("getPeersWithNodes")

    for _, c := range getPeersWithNodes.Nodes {
        d.routingTable.Insert(c, d.PingRequest)
    }

    if !d.peerStore.Contains(infoHash) {
        d.getPeersByNodes(infoHash, getPeersWithNodes.Nodes)
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
    tx := message.NewTransactionId()

    for _, node := range nodes {
        getPeersRequest := message.GetPeersRequest{
            T:        tx,
            Id:       d.selfNodeID,
            InfoHash: infoHash,
        }
        node.Send(d.conn, getPeersRequest.Encode())
    }

    d.eventDispatcher.AddEvent(tx.String(), Event{
        timeout:           time.Now(),
        maxTries:          1,
        duplicates:        len(nodes),
        CallbackOnTimeout: Callback{},
        Callback:          NewCallback(d.OnGetPeers, infoHash),
        Caller:            Callback{},
    })
}

func (d *DHT) GetPeers(infoHash datastructure.InfoHash) {
    nodes := d.routingTable.GetK(infoHash)
    d.getPeersByNodes(infoHash, nodes)
}

//-------------------------

// FIND NODES
func (d *DHT) OnFindNodesResponse(findNodes *message.FindNodeResponse) {
    //log.Debug("findNodes")

    for _, c := range findNodes.Nodes {
        d.routingTable.Insert(c, d.PingRequest)
    }
}

func (d *DHT) PopulateRT() {
    closestNodes := d.routingTable.GetClosestNodes()
    tx := message.NewTransactionId()
    totalGoodNodes := len(closestNodes)

    for _, node := range closestNodes {
        if !node.CanRequestFindNode(){
            totalGoodNodes -= 1
            continue
        }
        findNodeRequest := message.FindNodeRequest{
            T:      tx,
            Id:     d.selfNodeID,
            Target: d.selfNodeID,
        }
        node.Send(d.conn, findNodeRequest.Encode())
        d.routingTable.UpdateLastRequestFindNode(node)
    }

    fmt.Println("populate >>>> ", totalGoodNodes, d.routingTable.ClosestBucketFilled)

    if totalGoodNodes > 0 && d.routingTable.ClosestBucketFilled < 158{
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

    //exists := d.routingTable.UpdateNodeStatus(ping.Id)
    //if !exists {
    //    return
    //}
    //
    //if c, ok := d.pingPool[ping.T.String()]; ok {
    //    c <- true
    //}
}

func (d *DHT) PingRequest(pingChan chan bool) {
    tx := message.NewTransactionId()
    pingRequest := message.PingRequest{
        T:  tx,
        Id: d.selfNodeID,
    }
    _, err := d.conn.Write(pingRequest.Encode())
    if err != nil {
        log.Error("Failed to send ping request to ", )
        return
    }
    d.pingPool[tx.String()] = pingChan
}

func (d *DHT) OnPingRequest(msg message.PingRequest) {
    pingResponse := message.PingResponse{
        T:  msg.T,
        Id: d.selfNodeID,
    }
    _, err := d.conn.Write(pingResponse.Encode())
    if err != nil {
        log.Error("Failed to send ping response")
        return
    }
}

//----------------------------------------

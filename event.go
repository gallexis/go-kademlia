package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    ds "kademlia/datastructure"
    "kademlia/message"
    "net"
    "time"
)

func (d *DHT) OnAnnouncePeerResponse(announcePeer *message.AnnouncePeersResponse, addr net.UDPAddr) {
    log.Info("OnAnnouncePeerResponse", announcePeer)
}

func (d *DHT) OnAnnouncePeerRequest(announcePeer *message.AnnouncePeersRequest, addr net.UDPAddr) {
    log.Info("OnAnnouncePeerRequest", announcePeer)
}

// GetK Peers
func (d *DHT) onGetPeersResponse(infoHash ds.InfoHash, getPeers *message.GetPeersResponse, addr net.UDPAddr) {
    log.Info("!!! onGetPeersResponse !!!")

    d.peerStore.Add(infoHash, getPeers.Peers)
}

func (d *DHT) onGetPeersWithNodesResponse(infoHash ds.InfoHash, getPeersWithNodes *message.GetPeersResponseWithNodes, addr net.UDPAddr) {
    log.Debug("getPeersWithNodes", addr)

    for _, c := range getPeersWithNodes.Nodes {
        d.Insert(c)
    }

    if !d.peerStore.Contains(infoHash) {
        d.getPeersByNodes(infoHash, getPeersWithNodes.Nodes)
    }
}

func (d *DHT) OnGetPeers(infoHash ds.InfoHash, msg message.Message, addr net.UDPAddr) {
    switch v := msg.(type) {
    case *message.GetPeersResponseWithNodes:
        d.onGetPeersWithNodesResponse(infoHash, v, addr)
    case *message.GetPeersResponse:
        d.onGetPeersResponse(infoHash, v, addr)
    default:
        log.Debug("Error : default case for OnGetPeers")
        return
    }
}

func (d *DHT) getPeersByNodes(infoHash ds.InfoHash, nodes []ds.Node) {
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

func (d *DHT) GetPeers(infoHash ds.InfoHash) {
    nodes := d.routingTable.GetK(infoHash)
    d.getPeersByNodes(infoHash, nodes)
}

//-------------------------

// FIND NODES
func (d *DHT) OnFindNodesResponse(findNodes *message.FindNodeResponse, addr net.UDPAddr) {
    //log.Debug("findNodes")

    for _, c := range findNodes.Nodes {
        d.Insert(c)
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
func (d *DHT) OnPingResponse(node ds.Node ,ping *message.PingResponse, addr net.UDPAddr) {
    log.Info("OnPingResponse")
    d.routingTable.UpdateNodeStatus(ping.Id)
    d.PingPool <- node
}

func (d *DHT) OnPingRequest(msg *message.PingRequest, addr net.UDPAddr) {
    log.Info("OnPingRequest")
    pingResponse := message.PingResponse{
        T:  msg.T,
        Id: d.selfNodeID,
    }
    _, err := d.conn.WriteToUDP(pingResponse.Encode(), &addr)
    if err != nil {
        log.Error("Failed to send ping response")
    }
}
//----------------------------------------

// FIND NODE Request
func (d *DHT) OnFindNodeRequest(msg *message.FindNodeRequest, addr net.UDPAddr) {
    log.Info("OnFindNodeRequest")

    nodes := d.routingTable.GetK(msg.Id)

    findNodeResponse := message.FindNodeResponse{
        T:     msg.T,
        Id:    d.selfNodeID,
        Nodes: nodes,
    }
    _, err := d.conn.WriteToUDP(findNodeResponse.Encode(), &addr)
    if err != nil {
        log.Error("Failed to send findNode response")
    }
}
//----------------------------------------

// GetPeers Request
func (d *DHT) OnGetPeersRequest(msg *message.GetPeersRequest, addr net.UDPAddr) {
    log.Info("OnGetPeersRequest")
    var data []byte

    peers := d.peerStore.Get(msg.InfoHash)
    if len(peers) > 0{
        data = message.GetPeersResponse{
            T:     msg.T,
            Id:    d.selfNodeID,
            Token: message.Token("sdhh"),
            Peers: peers,
        }.Encode()
    } else{
        nodes := d.routingTable.GetK(msg.Id)

        data = message.GetPeersResponseWithNodes{
            T:     msg.T,
            Id:    d.selfNodeID,
            Token: message.Token("sdhh"),
            Nodes: nodes,
        }.Encode()
    }

    _, err := d.conn.WriteToUDP(data, &addr)
    if err != nil {
        log.Error("Failed to send GetPeersWithNodes response")
    }
}
//----------------------------------------

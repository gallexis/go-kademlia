package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "kademlia/Dispatcher"
    "kademlia/message"
    "net"
)

func (d *DHT) onFindNodeRequest(msg *message.FindNodeRequest, addr net.UDPAddr) {
    log.Info("onFindNodeRequest")

    nodes := d.routingTable.GetK(msg.Id)

    findNodeResponse := message.FindNodeResponse{
        T:     msg.T,
        Id:    d.selfNodeID,
        Nodes: nodes,
    }

    if _, err := d.conn.WriteToUDP(findNodeResponse.Encode(), &addr); err != nil {
        log.Error("Failed to send findNode response")
    }
}

func (d *DHT) onFindNodesResponse(findNodes *message.FindNodeResponse, addr net.UDPAddr) {
    log.Printf("findNodes %+v", addr)

    for _, c := range findNodes.Nodes {
        d.insert(c)
    }
}

func (d *DHT) populateRT() {
    log.Debug("populateRT")

    closestNodes := d.routingTable.GetClosestNodes()
    tx := message.NewTransactionId()
    totalGoodNodes := len(closestNodes)

    for _, node := range closestNodes {
        if !node.CanRequestFindNode() {
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

    if totalGoodNodes > 0 && d.routingTable.GetClosestBucketFilled() < 158 {
        d.eventDispatcher.AddEvent(tx.String(), Dispatcher.Event{
            Duplicates: len(closestNodes),
            OnTimeout:  Dispatcher.NewCallback(d.populateRT),
            OnResponse: Dispatcher.NewCallback(d.onFindNodesResponse),
        })
    }
    fmt.Println(d.routingTable.ClosestBucketFilled, d.routingTable)
}

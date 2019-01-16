package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    ds "kademlia/datastructure"
    "kademlia/message"
    "net"
    "time"
)

const (
    MinBucketFilled = 20
)

type DHT struct {
    selfNodeID         ds.NodeId
    routingTable       ds.RoutingTable
    pingPool           map[string]chan bool
    conn               *net.UDPConn
    bootstrapNodes     []string
    peerStore          ds.PeerStore
    eventDispatcher    Dispatcher
}

func NewDHT() DHT {
    nid := ds.NewNodeID()
    bootstrapNodes := []string{
        "router.utorrent.com:6881",
        "router.bittorrent.com:6881",
    }
    addr := net.UDPAddr{
        Port: 38568,
        IP:   net.ParseIP("0.0.0.0"),
    }
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Panicf("Some error %v", err)
    }

    return NewCustomDHT(nid, bootstrapNodes, conn)
}

func NewCustomDHT(nid ds.NodeId, bootstrapNodes []string, conn *net.UDPConn) DHT {
    return DHT{
        selfNodeID:         nid,
        routingTable:       ds.NewRoutingTable(nid),
        pingPool:           make(map[string]chan bool),
        conn:               conn,
        bootstrapNodes:     bootstrapNodes,
        peerStore:          nil,
        eventDispatcher:    NewDispatcher(),
    }
}

func (d *DHT) Bootstrap(bootstrapNode string) bool {
    raddr, err := net.ResolveUDPAddr("udp", bootstrapNode)
    if err != nil {
        log.Panic("Can't resolve")
    }

    conn, err := net.DialUDP("udp", nil, raddr)
    if err != nil {
        log.Panic("can't dial")
    }

    deadline := time.Now().Add(time.Second * 10)
    err = conn.SetReadDeadline(deadline)
    if err != nil {
        log.Panic("too long")
    }

    // Send FindNode request
    findNodeRequest := message.FindNodeRequest{
        T:      message.NewTransactionId(),
        Id:     d.selfNodeID,
        Target: d.selfNodeID,
    }
    _, err = conn.Write(findNodeRequest.Encode())
    if err != nil {
        log.Panic("NewTransactionId", err)
        return false
    }

    // GetResponse
    buffer := make([]byte, 1024)
    _, _, err = conn.ReadFrom(buffer)
    if err != nil {
        log.Panic("can't read", err.Error())
    }

    // Assert response if findNodeResponse
    g, _ := message.BytesToMessage(buffer)
    if g.Y != "r" || len(g.R.Nodes) <= 0 {
        log.Panic("Not event findNode Response")
    }

    // Decode FindNodeResponse
    findNodes := message.FindNodeResponse{}
    findNodes.Decode(g)

    // Update routing table with nodes received
    for _, c := range findNodes.Nodes {
        d.routingTable.Insert(c, func(chan bool) {})
    }

    return true
}

func (d *DHT) Receiver() {
    go func() {
        buffer := make([]byte, 1024)

        for {
            n, _, err := d.conn.ReadFromUDP(buffer)
            if err != nil {
                log.Printf("Some error %v", err)
                time.Sleep(time.Second * 1)
                continue
            }
            d.Router(buffer[:n])
        }
    }()
}

func (d *DHT) Send(data []byte, contact ds.Contact) {
    destAddr := net.UDPAddr{IP: contact.IP, Port: int(contact.Port)}
    _, err := d.conn.WriteToUDP(data, &destAddr)
    if err != nil {
        log.Error("DHT.Send", err.Error())
    }
}

func (d *DHT) Router(data []byte) {
    var msg message.Message

    g, ok := message.BytesToMessage(data)
    if !ok{
        return
    }

    switch g.Y {
    case "q":

        switch g.Q {
        case "ping":
            log.Info("PingRequest")
            msg = &message.PingRequest{}

        case "find_node":
            log.Info("FindNodeRequest")
            msg = &message.FindNodeRequest{}

        case "get_peers":
            log.Info("GetPeersRequest")
            msg = &message.GetPeersRequest{}

        case "announce_peer":
            log.Info("AnnouncePeersRequest")
            msg = &message.AnnouncePeersRequest{}

        default:
            log.Panic("q")
        }

        msg.Decode(g)
        fmt.Printf("Receive Query : %+v \n", msg)

    case "r":
        callback, exists := d.eventDispatcher.GetCallback(g.T)
        if !exists {
            return
        }

        switch {
        case len(g.R.Values) > 0 && len(g.R.Token) > 0:
            msg = &message.GetPeersResponse{}

        case len(g.R.Nodes) > 0 && len(g.R.Token) > 0:
            msg = &message.GetPeersResponseWithNodes{}

        case len(g.R.Nodes) > 0:
            msg = &message.FindNodeResponse{}

        case len(g.R.Id) > 0:
            msg = &message.PingResponse{}

        default:
            log.Panic("r", g.Y)
        }

        msg.Decode(g)
        callback.Call(msg)

    case "e":
        log.Info("Error:", g.E)

    }
}

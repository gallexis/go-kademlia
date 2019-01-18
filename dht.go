package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    ds "kademlia/datastructure"
    "kademlia/message"
    "net"
    "time"
)

type DHT struct {
    selfNodeID      ds.NodeId
    routingTable    ds.RoutingTable
    conn            *net.UDPConn
    bootstrapNodes  []string
    peerStore       ds.PeerStore
    eventDispatcher Dispatcher
    PingPool        chan ds.Node
    pingRequests    chan ds.Node
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
    pingRequests := make(chan ds.Node)

    return DHT{
        pingRequests:    pingRequests,
        selfNodeID:      nid,
        routingTable:    ds.NewRoutingTable(nid, pingRequests),
        conn:            conn,
        bootstrapNodes:  bootstrapNodes,
        peerStore:       make(ds.PeerStore),
        eventDispatcher: NewDispatcher(),
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
        _, _ = d.routingTable.Insert(c, true)
    }

    return true
}

func (d *DHT) Receiver() {
    go func() {
        buffer := make([]byte, 1024)

        for {
            n, udpAddr, err := d.conn.ReadFromUDP(buffer)
            if err != nil {
                log.Printf("Some error %v", err)
                time.Sleep(time.Second * 1)
                continue
            }
            d.Router(buffer[:n], *udpAddr)
        }
    }()
}

func (d *DHT) ManageRequest(request message.Message, addr net.UDPAddr) {
    switch v := request.(type) {
    case *message.PingRequest:
        d.OnPingRequest(v, addr)
    case *message.FindNodeRequest:
        d.OnFindNodeRequest(v, addr)
    case *message.GetPeersRequest:
        d.OnGetPeersRequest(v, addr)
    case *message.AnnouncePeersRequest:
        d.OnAnnouncePeerRequest(v, addr)
    default:
        fmt.Println("Unknown request")
    }
}

func (d *DHT) PendingPingPool() {
    go func() {
        for {
            select {
            case node := <-d.pingRequests:
                tx := message.NewTransactionId()
                d.eventDispatcher.AddEvent(tx.String(), Event{
                    timeout:           time.Now(),
                    maxTries:          2,
                    duplicates:        0,
                    CallbackOnTimeout: NewCallback(func() {
                        fmt.Println("Node is dead")
                        d.routingTable.Remove(node)
                    }),
                    Callback:          NewCallback(d.OnPingResponse, node),
                    Caller:            NewCallback(d.SendPingRequest, node, tx),
                })

                d.SendPingRequest(node, tx)
            }
        }

    }()
}

func (d *DHT) Insert(node ds.Node) {
    ok, err := d.routingTable.Insert(node, false)
    if ok {
        return
    } else if err != nil {
        log.Error(err)
        return
    }

    tx := message.NewTransactionId()

    // Not inserted because we are waiting for a ping response
    // add in pending pool
    // send ping
    d.eventDispatcher.AddEvent(tx.String(), Event{
        timeout:  time.Now(),
        maxTries: 2,
        CallbackOnTimeout: NewCallback(func() {
            ok, err := d.routingTable.Insert(node, true)
            if !ok {
                log.Error("should have inserted node properly")
            }
            if err != nil {
                log.Error("error when inserting new node", err)
            }
        }),
        Callback: NewCallback(d.OnPingResponse, node),
        Caller:   NewCallback(d.SendPingRequest, node, tx),
    })

    d.SendPingRequest(node, tx)
}

func (d *DHT) Router(data []byte, addr net.UDPAddr) {
    var msg message.Message

    g, ok := message.BytesToMessage(data)
    if !ok {
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
        d.ManageRequest(msg, addr)

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
        callback.Call(msg, addr)

    case "e":
        log.Info("Error:", g.E)

    }
}

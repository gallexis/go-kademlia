package main

import (
    log "github.com/sirupsen/logrus"
    ds "kademlia/datastructure"
    "kademlia/message"
    "net"
    "time"
)

type DHT struct {
    selfNodeID      ds.NodeId
    routingTable    ds.RoutingTable
    pingPool        map[string]chan bool
    conn            *net.UDPConn
    bootstrapNodes  []string
    peerStore       ds.PeerStore
    eventDispatcher Dispatcher
    latestBucketFilled int
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
        pingPool:           nil,
        conn:               conn,
        bootstrapNodes:     bootstrapNodes,
        peerStore:          nil,
        eventDispatcher:    NewDispatcher(),
        latestBucketFilled: ds.BitsInNodeID,
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
    g := message.BytesToMessage(buffer)
    if g.Y != "r" || len(g.R.Nodes) <= 0 {
        log.Panic("Not event findNode Response")
    }

    // Decode FindNodeResponse
    findNodes := message.FindNodeResponse{}
    findNodes.Decode(g.T, g.R)

    // Update routing table with nodes received
    for _, c := range findNodes.Nodes {
        d.routingTable.Insert(c, func(chan bool) {})
    }

    return true
}


func (d *DHT) Receiver() {
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
}

func (d *DHT) Send(data []byte, contact ds.Contact) {
    destAddr := net.UDPAddr{IP: contact.IP, Port: int(contact.Port)}
    _, err := d.conn.WriteToUDP(data, &destAddr)
    if err != nil {
        log.Error(err.Error())
    }
}


func (d *DHT) Router(data []byte) {
    g := message.BytesToMessage(data)

    switch g.Y {
    case "q":

        switch g.Q {
        case "ping":
            ping := message.PingRequest{}
            ping.Decode(g.T, g.A.Id)
            log.Info("PingRequest")

        case "find_node":
            findNodesRequest := message.FindNodeRequest{}
            findNodesRequest.Decode(g.T, g.A)
            log.Info("FindNodeRequest")

        case "get_peers":
            getPeers := message.GetPeersRequest{}
            getPeers.Decode(g.T, g.A)
            log.Info("GetPeersRequest")

        case "announce_peer":
            announcePeers := message.AnnouncePeersRequest{}
            announcePeers.Decode(g.T, g.A)
            log.Info("AnnouncePeersRequest")

        default:
            log.Panic("q")
        }

    case "r":

        switch {
        case len(g.R.Values) > 0 && len(g.R.Token) > 0:
            getPeers := message.GetPeersResponse{}
            getPeers.Decode(g.T, g.R)
            d.OnGetPeersResponse(getPeers)

        case len(g.R.Nodes) > 0 && len(g.R.Token) > 0:
            getPeers := message.GetPeersResponseWithNodes{}
            getPeers.Decode(g.T, g.R)
            d.OnGetPeersWithNodesResponse(getPeers)

        case len(g.R.Nodes) > 0:
            findNodes := message.FindNodeResponse{}
            findNodes.Decode(g.T, g.R)
            d.OnFindNodesResponse(findNodes)
            f, exists := d.eventDispatcher.EventExists(findNodes.T.String(), false)
            if exists {
                f()
            }

        case len(g.R.Id) > 0:
            ping := message.PingResponse{}
            ping.Decode(g.T, g.R.Id)
            d.OnPingResponse(ping)

            /*
            AnnouncePeersResponse == PingResponse

                announcePeers := AnnouncePeersResponse{}
                announcePeers.Decode(g.TransactionId, g.R)
                log.Println(announcePeers)
             */

        default:
            log.Panic("r", g.Y)
        }

    case "e":
        log.Info("Error:", g.E)

    }
}

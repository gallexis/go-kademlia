package main

import (
    "fmt"
    ds "kademlia/datastructure"
    "kademlia/message"
    "net"
    "time"
)

type DHT struct {
    nodeID         ds.NodeID
    routingTable   ds.RoutingTable
    handshakePool  map[string]time.Time
    pingPool       map[string]chan bool
    conn           *net.UDPConn
    bootstrapNodes []string
    port           string
    //peerstore
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
        fmt.Printf("Some error %v\n", err)
    }

    return NewCustomDHT(nid, bootstrapNodes, conn)
}

func NewCustomDHT(nid ds.NodeID, bootstrapNodes []string, conn *net.UDPConn) DHT {
    return DHT{
        nodeID:         nid,
        routingTable:   ds.NewRoutingTable(nid),
        bootstrapNodes: bootstrapNodes,
        conn:           conn,
        port:           "38568",
    }
}

func (d *DHT) Bootstrap(bootstrapNode string) bool {
    raddr, err := net.ResolveUDPAddr("udp", bootstrapNode)
    if err != nil {
        fmt.Println("can't resolve")
        return false
    }

    conn, err := net.DialUDP("udp", nil, raddr)
    if err != nil {
        fmt.Println("can't dial")
        return false
    }

    deadline := time.Now().Add(time.Second * 10)
    err = conn.SetReadDeadline(deadline)
    if err != nil {
        fmt.Println("too long")
        return false
    }

    // Send FindNode request
    tx := message.NewRandomBytes(2)
    _, err = conn.Write(message.FindNodeRequest{}.Encode(tx, d.nodeID, d.nodeID))
    if err != nil {
        return false
    }

    // GetResponse
    buffer := make([]byte, 2048)
    _, _, err = conn.ReadFrom(buffer)
    if err != nil {
        fmt.Println("can't read", err.Error())
        return false
    }

    // Assert response if findNodeResponse
    g := message.BytesToMessage(buffer)
    if g.Y != "r" || len(g.R.Nodes) <= 0 {
        fmt.Println("Not a findNode Response")
        return false
    }

    // Decode FindNodeResponse
    findNodes := message.FindNodeResponse{}
    findNodes.Decode(g.T, g.R)

    // Update routing table with nodes received
    for _, c := range findNodes.Contact {
        d.routingTable.Insert(c, func(chan bool){})
    }

    return true
}

func (d *DHT) PopulateRT(){
    smallest := 159

    for {
        time.Sleep(time.Second * 10)
        randomContacts := d.routingTable.GetClosest()
        newSmallest := d.routingTable.GetLatestBucketFilled()

        if newSmallest >= smallest{
            fmt.Println("exit")
            break
        } else{
            smallest = newSmallest
        }

        for _, contact := range randomContacts{
            tx := message.NewRandomBytes(2)
            d.Send(message.FindNodeRequest{}.Encode(tx, d.nodeID, d.nodeID), contact)
        }

        fmt.Println(">>> smallest:", smallest)
    }
}


func (d *DHT) Receiver() {
    buffer := make([]byte, 1024)

    for {
        n, _, err := d.conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        d.Router(buffer[:n])
    }
}

func (d *DHT) Send(data []byte, contact ds.Contact) {
    destAddr := net.UDPAddr{IP: contact.IP, Port: int(contact.Port)}
    _, err := d.conn.WriteToUDP(data, &destAddr)
    if err != nil {
        fmt.Println(err)
    }
}

func (d *DHT) OnAnnouncePeerResponse(announcePeer message.AnnouncePeersResponse) {
    fmt.Println(announcePeer)
}

func (d *DHT) OnFindNodesResponse(findNodes message.FindNodeResponse) {
    //fmt.Printf("findNodes: %+v \n", findNodes)
    fmt.Printf("findNodes\n")

    for _, c := range findNodes.Contact {
        d.routingTable.Insert(c, d.PingNode)
    }

    d.routingTable.Display()
    fmt.Println("-------------------")

}

func (d *DHT) OnGetPeersResponse(getPeers message.GetPeersResponse) {
    fmt.Printf("getPeers: %+v \n", getPeers)
}

func (d *DHT) OnGetPeersWithNodesResponse(getPeersWithNodes message.GetPeersResponseWithNodes) {
    fmt.Printf("getPeersWithNodes: %+v \n", getPeersWithNodes)
}

func (d *DHT) OnPingResponse(ping message.PingResponse) {
    fmt.Printf("OnPingResponse: %+v \n", ping)
    if c, ok := d.pingPool[ping.T.String()]; ok {
        c <- true
    }
}

func (d *DHT) PingNode(pingChan chan bool) {
    tx := message.NewRandomBytes(2)
    _, err := d.conn.Write(message.PingRequest{}.Encode(tx, d.nodeID))
    if err != nil {
        return
    }
    d.pingPool[tx.String()] = pingChan
}

func (d *DHT) Router(data []byte) {
    g := message.BytesToMessage(data)
    switch g.Y {
    case "q":

        switch g.Q {
        case "ping":
            ping := message.PingRequest{}
            ping.Decode(g.T, g.A.Id)
            fmt.Println("PingRequest")

        case "find_node":
            findNodesRequest := message.FindNodeRequest{}
            findNodesRequest.Decode(g.T, g.A)
            fmt.Println("FindNodeRequest")

        case "get_peers":
            getPeers := message.GetPeersRequest{}
            getPeers.Decode(g.T, g.A)
            fmt.Println("GetPeersRequest")

        case "announce_peer":
            announcePeers := message.AnnouncePeersRequest{}
            announcePeers.Decode(g.T, g.A)
            fmt.Println("AnnouncePeersRequest")

        default:
            panic("q")
        }

    case "r":

        switch {
        case len(g.R.Values) > 0 && g.R.Token != "":
            getPeers := message.GetPeersResponse{}
            getPeers.Decode(g.T, g.R)
            d.OnGetPeersResponse(getPeers)

        case len(g.R.Nodes) > 0 && g.R.Token != "":
            getPeers := message.GetPeersResponseWithNodes{}
            getPeers.Decode(g.T, g.R)
            d.OnGetPeersWithNodesResponse(getPeers)

        case len(g.R.Nodes) > 0:
            findNodes := message.FindNodeResponse{}
            findNodes.Decode(g.T, g.R)
            d.OnFindNodesResponse(findNodes)

        case len(g.R.Id) > 0:
            ping := message.PingResponse{}
            ping.Decode(g.T, g.R.Id)
            d.OnPingResponse(ping)

            /*
            AnnouncePeersResponse == PingResponse

                announcePeers := AnnouncePeersResponse{}
                announcePeers.Decode(g.RandomBytes, g.R)
                fmt.Println(announcePeers)
             */

        default:
            panic("r")
        }

    case "e":
        fmt.Println("Error:", g.E)

    }
}

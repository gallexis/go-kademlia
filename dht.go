package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	ds "kademlia/datastructure"
	"kademlia/dispatcher"
	"kademlia/message"
	"math/rand"
	"net"
	"time"
)

type DHT struct {
	selfNodeID      ds.NodeId
	routingTable    ds.RoutingTable
	conn            *net.UDPConn
	bootstrapNodes  []string
	peerStore       ds.PeerStore
	eventDispatcher dispatcher.Dispatcher
	token           message.Token
	PingPool        chan ds.Node
	pingRequests    chan ds.Node
	Callback        chan dispatcher.Callback
}

func NewDHT() DHT {
	nid := ds.NewNodeID()
	bootstrapNodes := []string{
		"router.utorrent.com:6881",
		"router.bittorrent.com:6881",
	}
	addr := net.UDPAddr{
		Port: 38569,
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
	callback := make(chan dispatcher.Callback)

	return DHT{
		selfNodeID:      nid,
		routingTable:    ds.NewRoutingTable(nid, pingRequests),
		conn:            conn,
		bootstrapNodes:  bootstrapNodes,
		peerStore:       ds.NewPeerStore(),
		eventDispatcher: dispatcher.NewDispatcher(callback),
		pingRequests:    pingRequests,
		Callback:        callback,
		token:           message.NewToken(), // todo : can receive tokens 5min older
	}
}

func (d *DHT) Init() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = fmt.Errorf("bootstrap has panicked: %v", t)
			case error:
				err = t
			default:
				err = fmt.Errorf("Unknown panic: %v", t)
			}
		}
	}()

	d.bootstrap()
	go d.eventDispatcher.Start()
	go d.callbackCaller()
	go d.receiver()
	go d.timer()
	d.populateRT()

	time.Sleep(time.Second * 5)

	//d.getPeers(ds.NewNodeIdFromString("57537D93A76F574369DC2E573E99C3840A9FD89D"))
	//d.getPeers(ds.NewNodeIdFromString("FC0CCE628DBE7EEA0CF655A6A13336791021F25F"))
	//d.getPeers(ds.NewNodeIdFromString("23E7A4876B36CE427A847A827306B4B2DC67304A"))
	d.getPeers(ds.NewNodeIdFromString("5EF929E35650741627DACA28E18A3DF0FC5A53DB"))
	//d.getPeers(ds.NewNodeIdFromString("D58952BDBBBFBA9DA444F8FE99DCF2C7F2E4AB77"))
	//d.getPeers(ds.NewNodeIdFromString("4EBF7D54EABA7380D46C05604B059FABAEA212F0"))

	go d.pendingPingPool()

	return
}

func (d *DHT) bootstrap() {
	hasBootstrapped := false
	var nodesTried []string
	for _, i := range rand.Perm(len(d.bootstrapNodes)) { // chose random node

		fmt.Println(d.bootstrapNodes[i])
		err := d.bootstrapNode(d.bootstrapNodes[i])
		if err == nil {
			hasBootstrapped = true
			break
		}
		nodesTried = append(nodesTried, d.bootstrapNodes[i])
	}
	if !hasBootstrapped {
		panic(fmt.Sprintf("bootstrapped nodes tried : %s", nodesTried))
	}
}

func (d *DHT) bootstrapNode(bootstrapNode string) error {
	var err error

	raddr, err := net.ResolveUDPAddr("udp", bootstrapNode)
	if err != nil {
		return fmt.Errorf("can't resolve UDP Addr : %v", err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return fmt.Errorf("can't dial UDP : %v", err)
	}

	deadline := time.Now().Add(time.Second * 10)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		return fmt.Errorf("network timeout : %v", err)
	}

	// Send FindNode request
	findNodeRequest := message.FindNodeRequest{
		T:      message.NewTransactionId(),
		Id:     d.selfNodeID,
		Target: d.selfNodeID,
	}
	_, err = conn.Write(findNodeRequest.Encode())
	if err != nil {
		return fmt.Errorf("can't send findNode request : %v", err)
	}

	// GetResponse
	buffer := make([]byte, 1024)
	_, _, err = conn.ReadFrom(buffer)
	if err != nil {
		return fmt.Errorf("can't read from socket : %v", err)
	}

	// Assert response if findNodeResponse
	g, _ := message.BytesToMessage(buffer)
	if g.Y != "r" || len(g.R.Nodes) <= 0 {
		return errors.New("incorrect findNode response")
	}

	// Decode FindNodeResponse
	findNodes := message.FindNodeResponse{}
	findNodes.Decode(g)

	if len(findNodes.Nodes) <= 0 {
		return errors.New("no nodes received")
	}

	// Update routing table with nodes received
	for _, c := range findNodes.Nodes {
		d.insert(c)
	}

	return nil
}

func (d *DHT) receiver() {
	buffer := make([]byte, 1024)

	for {
		n, udpAddr, err := d.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Some error %v", err)
			time.Sleep(time.Second * 1)
			continue
		}
		d.router(buffer[:n], *udpAddr)
	}
}

func (d *DHT) manageRequest(request message.Message, addr net.UDPAddr) {
	switch v := request.(type) {
	case *message.PingRequest:
		d.onPingRequest(v, addr)
	case *message.FindNodeRequest:
		d.onFindNodeRequest(v, addr)
	case *message.GetPeersRequest:
		d.onGetPeersRequest(v, addr)
	case *message.AnnouncePeersRequest:
		d.onAnnouncePeerRequest(v, addr)
	default:
		fmt.Println("Unknown request")
	}
}

func (d *DHT) timer() {
	displayRoutingTable := time.Tick(time.Second * 30)
	reGenerateToken := time.Tick(time.Minute * 10)

	for {
		select {
		case <-displayRoutingTable:
			fmt.Println(d.routingTable)

		case <-reGenerateToken:
			d.token = message.NewToken()

		}
	}
}

func (d *DHT) pendingPingPool() {
	for {
		select {
		case node := <-d.pingRequests:
			tx := message.NewTransactionId()
			d.eventDispatcher.AddEvent(tx.String(), dispatcher.Event{
				Retries: 1,
				OnTimeout: dispatcher.NewCallback(func() {
					d.routingTable.Remove(node)
				}),
				OnResponse: dispatcher.NewCallback(d.onPingResponse, node),
				OnRetry:    dispatcher.NewCallback(d.sendPingRequest, node, tx),
			})

			d.sendPingRequest(node, tx)
		}
	}
}

func (d *DHT) callbackCaller() {
	for {
		select {
		case callback := <-d.Callback:
			callback.Call()
		}
	}
}

func (d *DHT) insert(node ds.Node) {
	ok, err := d.routingTable.Insert(node, false)
	if ok {
		return
	} else if err != nil {
		log.Error(err)
		return
	}

	tx := message.NewTransactionId()

	// Not inserted because we are waiting for a ping response
	d.eventDispatcher.AddEvent(tx.String(), dispatcher.Event{
		Retries: 1,
		OnTimeout: dispatcher.NewCallback(func() {
			ok, err := d.routingTable.Insert(node, true)
			if !ok {
				log.Error("should have inserted node properly")
			}
			if err != nil {
				log.Error("error when inserting new node", err)
			}
		}),
		OnResponse: dispatcher.NewCallback(d.onPingResponse, node),
		OnRetry:    dispatcher.NewCallback(d.sendPingRequest, node, tx),
	})

	d.sendPingRequest(node, tx)
}

func (d *DHT) router(data []byte, addr net.UDPAddr) {
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
		d.manageRequest(msg, addr)

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
		callback.AddArgs(msg, addr)
		d.Callback <- callback

	case "e":
		log.Info("Error:", g.E)

	}
}

package main

import (
	log "github.com/sirupsen/logrus"
	"kademlia/datastructure"
	"kademlia/message"
	"net"
)

func (d *DHT) onAnnouncePeerResponse(announcePeer *message.AnnouncePeersResponse, addr net.UDPAddr) {
	log.Info("onAnnouncePeer ack", announcePeer)
}

func (d *DHT) onAnnouncePeerRequest(announcePeer *message.AnnouncePeersRequest, addr net.UDPAddr) {
	log.Info("onAnnouncePeer Request", announcePeer)

	if announcePeer.Token.Equals(d.token) {
		port := uint16(announcePeer.Port)

		if announcePeer.ImpliedPort == 1 {
			port = uint16(addr.Port)
		}

		d.peerStore.Add(announcePeer.InfoHash, datastructure.NewPeer(addr.IP, port))
		log.Debug("New peer announced : ", addr.IP, port)
	}

	log.Debug("announcePeer.Token !=  d.token ")
}

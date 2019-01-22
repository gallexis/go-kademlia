package main

import (
    log "github.com/sirupsen/logrus"
    "kademlia/message"
    "net"
)

func (d *DHT) OnAnnouncePeerResponse(announcePeer *message.AnnouncePeersResponse, addr net.UDPAddr) {
    log.Info("OnAnnouncePeerResponse", announcePeer)
}

func (d *DHT) OnAnnouncePeerRequest(announcePeer *message.AnnouncePeersRequest, addr net.UDPAddr) {
    log.Info("OnAnnouncePeerRequest", announcePeer)
}

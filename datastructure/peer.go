package datastructure

import (
    "encoding/binary"
    log "github.com/sirupsen/logrus"
    "net"
    "time"
)

type Peer struct {
    IP                  net.IP
    Port                uint16
    LastMessageReceived time.Time
    LastFindNodeRequest time.Time
}

func NewPeer(ip net.IP, port uint16) Peer {
    return Peer{
        IP:                  ip,
        Port:                port,
        LastMessageReceived: time.Now(),
        LastFindNodeRequest: time.Time{},
    }
}

func (p *Peer) Decode(data []byte) {
    ip := net.IP(data[:4])
    port := binary.BigEndian.Uint16(data[4:6])

    *p = NewPeer(ip, port)
}

func (p Peer) Encode() (b []byte) {
    b = append(b, p.IP.To4()...)
    binary.BigEndian.PutUint16(b, p.Port)
    return
}

func (p Peer) CanRequestFindNode() bool {
    now := time.Now()

    if p.LastFindNodeRequest.Add(time.Minute).After(now) {
        return false
    }

    return true
}

func (p Peer) IsGood() bool {
    return p.LastMessageReceived.Add(time.Minute * 15).After(time.Now())
}

func (p *Peer) UpdateLastRequestFindNode() {
    p.LastFindNodeRequest = time.Now()
}

func (p *Peer) UpdateLastMessageReceived() {
    p.LastMessageReceived = time.Now()
}

func (p Peer) Send(conn *net.UDPConn, data []byte) {
    destAddr := net.UDPAddr{IP: p.IP, Port: int(p.Port)}
    _, err := conn.WriteToUDP(data, &destAddr)
    if err != nil {
        log.Error("peer.Send", err.Error())
    }
}

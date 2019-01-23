package datastructure

import (
    "encoding/binary"
    log "github.com/sirupsen/logrus"
    "net"
    "sync"
    "time"
)

type Peer struct {
    sync.RWMutex
    ip                  net.IP
    port                uint16
    LastMessageReceived time.Time
    LastFindNodeRequest time.Time
}

func NewPeer(ip net.IP, port uint16) Peer {
    return Peer{
        ip:                  ip,
        port:                port,
        LastMessageReceived: time.Now(), // pretend the node is good for 1 minute on init
    }
}

func (p *Peer) Decode(data []byte) {
    ip := net.IP(data[:4])
    port := binary.BigEndian.Uint16(data[4:6])

    *p = NewPeer(ip, port)
}

func (p Peer) Encode() (b []byte) {
    b = append(b, p.ip.To4()...)
    binary.BigEndian.PutUint16(b, p.port)
    return
}

func (p *Peer) CanRequestFindNode() bool {
    p.RLock()
    defer p.RUnlock()

    now := time.Now()

    if p.LastFindNodeRequest.Add(time.Minute).After(now) {
        return false
    }

    return true
}

func (p *Peer) IsGood() bool {
    p.RLock()
    defer p.RUnlock()

    return p.LastMessageReceived.Add(time.Minute * 15).After(time.Now())
}

func (p *Peer) UpdateLastRequestFindNode() {
    p.Lock()
    defer p.Unlock()

    p.LastFindNodeRequest = time.Now()
}

func (p *Peer) UpdateLastMessageReceived() {
    p.Lock()
    defer p.Unlock()

    p.LastMessageReceived = time.Now()
}

func (p Peer) Send(conn *net.UDPConn, data []byte) {
    destAddr := net.UDPAddr{IP: p.ip, Port: int(p.port)}
    _, err := conn.WriteToUDP(data, &destAddr)
    if err != nil {
        log.Error("peer.Send", err.Error())
    }
}

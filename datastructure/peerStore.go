package datastructure

import "sync"

type PeerStore struct {
    sync.RWMutex
    store map[string][]Peer
}

func NewPeerStore() PeerStore{
    return PeerStore{
        store:   make(map[string][]Peer),
    }
}

func (p *PeerStore) Add(infoHash InfoHash, peers ...Peer) {
    p.Lock()
    defer p.Unlock()

    v, exists := p.store[infoHash.String()]
    if exists {
        v = append(v, peers...)
        p.store[infoHash.String()] = v
    } else {
        p.store[infoHash.String()] = peers
    }
}

func (p *PeerStore) Contains(infoHash InfoHash) bool {
    p.RLock()
    defer p.RUnlock()

    _, exists := p.store[infoHash.String()]
    return exists
}

func (p *PeerStore) Get(infoHash InfoHash) []Peer {
    p.RLock()
    defer p.RUnlock()

    v, _ := p.store[infoHash.String()]
    return v
}



package datastructure

type PeerStore map[string][]Peer

func (p *PeerStore) Add(infoHash InfoHash, peers []Peer) {
    v, exists := (*p)[infoHash.String()]
    if exists {
        v = append(v, peers...)
        (*p)[infoHash.String()] = v
    } else {
        (*p)[infoHash.String()] = peers
    }
}

func (p *PeerStore) Contains(infoHash InfoHash) bool {
    _, exists := (*p)[infoHash.String()]
    return exists
}

func (p *PeerStore) Get(infoHash InfoHash) []Peer {
    v, _ := (*p)[infoHash.String()]
    return v
}


